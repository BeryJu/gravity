package ms_dhcp

import (
	"context"
	"encoding/xml"
	"fmt"
	"net"
	"strconv"
	"strings"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Converter struct {
	a     *api.APIClient
	l     *zap.Logger
	in    DHCPServer
	scope string
}

type ConverterOption struct {
	apply func(*Converter)
}

func WithExistingScope(name string) ConverterOption {
	return ConverterOption{
		apply: func(c *Converter) {
			c.scope = name
		},
	}
}

func New(api *api.APIClient, input string, options ...ConverterOption) (*Converter, error) {
	var dhcps DHCPServer
	err := xml.Unmarshal([]byte(input), &dhcps)
	if err != nil {
		return nil, err
	}
	conv := &Converter{
		a:  api,
		in: dhcps,
		l:  extconfig.Get().Logger().Named("convert.ms_dhcp"),
	}
	for _, opt := range options {
		opt.apply(conv)
	}
	return conv, nil
}

func (c *Converter) Run(ctx context.Context) error {
	for _, scope := range c.in.IPv4.Scopes.Scope {
		if c.scope != "" && scope.Name != c.scope {
			continue
		}
		err := c.convertScope(scope, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Converter) convertScope(sc Scope, ctx context.Context) error {
	// Build CIDR
	m := net.IPMask(net.ParseIP(sc.SubnetMask).To4())
	ones, _ := m.Size()
	_, cidr, err := net.ParseCIDR(fmt.Sprintf("%s/%d", sc.ScopeId, ones))
	if err != nil {
		return errors.Wrap(err, "failed to parse CIDR")
	}
	// Build lease duration
	// saved as days:hours:minutes
	// rdur := strings.Split(scope.LeaseDuration, ":")
	// dur := time.Duration(0)
	// // days
	// day, err := strconv.Atoi(rdur[0])
	// if err != nil {
	// 	log.Println(err)
	// 	continue
	// }
	// dur += day * 24 * time.Hour
	gscope := api.DhcpAPIScopesPutInput{
		Default:    false,
		SubnetCidr: cidr.String(),
		Ipam: map[string]string{
			"type":        "internal",
			"range_start": sc.StartRange,
			"range_end":   sc.EndRange,
		},
		Options: []api.TypesDHCPOption{},
		Hook:    "",
	}
	for _, optv := range sc.OptionValues.OptionValue {
		tag, err := strconv.Atoi(optv.OptionId)
		if err != nil {
			c.l.Error("failed to convert optionID to int", zap.Error(err))
			continue
		}
		t := int32(tag)
		v := optv.Value[0]
		gscope.Options = append(gscope.Options, api.TypesDHCPOption{
			Tag:   *api.NewNullableInt32(&t),
			Value: *api.NewNullableString(&v),
		})
	}
	name := slug.Make(sc.Name)
	_, err = c.a.RolesDhcpAPI.DhcpPutScopes(ctx).Scope(name).DhcpAPIScopesPutInput(gscope).Execute()
	if err != nil {
		c.l.Warn("failed to convert scope", zap.Error(err))
		return err
	}
	c.l.Info("converted scope", zap.String("name", name))

	for _, res := range sc.Reservations.Reservation {
		l := c.convertReservation(name, ctx, res)
		if l != nil {
			c.l.Warn("failed to convert reservation", zap.Error(err))
		} else {
			c.l.Info("converted reservation", zap.String("name", res.Name))
		}
	}
	for _, l := range sc.Leases.Lease {
		ll := c.convertLease(name, ctx, l)
		if ll != nil {
			c.l.Warn("failed to convert lease", zap.Error(err))
		} else {
			c.l.Info("converted lease", zap.String("name", l.HostName))
		}
	}
	return nil
}

func (c *Converter) getIdentifier(clientId string) string {
	if strings.Count(clientId, "-") == 5 {
		return strings.ReplaceAll(clientId, "-", ":")
	}
	return strings.ReplaceAll(clientId, "-", "")
}

func (c *Converter) convertReservation(scope string, ctx context.Context, r Reservation) error {
	if ip := net.ParseIP(r.IPAddress); ip == nil {
		return fmt.Errorf("failed to parse IP")
	}
	lease := api.DhcpAPILeasesPutInput{
		Address:  r.IPAddress,
		Hostname: r.Name,
		Expiry:   api.PtrInt32(-1),
	}
	_, err := c.a.RolesDhcpAPI.DhcpPutLeases(ctx).Scope(scope).Identifier(c.getIdentifier(r.ClientId)).DhcpAPILeasesPutInput(lease).Execute()
	return err
}

func (c *Converter) convertLease(scope string, ctx context.Context, l Lease) error {
	if l.HostName == "BAD_ADDRESS" {
		return nil
	}
	if ip := net.ParseIP(l.IPAddress); ip == nil {
		return fmt.Errorf("failed to parse IP")
	}
	lease := api.DhcpAPILeasesPutInput{
		Address:  l.IPAddress,
		Hostname: l.HostName,
	}
	_, err := c.a.RolesDhcpAPI.DhcpPutLeases(ctx).Scope(scope).Identifier(c.getIdentifier(l.ClientId)).DhcpAPILeasesPutInput(lease).Execute()
	return err
}
