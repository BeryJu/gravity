package instance

import "beryju.io/ddet/pkg/roles"

func (i *Instance) dispatchEvent(topic string, ev *roles.Event) {
	i.eventHandlersM.RLock()
	defer i.eventHandlersM.RUnlock()
	handlers, ok := i.eventHandlers[topic]
	if !ok {
		return
	}
	for role, handlers := range handlers {
		i.log.WithField("topic", topic).WithField("role", role).Trace("dispatching event")
		for _, handler := range handlers {
			handler(ev)
		}
	}
}
