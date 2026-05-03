package backup_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/backup"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIBackupStart(t *testing.T) {
	tests.Setup(t)
	role := getRole()
	defer role.Stop()

	var output backup.BackupStatus
	assert.NoError(t, role.APIBackupStart().Interact(tests.Context(), backup.APIBackupStartInput{
		Wait: true,
	}, &output))
	assert.Equal(t, backup.BackupStatusSuccess, output.Status)
}

func TestAPIBackupStarNoWait(t *testing.T) {
	tests.Setup(t)
	role := getRole()
	defer role.Stop()

	var output backup.BackupStatus
	assert.NoError(t, role.APIBackupStart().Interact(tests.Context(), backup.APIBackupStartInput{
		Wait: false,
	}, &output))
	assert.Equal(t, backup.BackupStatusStarted, output.Status)
}

func TestAPIBackupStatus(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	rootInst.ForRole("backup", ctx)

	TestAPIBackupStart(t)

	role := getRole()
	defer role.Stop()

	var output backup.APIBackupStatusOutput
	assert.NoError(t, role.APIBackupStatus().Interact(ctx, struct{}{}, &output))
	assert.Equal(t, backup.BackupStatusSuccess, output.Status[0].Status)
	assert.Equal(t, extconfig.Get().Instance.Identifier, output.Status[0].Node)
}
