package options

import (
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dhcp/options/types"
	dhcpTypes "beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"google.golang.org/protobuf/proto"
)

func OptionsIP() []dhcpv4.OptionCode {
	return []dhcpv4.OptionCode{
		dhcpv4.OptionBroadcastAddress,
		dhcpv4.OptionRequestedIPAddress,
		dhcpv4.OptionServerIdentifier,
	}
}

func OptionsIPS() []dhcpv4.OptionCode {
	return []dhcpv4.OptionCode{
		dhcpv4.OptionRouter,
		dhcpv4.OptionNTPServers,
		dhcpv4.OptionDomainNameServer,
	}
}

func OptionsMask() []dhcpv4.OptionCode {
	return []dhcpv4.OptionCode{
		dhcpv4.OptionSubnetMask,
	}
}

func OptionsString() []dhcpv4.OptionCode {
	return []dhcpv4.OptionCode{
		dhcpv4.OptionHostName,
		dhcpv4.OptionDomainName,
		dhcpv4.OptionRootPath,
		dhcpv4.OptionBootfileName,
		dhcpv4.OptionTFTPServerName,
		dhcpv4.OptionClassIdentifier,
		dhcpv4.OptionUserClassInformation,
		// dhcpv4.OptionMessage,
	}
}

func Bootstrap(r roles.Instance) func(ev *roles.Event) {
	return func(ev *roles.Event) {
		opts := []*types.OptionDefinition{}
		for _, opt := range OptionsIP() {
			opts = append(opts, &types.OptionDefinition{
				Name:    opt.String(),
				Code:    int32(opt.Code()),
				IsArray: false,
				Type:    types.DataType_IP_ADDRESS,
			})
		}
		for _, opt := range OptionsIPS() {
			opts = append(opts, &types.OptionDefinition{
				Name:    opt.String(),
				Code:    int32(opt.Code()),
				IsArray: true,
				Type:    types.DataType_IP_ADDRESS,
			})
		}
		for _, opt := range OptionsMask() {
			opts = append(opts, &types.OptionDefinition{
				Name:    opt.String(),
				Code:    int32(opt.Code()),
				IsArray: false,
				Type:    types.DataType_IP_MASK,
			})
		}
		for _, opt := range OptionsString() {
			opts = append(opts, &types.OptionDefinition{
				Name:    opt.String(),
				Code:    int32(opt.Code()),
				IsArray: false,
				Type:    types.DataType_STRING,
			})
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
