package ms_dhcp

import (
	"context"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Converter struct {
	a  *api.APIClient
	in DHCPServer
	l  *zap.Logger
}

func New(api *api.APIClient, input string) (*Converter, error) {
	x, err := os.ReadFile(input)
	if err != nil {
		return nil, err
	}
	var dhcps DHCPServer
	err = xml.Unmarshal(x, &dhcps)
	if err != nil {
		return nil, err
	}
	return &Converter{
		a:  api,
		in: dhcps,
		l:  extconfig.Get().Logger().Named("convert.ms_dhcp"),
	}, nil
}

func (c *Converter) Run(ctx context.Context) {
	for _, scope := range c.in.IPv4.Scopes.Scope {
		c.convertScope(scope, ctx)
	}
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
			"type":  "internal",
			"start": sc.StartRange,
			"end":   sc.EndRange,
		},
		Options: []api.TypesDHCPOption{},
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
	_, err = c.a.RolesDhcpApi.DhcpPutScopes(ctx).Scope(sc.Name).DhcpAPIScopesPutInput(gscope).Execute()
	if err != nil {
		return err
	}

	for _, res := range sc.Reservations.Reservation {
		l := c.convertReservation(sc.Name, ctx, res)
		if l != nil {
			c.l.Warn("failed to convert reservation", zap.Error(err))
			continue
		}
	}
	for _, l := range sc.Leases.Lease {
		ll := c.convertLease(sc.Name, ctx, l)
		if ll != nil {
			c.l.Warn("failed to convert lease", zap.Error(err))
			continue
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
	lease := api.DhcpAPILeasesPutInput{
		Address:  r.IPAddress,
		Hostname: r.Name,
	}
	_, err := c.a.RolesDhcpApi.DhcpPutLeases(ctx).Scope(scope).Identifier(c.getIdentifier(r.ClientId)).DhcpAPILeasesPutInput(lease).Execute()
	return err
}

func (c *Converter) convertLease(scope string, ctx context.Context, l Lease) error {
	if l.HostName == "BAD_ADDRESS" {
		return nil
	}
	lease := api.DhcpAPILeasesPutInput{
		Address:  l.IPAddress,
		Hostname: l.HostName,
	}
	_, err := c.a.RolesDhcpApi.DhcpPutLeases(ctx).Scope(scope).Identifier(c.getIdentifier(l.ClientId)).DhcpAPILeasesPutInput(lease).Execute()
	return err
}
