package roles

import "github.com/beryju/dns-dhcp-etcd-thingy/internal/storage"

type Event[T any] struct {
	Payload      T
	Sync         bool
	SourceRoleId string
}

type Role interface {
	Start(config []byte) error
	Stop()
	HandleEvent(ev *Event[any])
}

type Instance interface {
	GetKV() *storage.Client
	DispatchEvent(ev *Event[any])
}
