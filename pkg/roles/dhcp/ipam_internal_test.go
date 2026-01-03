package dhcp_test

import (
	"fmt"
	"sync"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestIPAMInternal(t *testing.T) {
	defer tests.Setup(t)()
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)
	scope := testScope()
	ipam, err := dhcp.NewInternalIPAM(role, &scope)
	assert.NoError(t, err)
	assert.NotNil(t, ipam)
}

func TestIPAMInternal_NextFreeAddress(t *testing.T) {
	defer tests.Setup(t)()
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)
	scope := testScope()
	ipam, err := dhcp.NewInternalIPAM(role, &scope)
	assert.NoError(t, err)

	next := ipam.NextFreeAddress("test")
	assert.NotNil(t, next)
	assert.Equal(t, tests.MustParseNetIP(t, "10.200.0.100"), *next)
}

func TestIPAMInternal_NextFreeAddress_UniqueParallel(t *testing.T) {
	t.Skip()
	defer tests.Setup(t)()
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)

	iter := 100
	addrs := make(chan string, iter)

	scope := testScope()

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			scope.Name,
		).String(),
		tests.MustJSON(scope),
	))
	// Create fake leases to test against
	for i := 0; i < iter-10; i++ {
		lease := testLease()
		lease.Address = fmt.Sprintf("10.200.0.%d", 100+i)
		lease.Identifier = fmt.Sprintf("test_%d", i)
		lease.ScopeKey = scope.Name
		tests.PanicIfError(inst.KV().Put(
			ctx,
			inst.KV().Key(
				types.KeyRole,
				types.KeyLeases,
				lease.Identifier,
			).String(),
			tests.MustJSON(lease),
		))
	}

	role := dhcp.New(inst)
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()
	ipam, err := dhcp.NewInternalIPAM(role, &scope)
	assert.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(iter)
	tester := func(idx int) {
		next := ipam.NextFreeAddress(fmt.Sprintf("test_%d", idx))
		assert.NotNil(t, next, "Iter %d", idx)

		lease := testLease()
		lease.Address = next.String()
		lease.Identifier = fmt.Sprintf("test_%d", idx)
		lease.ScopeKey = scope.Name
		tests.PanicIfError(inst.KV().Put(
			ctx,
			inst.KV().Key(
				types.KeyRole,
				types.KeyLeases,
				lease.Identifier,
			).String(),
			tests.MustJSON(lease),
		))

		addrs <- next.String()
		wg.Done()
	}
	for i := range iter {
		go tester(i)
	}
	wg.Wait()
	close(addrs)
	ips := map[string]int{}
	assert.Len(t, addrs, iter)
	for ip := range addrs {
		ips[ip] += 1
	}
	dupes := []string{}
	for ip, occur := range ips {
		if occur > 1 {
			dupes = append(dupes, ip)
		}
	}
	assert.Len(t, maps.Keys(ips), iter, "Addresses: %v, duplicates: %v", maps.Keys(ips), dupes)
	assert.Len(t, dupes, 0, "Addresses: %v, duplicates: %v", maps.Keys(ips), dupes)
}
