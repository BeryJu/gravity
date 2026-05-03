package tsdb_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/tsdb"
	"beryju.io/gravity/pkg/roles/tsdb/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestAPIMetrics_Memory(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("metrics", ctx)
	tests.PanicIfError(inst.KV().Delete(
		ctx,
		inst.KV().Key(
			types.KeyRole,
		).Prefix(true).String(),
		clientv3.WithPrefix(),
	))

	role := tsdb.New(inst)
	assert.NoError(t, role.Start(ctx, []byte{}))
	inst.DispatchEvent(types.EventTopicTSDBWrite, roles.NewEvent(ctx, map[string]interface{}{}))

	var output types.APIMetricsGetOutput
	assert.NoError(t, role.APIMetrics().Interact(ctx, types.APIMetricsGetInput{
		Role:     "system",
		Category: "memory",
	}, &output))
	assert.Equal(t, extconfig.Get().Instance.Identifier, output.Records[0].Node)
	assert.Equal(t, 1, len(output.Records))
}
