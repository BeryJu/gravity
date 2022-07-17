package roles

import (
	"context"

	"beryju.io/gravity/pkg/storage"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Role interface {
	Start(ctx context.Context, config []byte) error
	Stop()
}

type Event struct {
	Payload EventPayload
	Context context.Context
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

type Instance interface {
	KV() *storage.Client
	Log() *log.Entry
	DispatchEvent(topic string, ev *Event)
	AddEventListener(topic string, handler EventHandler)
}
