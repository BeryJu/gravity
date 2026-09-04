package instance

import (
	"context"
	"sync"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/migrate"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/storage"
	"go.uber.org/zap"
)

type RoleInstance struct {
	log    *zap.Logger
	parent *Instance
	roleId string

	// A role is constructed once and started again on every configuration
	// change, so the fields belonging to a single run are swapped in place
	// rather than by replacing the instance: roles keep a reference to the
	// instance they were constructed with, and replacing it would leave them
	// reading the state of a run that has already been stopped.
	m        sync.RWMutex
	kv       *storage.Client
	context  context.Context
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
	ri.m.RLock()
	defer ri.m.RUnlock()
	return ri.kv
}

func (ri *RoleInstance) setKV(kv *storage.Client) {
	ri.m.Lock()
	defer ri.m.Unlock()
	ri.kv = kv
}

func (ri *RoleInstance) Log() *zap.Logger {
	return ri.log
}

func (ri *RoleInstance) Context() context.Context {
	ri.m.RLock()
	defer ri.m.RUnlock()
	return ri.context
}

func (ri *RoleInstance) Migrator() roles.RoleMigrator {
	ri.m.RLock()
	defer ri.m.RUnlock()
	return ri.migrator
}

// rebind points the instance at the context for a fresh run of the role. The
// migrator is rebuilt along with it because roles register their migrations
// again on every start, and would otherwise accumulate duplicates.
func (ri *RoleInstance) rebind(ctx context.Context) {
	ri.m.Lock()
	ri.context = ctx
	ri.m.Unlock()
	// Built outside the lock: the migrator reads back from the instance.
	migrator := migrate.New(ri)
	ri.m.Lock()
	ri.migrator = migrator
	ri.m.Unlock()
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
