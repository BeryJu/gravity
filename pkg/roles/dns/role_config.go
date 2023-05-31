package dns

import (
	"context"

	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/storage"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *Role) decodeRoleConfig(raw []byte) *types.DNSRoleConfig {
	def := types.DNSRoleConfig{
		Port: 53,
	}
	if len(raw) < 1 {
		return &def
	}
	conf, err := storage.Parse(raw, &def)
	if err != nil {
		r.log.Warn("failed to decode role config", zap.Error(err))
	}
	return conf
}

type APIRoleConfigOutput struct {
	Config *types.DNSRoleConfig `json:"config" required:"true"`
}

func (r *Role) APIRoleConfigGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIRoleConfigOutput) error {
		output.Config = r.cfg
		return nil
	})
	u.SetName("dns.get_role_config")
	u.SetTitle("DNS role config")
	u.SetTags("roles/dns")
	return u
}

type APIRoleConfigInput struct {
	Config *types.DNSRoleConfig `json:"config" required:"true"`
}

func (r *Role) APIRoleConfigPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIRoleConfigInput, output *struct{}) error {
		jc, err := proto.Marshal(input.Config)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		_, err = r.i.KV().Put(
			ctx,
			r.i.KV().Key(
				instanceTypes.KeyInstance,
				instanceTypes.KeyRole,
				types.KeyRole,
			).String(),
			string(jc),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dns.put_role_config")
	u.SetTitle("DNS role config")
	u.SetTags("roles/dns")
	u.SetExpectedErrors(status.InvalidArgument, status.Internal)
	return u
}
