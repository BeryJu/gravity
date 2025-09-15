package instance

import (
	"context"
	"encoding/json"
	"slices"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/types"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type InstanceInfo struct {
	Version    string   `json:"version" required:"true"`
	Roles      []string `json:"roles" required:"true"`
	Identifier string   `json:"identifier" required:"true"`
	IP         string   `json:"ip" required:"true"`
}

func (i *Instance) getInfo() *InstanceInfo {
	roles := maps.Keys(i.roles)
	if i.etcd != nil {
		roles = append(roles, "etcd")
	}
	slices.Sort(roles)
	return &InstanceInfo{
		Version:    extconfig.FullVersion(),
		Roles:      roles,
		Identifier: extconfig.Get().Instance.Identifier,
		IP:         extconfig.Get().Instance.IP,
	}
}

func (i *Instance) keepAliveInstanceInfo(ctx context.Context) {
	restarter := func() {
		if i.instanceSession != nil {
			<-i.instanceSession.Done()
		}
		i.keepAliveInstanceInfo(ctx)
	}
	defer func() {
		i.putInstanceInfo(ctx)
		go restarter()
	}()
	sess, err := concurrency.NewSession(i.kv.Client)
	if err != nil {
		i.log.Warn("failed to setup etcd lease session", zap.Error(err))
		return
	}
	i.instanceSession = sess
}

func (i *Instance) putInstanceInfo(ctx context.Context, opts ...clientv3.OpOption) {
	ji, err := json.Marshal(i.getInfo())
	if err != nil {
		i.log.Warn("failed to get instance info", zap.Error(err))
		return
	}
	if i.instanceSession != nil {
		opts = append(opts, clientv3.WithLease(i.instanceSession.Lease()))
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
