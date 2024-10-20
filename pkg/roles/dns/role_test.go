package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func RoleConfig() []byte {
	return []byte(tests.MustJSON(dns.RoleConfig{
		Port: 1054,
	}))
}

func TestRoleDNSStartNoConfig(t *testing.T) {
	defer tests.Setup(t)()
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte{}))
	role.Stop()
}

func TestRoleDNSStartEmptyConfig(t *testing.T) {
	defer tests.Setup(t)()
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte("{}")))
	role.Stop()
}
