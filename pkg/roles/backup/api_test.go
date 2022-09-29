package backup_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/backup"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func Test_API_apiHandlerBackupStart(t *testing.T) {
	rootInst := instance.NewInstance()
	inst := rootInst.ForRole("backup")
	role := backup.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()

	cfg := tests.MustJSON(&backup.RoleConfig{
		Endpoint:  "http://localhost:9000",
		AccessKey: TestingMinioAccessKey,
		SecretKey: TestingMinioSecretKey,
		Bucket:    "gravity",
		Path:      "foo",
		CronExpr:  "",
	})
	assert.Nil(t, role.Start(ctx, []byte(cfg)))
	defer role.Stop()

	var output backup.BackupStatus
	role.APIHandlerBackupStart().Interact(tests.Context(), backup.BackupStartInput{
		Wait: true,
	}, &output)
	assert.NotNil(t, output)
}
