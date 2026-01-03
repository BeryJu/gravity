package backup_test

import (
	"testing"

	"beryju.io/gravity/pkg/roles/backup"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIRoleConfigGet(t *testing.T) {
	tests.Setup(t)
	role := getRole()
	defer role.Stop()

	var output backup.APIRoleConfigOutput
	assert.NoError(t, role.APIRoleConfigGet().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPIRoleConfigPut(t *testing.T) {
	tests.Setup(t)
	role := getRole()
	defer role.Stop()

	assert.NoError(t, role.APIRoleConfigPut().Interact(tests.Context(), backup.APIRoleConfigInput{
		Config: backup.RoleConfig{
			Endpoint: "foo",
		},
	}, &struct{}{}))
}
