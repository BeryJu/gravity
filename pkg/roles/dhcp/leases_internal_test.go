package dhcp

import (
	"context"
	"errors"
	"testing"

	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/storage"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestFindLeaseInStore_GetError(t *testing.T) {
	ctx := setupDHCPInternalTest(t)
	inst := newDHCPTestInstance(ctx)
	inst.kv = inst.KV().WithHooks(storage.StorageHook{
		GetPre: func(context.Context, string, ...clientv3.OpOption) error {
			return errors.New("boom")
		},
	})
	role := New(inst)

	req := role.NewRequest4(&dhcpv4.DHCPv4{
		ClientHWAddr: []byte{0xb2, 0xb7, 0x86, 0x2c, 0xd3, 0xfa},
	})

	assert.Nil(t, role.FindLeaseInStore(req))
}

func TestFindLeaseInStore_EmptyResult(t *testing.T) {
	ctx := setupDHCPInternalTest(t)
	inst := newDHCPTestInstance(ctx)
	role := New(inst)

	req := role.NewRequest4(&dhcpv4.DHCPv4{
		ClientHWAddr: []byte{0xb2, 0xb7, 0x86, 0x2c, 0xd3, 0xfa},
	})

	assert.Nil(t, role.FindLeaseInStore(req))
}

func TestFindLeaseInStore_ParseError(t *testing.T) {
	ctx := setupDHCPInternalTest(t)
	inst := newDHCPTestInstance(ctx)
	role := New(inst)

	panicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyLeases,
			"b2:b7:86:2c:d3:fa",
		).String(),
		"{",
	))

	req := role.NewRequest4(&dhcpv4.DHCPv4{
		ClientHWAddr: []byte{0xb2, 0xb7, 0x86, 0x2c, 0xd3, 0xfa},
	})

	assert.Nil(t, role.FindLeaseInStore(req))
}
