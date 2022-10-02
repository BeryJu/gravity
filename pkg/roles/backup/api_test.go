package backup_test

import (
	"testing"

	"beryju.io/gravity/pkg/roles/backup"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIBackupStart(t *testing.T) {
	role := getRole()
	defer role.Stop()

	var output backup.BackupStatus
	assert.NoError(t, role.APIBackupStart().Interact(tests.Context(), backup.APIBackupStartInput{
		Wait: true,
	}, &output))
	assert.NotNil(t, output)
}
