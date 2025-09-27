package instance

import (
	"context"
	"encoding/json"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/migrate"
	"beryju.io/gravity/pkg/instance/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type APIClusterInfoOutput struct {
	ClusterVersion      string         `json:"clusterVersion" required:"true"`
	ClusterVersionShort string         `json:"clusterVersionShort" required:"true"`
	Instances           []InstanceInfo `json:"instances" required:"true"`
}

func (i *Instance) APIClusterInfo() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIClusterInfoOutput) error {
		ri := i.ForRole("api", ctx)
		m := migrate.New(ri)
		cv, err := m.GetClusterVersion(ctx)
		if err != nil {
			return status.Internal
		}
		output.ClusterVersion = cv.String()
		sv, _ := cv.SetMetadata("")
		output.ClusterVersionShort = sv.String()

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
				i.log.Warn("failed to parse instance info", zap.Error(err))
				continue
			}
			output.Instances = append(output.Instances, inst)
		}
		return nil
	})
	u.SetName("cluster.get_cluster_info")
	u.SetTitle("Cluster")
	u.SetTags("cluster")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APIInstanceInfo struct {
	Version   string `json:"version" required:"true"`
	BuildHash string `json:"buildHash" required:"true"`

	Dirs *extconfig.ExtConfigDirs `json:"dirs" required:"true"`

	InstanceIdentifier string `json:"instanceIdentifier" required:"true"`
	InstanceIP         string `json:"instanceIP" required:"true"`
}

func (i *Instance) APIInstanceGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIInstanceInfo) error {
		output.Version = extconfig.Version
		output.BuildHash = extconfig.BuildHash
		output.Dirs = extconfig.Get().Dirs()
		output.InstanceIP = extconfig.Get().Instance.IP
		output.InstanceIdentifier = extconfig.Get().Instance.Identifier
		return nil
	})
	u.SetName("cluster.get_instance_info")
	u.SetTitle("Instance")
	u.SetTags("cluster/instances")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APIInstancesPutInput struct {
	Identifier string `query:"identifier" required:"true"`

	Roles []string `json:"roles" required:"true"`
}

func (i *Instance) APIInstancePut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIInstancesPutInput, output *struct{}) error {
		_, err := i.kv.Put(
			ctx,
			i.kv.Key(
				types.KeyInstance,
				i.identifier,
				types.KeyRoles,
			).String(),
			string(strings.Join(input.Roles, types.RoleSeparator)),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("cluster.put_instance")
	u.SetTitle("Instance")
	u.SetTags("cluster/instances")
	u.SetExpectedErrors(status.Internal)
	return u
}
