package backup_test

import (
	"testing"

	"beryju.io/gravity/pkg/roles/backup"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func Test_APIHandlerRoleConfigGet(t *testing.T) {
	role := getRole()
	defer role.Stop()

	var output backup.RoleBackupConfigOutput
	assert.NoError(t, role.APIHandlerRoleConfigGet().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}

func Test_APIHandlerRoleConfigPut(t *testing.T) {
	role := getRole()
	defer role.Stop()

	var output struct{}
	assert.NoError(t, role.APIHandlerRoleConfigPut().Interact(tests.Context(), backup.RoleBackupConfigInput{
		Config: backup.RoleConfig{
			Endpoint: "foo",
		},
	}, &output))
}
