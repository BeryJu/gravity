package dhcp

import (
	"context"
	"encoding/json"
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/storage"
	"github.com/getsentry/sentry-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type dhcpTestMigrator struct{}

func (dhcpTestMigrator) AddMigration(roles.Migration) {}

func (dhcpTestMigrator) Run(context.Context) (*storage.Client, error) {
	return extconfig.Get().EtcdClient(), nil
}

type dhcpTestInstance struct {
	ctx context.Context
	log *zap.Logger
	kv  *storage.Client
}

func newDHCPTestInstance(ctx context.Context) *dhcpTestInstance {
	return &dhcpTestInstance{
		ctx: ctx,
		log: extconfig.Get().Logger().Named("role.dhcp.test"),
		kv:  extconfig.Get().EtcdClient(),
	}
}

func (i *dhcpTestInstance) KV() *storage.Client {
	return i.kv
}

func (i *dhcpTestInstance) Log() *zap.Logger {
	return i.log
}

func (i *dhcpTestInstance) DispatchEvent(string, *roles.Event) {}

func (i *dhcpTestInstance) AddEventListener(string, roles.EventHandler) {}

func (i *dhcpTestInstance) Context() context.Context {
	return i.ctx
}

func (i *dhcpTestInstance) ExecuteHook(roles.HookOptions, ...interface{}) interface{} {
	return nil
}

func (i *dhcpTestInstance) Migrator() roles.RoleMigrator {
	return dhcpTestMigrator{}
}

func setupDHCPInternalTest(t testing.TB) context.Context {
	t.Helper()

	ctx, cancel := context.WithCancel(t.Context())
	tx := sentry.StartTransaction(ctx, "test")

	_, err := extconfig.Get().EtcdClient().Delete(tx.Context(), "/", clientv3.WithPrefix())
	if err != nil {
		t.Fatalf("failed to reset etcd: %v", err)
	}

	t.Cleanup(func() {
		tx.Finish()
		cancel()
	})

	return tx.Context()
}

func mustJSON(v any) string {
	raw, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(raw)
}

func panicIfError(args ...any) {
	for _, arg := range args {
		if err, ok := arg.(error); ok && err != nil {
			panic(err)
		}
	}
}
