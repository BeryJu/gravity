package api_test

import (
	"encoding/base64"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("api")
	role := api.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output api.APIExportOutput

	extconfig.Get().EtcdClient().Delete(
		ctx,
		"/",
		clientv3.WithPrefix(),
	)
	_, err := extconfig.Get().EtcdClient().Put(
		tests.Context(),
		"/foo",
		"bar",
	)
	assert.NoError(t, err)

	err = role.APIClusterExport().Interact(ctx, struct{}{}, &output)
	assert.NoError(t, err)
	assert.Equal(t, api.APIExportOutput{
		Entries: []api.APITransportEntry{
			{
				Key:   "/foo",
				Value: "YmFy",
			},
		},
	}, output)
}

func TestImport(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("api")
	role := api.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	entries := api.APIImportInput{
		Entries: []api.APITransportEntry{
			{
				Key:   "foo",
				Value: base64.StdEncoding.EncodeToString([]byte("bar")),
			},
			{
				Key:   "foo",
				Value: "bar",
			},
		},
	}
	// var output struct{}

	err := role.APIClusterImport().Interact(ctx, entries, &struct{}{})
	assert.NoError(t, err)
	res, err := extconfig.Get().EtcdClient().Get(
		tests.Context(),
		"foo",
	)
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(res.Kvs[0].Value))
}
