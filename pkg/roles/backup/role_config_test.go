package backup_test

import (
	"testing"

	"beryju.io/gravity/pkg/roles/backup"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func Test_APIRoleConfigGet(t *testing.T) {
	role := getRole()
	defer role.Stop()

	var output backup.APIRoleBackupConfigOutput
	assert.NoError(t, role.APIRoleConfigGet().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}

func Test_APIRoleConfigPut(t *testing.T) {
	role := getRole()
	defer role.Stop()

	var output struct{}
	assert.NoError(t, role.APIRoleConfigPut().Interact(tests.Context(), backup.APIRoleBackupConfigInput{
		Config: backup.RoleConfig{
			Endpoint: "foo",
		},
	}, &output))
}
