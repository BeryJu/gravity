package dhcp

import (
	"context"
	"errors"
	"testing"
	"time"

	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/storage"
	"beryju.io/gravity/pkg/tests"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func testLeaseScope(t *testing.T, ctx context.Context, inst *dhcpTestInstance, role *Role) *Scope {
	t.Helper()
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test",
		).String(),
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

	scope, ok := role.scopes.GetPrefix("test")
	assert.True(t, ok)
	assert.NotNil(t, scope)
	return scope
}

func TestLease_IsReservation(t *testing.T) {
	l := &Lease{}
	assert.False(t, l.IsReservation())

	l.Reservation = true
	assert.True(t, l.IsReservation())
}

func TestLeaseFromKV_MigratesLegacyReservation(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()
	inst := newDHCPTestInstance(ctx)
	role := New(inst)
	testLeaseScope(t, ctx, inst, role)

	identifier := "b2:b7:86:2c:d3:fa"
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyLeases,
			identifier,
		).String(),
		`{"address":"10.100.0.100","hostname":"legacy","scopeKey":"test","expiry":-1,"description":""}`,
	))

	res, err := inst.KV().Get(ctx, inst.KV().Key(
		types.KeyRole,
		types.KeyLeases,
		identifier,
	).String())
	tests.PanicIfError(err)
	assert.Len(t, res.Kvs, 1)

	l, err := role.leaseFromKV(res.Kvs[0])
	assert.NoError(t, err)
	assert.True(t, l.Reservation)
	assert.True(t, l.IsReservation())
}

func TestLeaseFromKV_DoesNotMarkRegularLeaseAsReservation(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()
	inst := newDHCPTestInstance(ctx)
	role := New(inst)
	testLeaseScope(t, ctx, inst, role)

	identifier := "b2:b7:86:2c:d3:fa"
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyLeases,
			identifier,
		).String(),
		`{"address":"10.100.0.100","hostname":"regular","scopeKey":"test","expiry":1234,"description":""}`,
	))

	res, err := inst.KV().Get(ctx, inst.KV().Key(
		types.KeyRole,
		types.KeyLeases,
		identifier,
	).String())
	tests.PanicIfError(err)
	assert.Len(t, res.Kvs, 1)

	l, err := role.leaseFromKV(res.Kvs[0])
	assert.NoError(t, err)
	assert.False(t, l.IsReservation())
}

func TestLeasePut_ReservationTracksExpiryWithoutEtcdLease(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()
	inst := newDHCPTestInstance(ctx)
	role := New(inst)
	scope := testLeaseScope(t, ctx, inst, role)

	lease := role.NewLease("b2:b7:86:2c:d3:fa")
	lease.scope = scope
	lease.ScopeKey = scope.Name
	lease.Address = "10.100.0.100"
	lease.Reservation = true

	tests.PanicIfError(lease.Put(ctx, 3600))

	assert.Greater(t, lease.Expiry, time.Now().Unix(), "expiry should still be tracked for informational purposes")

	res, err := inst.KV().Get(ctx, inst.KV().Key(
		types.KeyRole,
		types.KeyLeases,
		lease.Identifier,
	).String())
	tests.PanicIfError(err)
	assert.Len(t, res.Kvs, 1)
	assert.Zero(t, res.Kvs[0].Lease, "reservations must not expire from etcd")
}

func TestLeasePut_NonReservationAttachesEtcdLease(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()
	inst := newDHCPTestInstance(ctx)
	role := New(inst)
	scope := testLeaseScope(t, ctx, inst, role)

	lease := role.NewLease("b2:b7:86:2c:d3:fa")
	lease.scope = scope
	lease.ScopeKey = scope.Name
	lease.Address = "10.100.0.100"

	tests.PanicIfError(lease.Put(ctx, 3600))

	res, err := inst.KV().Get(ctx, inst.KV().Key(
		types.KeyRole,
		types.KeyLeases,
		lease.Identifier,
	).String())
	tests.PanicIfError(err)
	assert.Len(t, res.Kvs, 1)
	assert.NotZero(t, res.Kvs[0].Lease, "regular leases must expire from etcd")
}

func TestCreateLeaseIfAbsent_ReservationDoesNotAttachEtcdLease(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()
	inst := newDHCPTestInstance(ctx)
	role := New(inst)
	scope := testLeaseScope(t, ctx, inst, role)

	lease := role.NewLease("b2:b7:86:2c:d3:fa")
	lease.scope = scope
	lease.ScopeKey = scope.Name
	lease.Address = "10.100.0.100"
	lease.Reservation = true

	created, isNew, err := role.CreateLeaseIfAbsent(ctx, lease, 3600)
	assert.NoError(t, err)
	assert.True(t, isNew)
	assert.NotNil(t, created)

	res, err := inst.KV().Get(ctx, inst.KV().Key(
		types.KeyRole,
		types.KeyLeases,
		lease.Identifier,
	).String())
	tests.PanicIfError(err)
	assert.Len(t, res.Kvs, 1)
	assert.Zero(t, res.Kvs[0].Lease, "reservations must not expire from etcd")
}

func TestFindLeaseInStore_GetError(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()
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
	tests.Setup(t)
	ctx := tests.Context()
	inst := newDHCPTestInstance(ctx)
	role := New(inst)

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
