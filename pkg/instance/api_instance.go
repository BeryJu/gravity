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

func (i *Instance) setupInstanceAPI() {
	i.ForRole("instance").AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/instances", i.APIInstances())
		svc.Get("/api/v1/info", i.APIInstanceInfo())
	})
}

type APIInstancesOutput struct {
	Instances []InstanceInfo `json:"instances" required:"true"`
}

func (i *Instance) APIInstances() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIInstancesOutput) error {
		prefix := i.kv.Key(types.KeyInstance).Prefix(true).String()
		instances, err := i.kv.Get(
			ctx,
			prefix,
			clientv3.WithPrefix(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, ri := range instances.Kvs {
			// We only want one level
			relKey := strings.TrimPrefix(string(ri.Key), prefix)
			if strings.Contains(relKey, "/") {
				continue
			}
			var inst InstanceInfo
			err := json.Unmarshal(ri.Value, &inst)
			if err != nil {
				i.log.WithError(err).Warning("failed to parse instance info")
				continue
			}
			output.Instances = append(output.Instances, inst)
		}
		return nil
	})
	u.SetName("root.get_instances")
	u.SetTitle("Instances")
	u.SetTags("instances")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APIInstanceInfo struct {
	Version   string `json:"version" required:"true"`
	BuildHash string `json:"buildHash" required:"true"`

	Dirs *extconfig.ExtConfigDirs `json:"dirs" required:"true"`

	CurrentInstanceIdentifier string `json:"currentInstanceIdentifier" required:"true"`
	CurrentInstanceIP         string `json:"currentInstanceIP" required:"true"`
}

func (i *Instance) APIInstanceInfo() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIInstanceInfo) error {
		output.Version = extconfig.Version
		output.BuildHash = extconfig.BuildHash
		output.Dirs = extconfig.Get().Dirs()
		output.CurrentInstanceIP = extconfig.Get().Instance.IP
		output.CurrentInstanceIdentifier = extconfig.Get().Instance.Identifier
		return nil
	})
	u.SetName("root.get_info")
	u.SetTitle("Instances")
	u.SetTags("instances")
	u.SetExpectedErrors(status.Internal)
	return u
}
