package instance

import (
	"encoding/json"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type InstanceInfo struct {
	Version    string `json:"version" required:"true"`
	Roles      string `json:"roles" required:"true"`
	Identifier string `json:"identifier" required:"true"`
	IP         string `json:"ip" required:"true"`
}

func (i *Instance) getInfo() *InstanceInfo {
	return &InstanceInfo{
		Version:    extconfig.FullVersion(),
		Roles:      extconfig.Get().BootstrapRoles,
		Identifier: extconfig.Get().Instance.Identifier,
		IP:         extconfig.Get().Instance.IP,
	}
}

func (i *Instance) writeInstanceInfo() {
	ji, err := json.Marshal(i.getInfo())
	if err != nil {
		i.log.WithError(err).Warning("failed to get instance info")
		return
	}
	lease, err := i.kv.Lease.Grant(i.rootContext, 100)
	if err != nil {
		i.log.WithError(err).Warning("failed to grant lease")
		return
	}
	i.kv.Put(
		i.rootContext,
		i.kv.Key(
			types.KeyInstance,
			extconfig.Get().Instance.Identifier,
		).String(),
		string(ji),
		clientv3.WithLease(lease.ID),
	)
	keepAlive, err := i.kv.KeepAlive(i.rootContext, lease.ID)
	if err != nil {
		i.log.WithError(err).Warning("failed to grant lease")
		return
	}
	go func() {
		for range keepAlive {
			// eat messages until keep alive channel closes
		}
	}()
}
