package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIRoleConfigGet(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	var output dns.APIRoleConfigOutput
	assert.NoError(t, role.APIRoleConfigGet().Interact(ctx, struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPIRoleConfigPut(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	assert.NoError(t, role.APIRoleConfigPut().Interact(ctx, dns.APIRoleConfigInput{
		Config: dns.RoleConfig{
			Port: 0,
		},
	}, &struct{}{}))
}
