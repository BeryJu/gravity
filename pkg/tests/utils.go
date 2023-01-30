package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"testing"
	"time"

	"beryju.io/gravity/pkg/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
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
	args := []clientv3.OpOption{}
	if key.IsPrefix() {
		args = append(args, clientv3.WithPrefix())
	}
	values, err := c.Get(
		Context(),
		key.String(),
		args...,
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

func WaitForPort(listen string) {
	max := 30
	try := 0
	for {
		ln, err := net.Listen("tcp", listen)
		if ln != nil {
			_ = ln.Close()
		}
		if err == nil {
			return
		}
		try += 1
		if try >= max {
			panic(fmt.Errorf("failed to wait for port '%s' to be listening", listen))
		}
		time.Sleep(1 * time.Millisecond)
	}
}
