package etcd_test

import (
	"os"
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/etcd"
	"beryju.io/gravity/pkg/storage"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestEmbeddedEtcd_Start(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("etcd", ctx)
	etcdRole := etcd.New(inst)

	defer func() {
		etcdRole.Stop()
		_ = os.RemoveAll(etcdRole.Config().Dir)
	}()

	err := etcdRole.Start(ctx, []byte{})
	assert.NoError(t, err)

	c := storage.NewClient(
		"/gravity-test",
		extconfig.Get().Logger(),
		true,
		"localhost:2379",
	)
	m, err := c.MemberList(ctx)
	assert.NoError(t, err)
	assert.Len(t, m.Members, 1)
}
