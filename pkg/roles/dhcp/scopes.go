package dhcp

import (
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"time"

	"beryju.io/ddet/pkg/ipam"
	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/dhcp/types"
	"github.com/insomniacslk/dhcp/dhcpv4"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

type Option struct {
	Tag     *uint8   `json:"tag"`
	TagName *string  `json:"tagName"`
	Value   *string  `json:"value"`
	Value64 []string `json:"value64"`
}

var TagMap map[string]uint8 = map[string]uint8{
	"subnet_mask": dhcpv4.OptionSubnetMask.Code(),
	"router":      dhcpv4.OptionRouter.Code(),
	"time_server": dhcpv4.OptionTimeServer.Code(),
	"name_server": dhcpv4.OptionNameServer.Code(),
	"domain_name": dhcpv4.OptionDomainName.Code(),
	"bootfile":    dhcpv4.OptionBootfileName.Code(),
	"tftp_server": dhcpv4.OptionTFTPServerName.Code(),
}

type Scope struct {
	Name string `json:"-"`

	SubnetCIDR string    `json:"subnetCidr"`
	Default    bool      `json:"default"`
	Options    []*Option `json:"options"`
	TTL        int64     `json:"ttl"`
	Range      struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"range"`
	DNS struct {
		Zone              string   `json:"zone"`
		Search            []string `json:"search"`
		AddZoneInHostname bool     `json:"addZoneInHostname"`
	} `json:"dns"`

	etcdKey string
	ipam    ipam.IPAM
	inst    roles.Instance
	role    *DHCPRole
	log     *log.Entry
}

func (r *DHCPRole) scopeFromKV(raw *mvccpb.KeyValue) (*Scope, error) {
	s := &Scope{
		inst: r.i,
		role: r,
		TTL:  int64((7 * 24 * time.Hour).Seconds()),
	}
	err := json.Unmarshal(raw.Value, &s)
	if err != nil {
		return nil, err
	}
	prefix := r.i.KV().Key(types.KeyRole, types.KeyScopes, "")
	s.Name = strings.TrimPrefix(string(raw.Key), prefix)
	// Get full etcd key without leading slash since this usually gets passed to Instance Key()
	s.etcdKey = string(raw.Key)[1:]

	s.log = r.log.WithField("scope", s.Name)

	// TODO: other IPAMs
	var ipamInst ipam.IPAM
	ipamInst, err = ipam.NewInternalIPAM(s.SubnetCIDR, s.Range.Start, s.Range.End)
	if err != nil {
		return nil, fmt.Errorf("failed to create ipam: %w", err)
	}
	s.ipam = ipamInst

	return s, nil
}

func (s *Scope) match(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) bool {
	return false
}

func (s *Scope) createLeaseFor(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) *Lease {
	ident := m.ClientHWAddr.String()
	lease := &Lease{
		Identifier: ident,

		Hostname: m.HostName(),
		Address:  s.ipam.NextFreeAddress().String(),
		ScopeKey: s.Name,

		inst:  s.inst,
		log:   s.log.WithField("lease", ident),
		scope: s,
	}
	if requestIp, ok := netip.AddrFromSlice(m.Options.Get(dhcpv4.OptionRequestedIPAddress)); ok {
		s.log.WithField("ip", requestIp).Debug("checking requested IP")
		if s.ipam.IsIPFree(requestIp) {
			s.log.WithField("ip", requestIp).Debug("requested IP is free")
			lease.Address = requestIp.String()
		}
	}
	return lease
}
