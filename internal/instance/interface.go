package instance

import (
	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles"
	"github.com/beryju/dns-dhcp-etcd-thingy/internal/storage"
)

func (i *Instance) GetKV() *storage.Client {
	return i.kv
}

func (i *Instance) DispatchEvent(ev *roles.Event[any]) {
}
