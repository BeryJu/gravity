package instance_test

import (
	"encoding/base64"
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func Test_Export(t *testing.T) {
	_, err := extconfig.Get().EtcdClient().Put(
		tests.Context(),
		"foo",
		"bar",
	)
	assert.NoError(t, err)

	rootInst := instance.New()
	_, err = rootInst.Export()
	assert.NoError(t, err)
}

func Test_Import(t *testing.T) {
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
