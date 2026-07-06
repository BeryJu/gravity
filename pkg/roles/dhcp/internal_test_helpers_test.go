package dhcp

import (
	"context"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/storage"
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

