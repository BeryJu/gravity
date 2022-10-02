package tests

import (
	"context"
	"encoding/json"
	"runtime"
	"testing"
	"time"

	"beryju.io/gravity/pkg/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func MustJSON(in interface{}) string {
	j, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func Context() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second) //nolint
	return ctx
}

func RandomString() string {
	return uuid.New().String()
}

func AssertEtcd(t *testing.T, c *storage.Client, key *storage.Key, expected ...interface{}) {
	values, err := c.Get(
		Context(),
		key.String(),
	)
	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(values.Kvs))
	for idx, res := range expected {
		assert.Equal(t, MustJSON(res), string(values.Kvs[idx].Value))
	}
}

func HasLocalDocker() bool {
	return runtime.GOOS == "linux"
}
