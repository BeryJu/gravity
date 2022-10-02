package instance_test

import (
	"encoding/base64"
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	_, err := extconfig.Get().EtcdClient().Put(
		tests.Context(),
		"/foo",
		"bar",
	)
	assert.NoError(t, err)

	rootInst := instance.New()
	entries, err := rootInst.Export()
	assert.NoError(t, err)
	assert.True(t, len(entries) > 0)
}

func TestImport(t *testing.T) {
	rootInst := instance.New()
	entries := []instance.ExportEntry{
		{
			Key:   "foo",
			Value: base64.StdEncoding.EncodeToString([]byte("bar")),
		},
		{
			Key:   "foo",
			Value: "bar",
		},
	}
	err := rootInst.Import(entries)
	assert.NoError(t, err)
	res, err := extconfig.Get().EtcdClient().Get(
		tests.Context(),
		"foo",
	)
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(res.Kvs[0].Value))
}
