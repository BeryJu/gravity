package instance

import (
	"context"
	"encoding/json"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	"github.com/swaggest/rest/web"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type InstanceInfo struct {
	Version    string `json:"version"`
	Roles      string `json:"roles"`
	Identifier string `json:"identifier"`
	IP         string `json:"ip"`
}

func (i *Instance) getInfo() *InstanceInfo {
	return &InstanceInfo{
		Version:    extconfig.FullVersion(),
		Roles:      extconfig.Get().BootstrapRoles,
		Identifier: extconfig.Get().Instance.Identifier,
		IP:         extconfig.Get().Instance.IP,
	}
}

func (i *Instance) apiHandlerInstances() usecase.Interactor {
	type instancesOutput struct {
		Instances []InstanceInfo `json:"instances"`
	}
	u := usecase.NewIOI(new(struct{}), new(instancesOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*instancesOutput)
		)
		instances, err := i.kv.Get(
			ctx,
			i.kv.Key(types.KeyInstance).Prefix(true).String(),
			clientv3.WithPrefix(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, ri := range instances.Kvs {
			// We only want one level
			if strings.Contains(string(ri.Key), "/") {
				continue
			}
			var inst InstanceInfo
			err := json.Unmarshal(ri.Value, &inst)
			if err != nil {
				i.log.WithError(err).Warning("failed to parse instance info")
				continue
			}
			out.Instances = append(out.Instances, inst)
		}
		return nil
	})
	u.SetName("root.get_instances")
	u.SetTitle("Instances")
	u.SetTags("instances")
	u.SetExpectedErrors(status.Internal)
	return u
}

func (i *Instance) writeInstanceInfo() {
	i.ForRole("instance_info").AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/instances", i.apiHandlerInstances())
	})
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
