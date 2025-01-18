package roles

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/storage"
	"github.com/Masterminds/semver/v3"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var ErrRoleNotConfigured = errors.New("role: Not configured")

type Role interface {
	Start(ctx context.Context, config []byte) error
	Stop()
}

type MigratableRole interface {
	Role
	RegisterMigrations()
}

type Event struct {
	Context context.Context
	topic   string
	Payload EventPayload
}

func (ev *Event) WithTopic(topic string) *Event {
	ev.topic = topic
	return ev
}

func (ev *Event) String() string {
	return ev.topic
}

type EventPayload struct {
	Data                 map[string]interface{}
	RelatedObjectKey     *storage.Key
	RelatedObjectOptions []clientv3.OpOption
}

func NewEvent(ctx context.Context, data map[string]interface{}) *Event {
	return &Event{
		Context: ctx,
		Payload: EventPayload{
			Data:                 data,
			RelatedObjectOptions: make([]clientv3.OpOption, 0),
		},
	}
}

type EventHandler func(ev *Event)

type HookOptions struct {
	Method string
	Source string
	Env    map[string]interface{}
}

type Migration interface {
	Check(clusterVersion *semver.Version, ctx context.Context) (bool, error)
	Hook(context.Context) (*storage.Client, error)
	Cleanup(context.Context) error
	Name() string
}

type RoleMigrator interface {
	AddMigration(Migration)
	Run(ctx context.Context) (*storage.Client, error)
}

type Instance interface {
	KV() *storage.Client
	Log() *zap.Logger
	DispatchEvent(topic string, ev *Event)
	AddEventListener(topic string, handler EventHandler)
	Context() context.Context
	ExecuteHook(HookOptions, ...interface{}) interface{}
	Migrator() RoleMigrator
}

type RoleConstructor func(Instance) Role

var roleRegistry map[string]RoleConstructor = make(map[string]RoleConstructor)

func Register(name string, constructor RoleConstructor) {
	roleRegistry[name] = constructor
}

func GetRole(name string) RoleConstructor {
	return roleRegistry[name]
}
