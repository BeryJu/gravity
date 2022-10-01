package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func Test_APIHandlerRoleConfigGet(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("dns")
	role := dns.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output dns.APIRoleDNSConfigOutput
	assert.NoError(t, role.APIHandlerRoleConfigGet().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}

func Test_APIHandlerRoleConfigPut(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("dns")
	role := dns.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output struct{}
	assert.NoError(t, role.APIHandlerRoleConfigPut().Interact(tests.Context(), dns.APIRoleDNSConfigInput{
		Config: dns.RoleConfig{
			Port: 1054,
		},
	}, &output))
}
