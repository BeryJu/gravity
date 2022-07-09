package instance

import "beryju.io/ddet/pkg/roles"

func (i *Instance) dispatchEvent(topic string, ev *roles.Event) {
	i.eventHandlersM.RLock()
	defer i.eventHandlersM.RUnlock()
	handlers, ok := i.eventHandlers[topic]
	if !ok {
		return
	}
	for _, handlers := range handlers {
		for _, handler := range handlers {
			handler(ev)
		}
	}
}
