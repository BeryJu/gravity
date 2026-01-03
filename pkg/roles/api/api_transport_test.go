package api_test

import (
	"encoding/base64"
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)
	role := api.New(inst)
	tests.PanicIfError(role.Start(ctx, []byte{}))
	defer role.Stop()

	var output api.APIExportOutput

	_, err := extconfig.Get().EtcdClient().Put(
		ctx,
		"/foo",
		"bar",
	)
	assert.NoError(t, err)

	err = role.APIClusterExport().Interact(ctx, api.APIExportInput{
		Safe: true,
	}, &output)
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
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)
	role := api.New(inst)
	tests.PanicIfError(role.Start(ctx, []byte{}))
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

	err := role.APIClusterImport().Interact(ctx, entries, &struct{}{})
	assert.NoError(t, err)
	res, err := extconfig.Get().EtcdClient().Get(
		ctx,
		"foo",
	)
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(res.Kvs[0].Value))
}
