package backup_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/backup"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

const (
	TestingMinioAccessKey = "gravity"
	TestingMinioSecretKey = "gravity-key"
)

func getRole() *backup.Role {
	rootInst := instance.New()
	inst := rootInst.ForRole("backup")
	role := backup.New(inst)
	ctx := tests.Context()

	cfg := tests.MustJSON(&backup.RoleConfig{
		Endpoint:  "http://localhost:9000",
		AccessKey: TestingMinioAccessKey,
		SecretKey: TestingMinioSecretKey,
		Bucket:    "gravity",
		Path:      "foo",
		CronExpr:  "",
	})
	role.Start(ctx, []byte(cfg))
	return role
}

func TestRoleStartNoConfig(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("backup")
	role := backup.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Equal(t, roles.ErrRoleNotConfigured, role.Start(ctx, []byte{}))
	role.Stop()
}

func TestRoleStart(t *testing.T) {
	assert.NotNil(t, getRole())
}

func TestSaveBackup(t *testing.T) {
	rootInst := instance.New()
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
	status := role.SaveSnapshot()
	assert.Equal(t, "", status.Error)
	assert.Equal(t, "success", status.Status)
}
