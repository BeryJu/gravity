package dhcp

import (
	"context"
	"errors"
	"testing"

	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/storage"
	"beryju.io/gravity/pkg/tests"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func startLeaseTestRole(t *testing.T, inst *dhcpTestInstance, role *Role, ctx context.Context) {
	t.Helper()
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "test").String(),
		tests.MustJSON(Scope{
			SubnetCIDR: "10.100.0.0/24",
			Default:    true,
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.100.0.100",
				"range_end":   "10.100.0.250",
			},
		}),
	))
	tests.PanicIfError(role.Start(ctx, []byte(tests.MustJSON(RoleConfig{Port: 0}))))
	t.Cleanup(role.Stop)
}

func TestFindLeaseInStore_GetError(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()
	inst := newDHCPTestInstance(ctx)
	role := New(inst)
	startLeaseTestRole(t, inst, role, ctx)
	inst.kv = inst.KV().WithHooks(storage.StorageHook{
		GetPre: func(context.Context, string, ...clientv3.OpOption) error {
			return errors.New("boom")
		},
	})

	req := role.NewRequest4(&dhcpv4.DHCPv4{
		ClientHWAddr: []byte{0xb2, 0xb7, 0x86, 0x2c, 0xd3, 0xfa},
	})

	assert.Nil(t, role.FindLeaseInStore(req))
}

func TestFindLeaseInStore_EmptyResult(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()
	inst := newDHCPTestInstance(ctx)
	role := New(inst)
	startLeaseTestRole(t, inst, role, ctx)

	req := role.NewRequest4(&dhcpv4.DHCPv4{
		ClientHWAddr: []byte{0xb2, 0xb7, 0x86, 0x2c, 0xd3, 0xfa},
	})

	assert.Nil(t, role.FindLeaseInStore(req))
}

func TestFindLeaseInStore_ParseError(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()
	inst := newDHCPTestInstance(ctx)
	role := New(inst)
	startLeaseTestRole(t, inst, role, ctx)

	tests.PanicIfError(inst.KV().Put(
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

func TestFindLeaseInStore_Reservation(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()
	inst := newDHCPTestInstance(ctx)
	role := New(inst)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "test").String(),
		tests.MustJSON(Scope{
			SubnetCIDR: "10.100.0.0/24",
			Default:    true,
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.100.0.100",
				"range_end":   "10.100.0.250",
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		role.reservationKey("test", "b2:b7:86:2c:d3:fa").String(),
		tests.MustJSON(Lease{
			ScopeKey: "test",
			Address:  "10.100.0.150",
			Expiry:   -1,
		}),
	))
	tests.PanicIfError(role.Start(ctx, []byte(tests.MustJSON(RoleConfig{Port: 0}))))
	defer role.Stop()

	req := role.NewRequest4(&dhcpv4.DHCPv4{
		ClientHWAddr: []byte{0xb2, 0xb7, 0x86, 0x2c, 0xd3, 0xfa},
	})
	lease := role.FindLeaseInStore(req)
	assert.NotNil(t, lease)
	assert.Equal(t, "10.100.0.150", lease.Address)
	assert.True(t, lease.IsReservation())
}

func TestFindLeaseInStore_MigratesLegacyReservation(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()
	inst := newDHCPTestInstance(ctx)
	role := New(inst)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "test").String(),
		tests.MustJSON(Scope{
			SubnetCIDR: "10.100.0.0/24",
			Default:    true,
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.100.0.100",
				"range_end":   "10.100.0.250",
			},
		}),
	))
	identifier := "b2:b7:86:2c:d3:fa"
	tests.PanicIfError(inst.KV().Put(
		ctx,
		role.leaseKey(identifier).String(),
		tests.MustJSON(Lease{
			ScopeKey: "test",
			Address:  "10.100.0.150",
			Expiry:   -1,
		}),
	))
	tests.PanicIfError(role.Start(ctx, []byte(tests.MustJSON(RoleConfig{Port: 0}))))
	defer role.Stop()

	req := role.NewRequest4(&dhcpv4.DHCPv4{
		ClientHWAddr: []byte{0xb2, 0xb7, 0x86, 0x2c, 0xd3, 0xfa},
	})
	lease := role.FindLeaseInStore(req)
	assert.NotNil(t, lease)
	assert.Equal(t, "10.100.0.150", lease.Address)
	tests.AssertEtcd(t, inst.KV(), role.leaseKey(identifier))
	tests.AssertEtcd(t, inst.KV(), role.reservationKey("test", identifier), Lease{
		ScopeKey: "test",
		Address:  "10.100.0.150",
		Expiry:   -1,
	})
}
