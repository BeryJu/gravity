package instance

import (
	"context"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/migrate"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/storage"
	"go.uber.org/zap"
)

type RoleInstance struct {
	kv       *storage.Client
	context  context.Context
	log      *zap.Logger
	parent   *Instance
	roleId   string
	migrator *migrate.Migrator
}

func (i *Instance) ForRole(roleId string, ctx context.Context) *RoleInstance {
	ri := &RoleInstance{
		log: extconfig.Get().Logger().Named("role." + roleId).WithOptions(
			extconfig.SetLevel(extconfig.Get().LogLevelFor(roleId)),
		),
		roleId:  roleId,
		parent:  i,
		context: ctx,
		kv:      i.kv,
	}
	ri.migrator = migrate.New(ri)
	return ri
}

func (ri *RoleInstance) KV() *storage.Client {
	return ri.kv
}

func (ri *RoleInstance) Log() *zap.Logger {
	return ri.log
}

func (ri *RoleInstance) Context() context.Context {
	return ri.context
}

func (ri *RoleInstance) Migrator() roles.RoleMigrator {
	return ri.migrator
}

func (ri *RoleInstance) DispatchEvent(topic string, ev *roles.Event) {
	l := ri.log
	if extconfig.Get().Debug {
		l = l.With(zap.Any("payload", ev.Payload.Data))
	}
	l.Debug("dispatching event", zap.String("topic", topic))
	if ev.Context == nil {
		ev.Context = context.TODO()
	}
	ri.parent.DispatchEvent(topic, ev.WithTopic(topic))
}

func (ri *RoleInstance) AddEventListener(topic string, handler roles.EventHandler) {
	ri.parent.eventHandlersM.RLock()
	topicHandlers, ok := ri.parent.eventHandlers[topic]
	ri.parent.eventHandlersM.RUnlock()
	if !ok {
		topicHandlers = make(map[string][]roles.EventHandler)
	}
	roleHandlers, ok := topicHandlers[ri.roleId]
	if !ok {
		roleHandlers = make([]roles.EventHandler, 0)
	}
	roleHandlers = append(roleHandlers, handler)
	topicHandlers[ri.roleId] = roleHandlers
	ri.parent.eventHandlersM.Lock()
	defer ri.parent.eventHandlersM.Unlock()
	ri.parent.eventHandlers[topic] = topicHandlers
}
