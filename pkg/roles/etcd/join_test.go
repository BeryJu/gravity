package etcd

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/tests"
)

func getRole() *Role {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("etcd", ctx)
	role := New(inst)

	return role
}

func TestClusterCanJoin(t *testing.T) {

}
