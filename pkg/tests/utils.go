package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"runtime"
	"strings"
	"testing"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/storage"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	testSpan          *sentry.Span
	testContextCancel context.CancelFunc
)

func PanicIfError(args ...interface{}) {
	for _, arg := range args {
		if e, ok := arg.(error); ok && e != nil {
			panic(arg)
		}
	}
}

func MustJSON(in interface{}) string {
	j, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func MustProto(in protoreflect.ProtoMessage) []byte {
	r, err := proto.Marshal(in)
	if err != nil {
		panic(err)
	}
	return r
}

func MustParseNetIP(t *testing.T, r string) netip.Addr {
	i, err := netip.ParseAddr(r)
	assert.NoError(t, err)
	return i
}

func Context() context.Context {
	return testSpan.Context()
}

func RandomString(prefix ...string) string {
	str := append(prefix, uuid.New().String())
	return strings.Join(str, "-")
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
		if rb, ok := res.([]byte); ok {
			assert.Equal(t, rb, values.Kvs[idx].Value)
		} else if rb, ok := res.(string); ok {
			assert.Equal(t, rb, string(values.Kvs[idx].Value))
		} else if rb, ok := res.(protoreflect.ProtoMessage); ok {
			assert.Equal(t, MustProto(rb), values.Kvs[idx].Value)
		} else {
			assert.Equal(t, MustJSON(res), string(values.Kvs[idx].Value))
		}
	}
}

func ResetEtcd(t *testing.T) {
	ctx := Context()
	_, err := extconfig.Get().EtcdClient().Delete(
		ctx,
		"/",
		clientv3.WithPrefix(),
	)
	assert.NoError(t, err)
}

func Setup(t *testing.T) func() {
	ctx, cn := context.WithCancel(context.Background())
	testSpan = sentry.StartTransaction(ctx, "test")
	testContextCancel = cn
	ResetEtcd(t)
	return func() {
		testContextCancel()
	}
}

func HasLocalDocker() bool {
	return runtime.GOOS == "linux"
}

func Listen(port int32) string {
	if runtime.GOOS == "darwin" {
		return fmt.Sprintf(":%d", port)
	}
	return extconfig.Get().Listen(port)
}

func WaitForPort(port int32) {
	max := 30
	try := 0
	listen := Listen(port)
	time.Sleep(500 * time.Millisecond)
	for {
		ln, err := net.Listen("tcp", listen)
		if ln != nil {
			_ = ln.Close()
		}
		if err != nil {
			return
		}
		try += 1
		if try >= max {
			panic(fmt.Errorf("failed to wait for port '%s' to be listening", listen))
		}
		time.Sleep(1 * time.Millisecond)
	}
}
