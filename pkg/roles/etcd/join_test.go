package etcd_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/etcd"
	"beryju.io/gravity/pkg/tests"
)

func getRole() *etcd.Role {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("etcd", ctx)
	role := etcd.New(inst)

	return role
}

func TestClusterCanJoin(t *testing.T) {
	defer tests.Setup(t)()
	r := getRole()

	r.Stop()
}
