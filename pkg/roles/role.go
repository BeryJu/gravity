package roles

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/storage"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var ErrRoleNotConfigured = errors.New("role: Not configured")

type Role interface {
	Start(ctx context.Context, config []byte) error
	Stop()
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

type Instance interface {
	KV() *storage.Client
	Log() *zap.Logger
	DispatchEvent(topic string, ev *Event)
	AddEventListener(topic string, handler EventHandler)
	Context() context.Context
}
