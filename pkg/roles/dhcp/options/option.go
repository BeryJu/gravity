package options

import (
	"net"

	"beryju.io/gravity/pkg/roles/dhcp/options/types"
	"beryju.io/gravity/pkg/storage/watcher"
	"github.com/insomniacslk/dhcp/dhcpv4"
)

type Option struct {
	*types.Option
	def *types.OptionDefinition
}

func ApplyToResponse(
	options []*types.Option,
	defs *watcher.Watcher[*types.OptionDefinition],
	req *dhcpv4.DHCPv4,
	reply *dhcpv4.DHCPv4,
) {
	for _, ropt := range options {
		def, found := defs.GetPrefix(ropt.Def)
		if !found {
			continue
		}
		opt := &Option{
			Option: ropt,
			def:    def,
		}
		reply.UpdateOption(opt.toDHCPv4())
	}
}

func (op *Option) toDHCPv4() dhcpv4.Option {
	v := []byte{}
	for idx, vv := range op.Data {
		if idx > 1 && !op.def.IsArray {
			break
		}
		v = append(v, op.encodeSingle(vv)...)
	}
	return dhcpv4.OptGeneric(dhcpv4.GenericOptionCode(op.def.Code), v)
}

func (op *Option) encodeSingle(raw *types.OptionData) []byte {
	switch op.def.Type {
	case types.DataType_BYTE:
		return raw.GetDataByte()
	case types.DataType_STRING:
		return dhcpv4.String(raw.GetDataString()).ToBytes()
	case types.DataType_IP_ADDRESS:
		return dhcpv4.IP(net.IP(raw.GetIp()))
	case types.DataType_IP_MASK:
		return dhcpv4.IPMask(net.IPMask(raw.GetIp()))
	}
	return []byte{} // unreachable
}
