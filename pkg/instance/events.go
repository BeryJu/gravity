package instance

import "beryju.io/gravity/pkg/roles"

func (i *Instance) DispatchEvent(topic string, ev *roles.Event) {
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
