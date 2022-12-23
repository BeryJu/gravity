package ms_dhcp

import (
	"context"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"go.uber.org/zap"
)

type Converter struct {
	*instance.RoleInstance
	d  *dhcp.Role
	in DHCPServer
}

func New(inst *instance.RoleInstance, input string) (*Converter, error) {
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
		RoleInstance: inst,
		d:            dhcp.New(inst),
		in:           dhcps,
	}, nil
}

func (c *Converter) Run(ctx context.Context) {
	for _, scope := range c.in.IPv4.Scopes.Scope {
		c.convertScope(scope, ctx)
	}
}

func (c *Converter) convertScope(sc Scope, ctx context.Context) {
	// Build CIDR
	m := net.IPMask(net.ParseIP(sc.SubnetMask).To4())
	ones, _ := m.Size()
	_, cidr, err := net.ParseCIDR(fmt.Sprintf("%s/%d", sc.ScopeId, ones))
	if err != nil {
		c.Log().Warn("failed to parse cidr", zap.Error(err))
		return
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
	gscope := c.d.NewScope(sc.Name)
	gscope.SubnetCIDR = cidr.String()
	// gscope.DNS.
	gscope.IPAM = map[string]string{
		"type":  "internal",
		"start": sc.StartRange,
		"end":   sc.EndRange,
	}
	for _, optv := range sc.OptionValues.OptionValue {
		tag, err := strconv.Atoi(optv.OptionId)
		if err != nil {
			c.Log().Error("failed to convert optionID to int", zap.Error(err))
			continue
		}
		t := uint8(tag)
		v := optv.Value[0]
		gscope.Options = append(gscope.Options, &types.Option{
			Tag:   &t,
			Value: &v,
		})
	}
	gscope.Put(ctx, 0)

	for _, res := range sc.Reservations.Reservation {
		l := c.convertReservation(gscope, res)
		if l != nil {
			l.Put(ctx, 0)
		}
	}
	for _, l := range sc.Leases.Lease {
		ll := c.convertLease(gscope, l)
		if ll != nil {
			ll.Put(ctx, 0)
		}
	}
}

func (c *Converter) getIdentifier(clientId string) string {
	if strings.Count(clientId, "-") == 5 {
		return strings.ReplaceAll(clientId, "-", ":")
	}
	return strings.ReplaceAll(clientId, "-", "")
}

func (c *Converter) convertReservation(gs *dhcp.Scope, r Reservation) *dhcp.Lease {
	lease := c.d.NewLease(c.getIdentifier(r.ClientId))
	lease.Address = r.IPAddress
	lease.Hostname = r.Name
	lease.ScopeKey = gs.Name
	lease.AddressLeaseTime = ""
	return lease
}

func (c *Converter) convertLease(gs *dhcp.Scope, l Lease) *dhcp.Lease {
	if l.HostName == "BAD_ADDRESS" {
		return nil
	}
	lease := c.d.NewLease(c.getIdentifier(l.ClientId))
	lease.Address = l.IPAddress
	lease.Hostname = l.HostName
	lease.ScopeKey = gs.Name
	lease.AddressLeaseTime = ""
	return lease
}
