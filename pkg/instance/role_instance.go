package instance

import (
	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/storage"
	log "github.com/sirupsen/logrus"
)

type RoleInstance struct {
	kv     *storage.Client
	log    *log.Entry
	roleId string
	parent *Instance
}

func (i *Instance) ForRole(roleId string) *RoleInstance {
	in := &RoleInstance{
		log:    log.WithField("forRole", roleId),
		kv:     i.kv,
		roleId: roleId,
		parent: i,
	}
	return in
}

func (ri *RoleInstance) GetKV() *storage.Client {
	return ri.kv
}

func (ri *RoleInstance) DispatchEvent(topic string, ev *roles.Event) {
	ri.parent.dispatchEvent(topic, ev.WithTopic(topic))
}

func (ri *RoleInstance) AddEventListener(topic string, handler roles.EventHandler) {
	ri.parent.eventHandlersM.Lock()
	defer ri.parent.eventHandlersM.Unlock()
	topicHandlers, ok := ri.parent.eventHandlers[topic]
	if !ok {
		topicHandlers = make(map[string][]roles.EventHandler)
	}
	roleHandlers, ok := topicHandlers[ri.roleId]
	if !ok {
		roleHandlers = make([]roles.EventHandler, 0)
	}
	roleHandlers = append(roleHandlers, handler)
	topicHandlers[ri.roleId] = roleHandlers
	ri.parent.eventHandlers[topic] = topicHandlers
}
