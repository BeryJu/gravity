package roles

import (
	"beryju.io/ddet/pkg/storage"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Role interface {
	Start(config []byte) error
	Stop()
}

type Event struct {
	Payload EventPayload
	topic   string
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
	RelatedObjectKey     string
	RelatedObjectOptions []clientv3.OpOption
}

func NewEvent(data map[string]interface{}) *Event {
	return &Event{
		Payload: EventPayload{
			Data:                 data,
			RelatedObjectOptions: make([]clientv3.OpOption, 0),
		},
	}
}

type EventHandler func(ev *Event)

type Instance interface {
	GetKV() *storage.Client
	DispatchEvent(topic string, ev *Event)
	AddEventListener(topic string, handler EventHandler)
}
