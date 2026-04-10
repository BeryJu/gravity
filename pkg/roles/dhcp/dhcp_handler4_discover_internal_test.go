package dhcp

import (
	"testing"
	"time"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/assert"
)

func TestDHCPDiscover_ReusesExistingLeaseWithoutDowngradingTTL(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := New(inst)

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

	tests.PanicIfError(role.Start(ctx, []byte(tests.MustJSON(RoleConfig{
		Port:                  0,
		LeaseNegotiateTimeout: 30,
	}))))
	defer role.Stop()

	scope, ok := role.scopes.GetPrefix("test")
	assert.True(t, ok)
	assert.NotNil(t, scope)

	lease := role.NewLease("b2:b7:86:2c:d3:fa")
	lease.scope = scope
	lease.ScopeKey = scope.Name
	lease.Address = "10.100.0.100"
	tests.PanicIfError(lease.Put(ctx, 3600))

	assert.Eventually(t, func() bool {
		match, ok := role.leases.GetPrefix(lease.Identifier)
		return ok && match != nil && match.Address == lease.Address
	}, time.Second, 10*time.Millisecond)

	role.leases.mutex.Lock()
	role.leases.entries = map[string]*Lease{}
	role.leases.mutex.Unlock()

	req := &dhcpv4.DHCPv4{
		OpCode:       dhcpv4.OpcodeBootRequest,
		ClientHWAddr: []byte{0xb2, 0xb7, 0x86, 0x2c, 0xd3, 0xfa},
	}
	req.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeDiscover))

	req4 := role.NewRequest4(req)
	res := role.HandleDHCPDiscover4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, lease.Address, res.YourIPAddr.String())

	stored := role.FindLeaseInStore(req4)
	assert.NotNil(t, stored)
	assert.Greater(t, stored.Expiry, time.Now().Add(5*time.Minute).Unix())
}
