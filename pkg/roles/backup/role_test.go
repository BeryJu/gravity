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
	ctx := tests.Context()
	inst := rootInst.ForRole("backup", ctx)
	role := backup.New(inst)

	cfg := tests.MustJSON(&backup.RoleConfig{
		Endpoint:  "http://localhost:9001",
		AccessKey: TestingMinioAccessKey,
		SecretKey: TestingMinioSecretKey,
		Bucket:    "gravity",
		Path:      "foo",
		CronExpr:  "",
	})
	tests.PanicIfError(role.Start(ctx, []byte(cfg)))
	return role
}

func TestRoleStartNoConfig(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("backup", ctx)
	role := backup.New(inst)
	assert.NotNil(t, role)
	assert.Equal(t, roles.ErrRoleNotConfigured, role.Start(ctx, []byte{}))
	role.Stop()
}

func TestRoleStart(t *testing.T) {
	tests.Setup(t)
	assert.NotNil(t, getRole())
}

func TestSaveBackup(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("backup", ctx)
	role := backup.New(inst)
	assert.NotNil(t, role)

	cfg := tests.MustJSON(&backup.RoleConfig{
		Endpoint:  "http://localhost:9001",
		AccessKey: TestingMinioAccessKey,
		SecretKey: TestingMinioSecretKey,
		Bucket:    "gravity",
		Path:      "foo",
		CronExpr:  "",
	})
	assert.Nil(t, role.Start(ctx, []byte(cfg)))
	defer role.Stop()
	status := role.SaveSnapshot(ctx)
	assert.Equal(t, "", status.Error)
	assert.Equal(t, "success", status.Status)
}
