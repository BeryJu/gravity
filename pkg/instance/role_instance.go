package instance

import (
	"context"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/storage"
	log "github.com/sirupsen/logrus"
)

type RoleInstance struct {
	log    *log.Entry
	roleId string
	parent *Instance
}

func (i *Instance) ForRole(roleId string) *RoleInstance {
	in := &RoleInstance{
		log:    i.log.WithField("forRole", roleId),
		roleId: roleId,
		parent: i,
	}
	return in
}

func (ri *RoleInstance) KV() *storage.Client {
	return ri.parent.kv
}

func (ri *RoleInstance) Log() *log.Entry {
	return ri.log
}

func (ri *RoleInstance) DispatchEvent(topic string, ev *roles.Event) {
	ri.log.WithField("topic", topic).Debug("dispatching event")
	if ev.Context == nil {
		ev.Context = context.Background()
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
