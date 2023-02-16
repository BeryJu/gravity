package instance

import (
	"context"
	"encoding/json"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/types"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type InstanceInfo struct {
	Version    string `json:"version" required:"true"`
	Roles      string `json:"roles" required:"true"`
	Identifier string `json:"identifier" required:"true"`
	IP         string `json:"ip" required:"true"`
}

func (i *Instance) getInfo() *InstanceInfo {
	roles := maps.Keys(i.roles)
	if i.etcd != nil {
		roles = append(roles, "etcd")
	}
	slices.Sort(roles)
	return &InstanceInfo{
		Version:    extconfig.FullVersion(),
		Roles:      strings.Join(roles, ";"),
		Identifier: extconfig.Get().Instance.Identifier,
		IP:         extconfig.Get().Instance.IP,
	}
}

func (i *Instance) keepAliveInstanceInfo(ctx context.Context) {
	if i.instanceInfoLease == nil {
		lease, err := i.kv.Lease.Grant(ctx, 100)
		if err != nil {
			i.log.Warn("failed to grant lease", zap.Error(err))
			return
		}
		i.instanceInfoLease = &lease.ID
	}
	keepAlive, err := i.kv.KeepAlive(ctx, *i.instanceInfoLease)
	if err != nil {
		i.log.Warn("failed to grant lease", zap.Error(err))
		return
	}
	go func() {
		for range keepAlive {
			// eat messages until keep alive channel closes
		}
	}()
}

func (i *Instance) putInstanceInfo(ctx context.Context) {
	ji, err := json.Marshal(i.getInfo())
	if err != nil {
		i.log.Warn("failed to get instance info", zap.Error(err))
		return
	}
	opts := []clientv3.OpOption{}
	if i.instanceInfoLease != nil {
		opts = append(opts, clientv3.WithLease(*i.instanceInfoLease))
	}
	_, err = i.kv.Put(
		ctx,
		i.kv.Key(
			types.KeyInstance,
			extconfig.Get().Instance.Identifier,
		).String(),
		string(ji),
		opts...,
	)
	if err != nil {
		i.log.Warn("failed to put instance info", zap.Error(err))
	}
}
