package options

import (
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dhcp/options/types"
	dhcpTypes "beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"google.golang.org/protobuf/proto"
)

func Bootstrap(r roles.Instance) func(ev *roles.Event) {
	return func(ev *roles.Event) {
		opts := []*types.OptionDefinition{
			{
				Name:        "Subnet Mask",
				Code:        int32(dhcpv4.OptionSubnetMask.Code()),
				IsArray:     false,
				Type:        types.DataType_IP_ADDRESS,
				Description: "Set the subnet mask",
			},
			{
				Name:        "Router",
				Code:        int32(dhcpv4.OptionRouter.Code()),
				IsArray:     false,
				Type:        types.DataType_IP_ADDRESS,
				Description: "Set the router",
			},
			{
				Name:        "Time server",
				Code:        int32(dhcpv4.OptionTimeServer.Code()),
				IsArray:     true,
				Type:        types.DataType_IP_ADDRESS,
				Description: "Set the time server(s)",
			},
			{
				Name:        "Domain name server",
				Code:        int32(dhcpv4.OptionDomainNameServer.Code()),
				IsArray:     true,
				Type:        types.DataType_IP_ADDRESS,
				Description: "Set the DNS server(s)",
			},
			{
				Name:        "Domain name",
				Code:        int32(dhcpv4.OptionDomainName.Code()),
				IsArray:     false,
				Type:        types.DataType_STRING,
				Description: "Set the DNS domain name",
			},
			{
				Name:        "Bootfile",
				Code:        int32(dhcpv4.OptionBootfileName.Code()),
				IsArray:     false,
				Type:        types.DataType_STRING,
				Description: "Set the PXE Bootfile",
			},
			{
				Name:        "TFTP Server",
				Code:        int32(dhcpv4.OptionTFTPServerName.Code()),
				IsArray:     false,
				Type:        types.DataType_STRING,
				Description: "Set the TFTP Server name",
			},
		}
		for _, opt := range opts {
			rv, err := proto.Marshal(opt)
			if err != nil {
				continue
			}
			_, err = r.KV().Put(ev.Context, r.KV().Key(
				dhcpTypes.KeyRole,
				dhcpTypes.KeyOptionDefinitions,
				opt.Name,
			).String(), string(rv))
			if err != nil {
				continue
			}
		}
	}
}
