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

func TestRole_Start_NoConfig(t *testing.T) {
	rootInst := instance.NewInstance()
	inst := rootInst.ForRole("backup")
	role := backup.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Equal(t, roles.ErrRoleNotConfigured, role.Start(ctx, []byte{}))
	role.Stop()
}

func TestRole_Start(t *testing.T) {
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
	role.Stop()
}

func TestRole_SaveBackup(t *testing.T) {
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
	status := role.SaveSnapshot()
	assert.Equal(t, nil, status.Error)
	assert.Equal(t, "success", status.Status)
}
