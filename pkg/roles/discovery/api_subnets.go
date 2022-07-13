package discovery

import (
	"context"

	"beryju.io/gravity/pkg/roles/discovery/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *DiscoveryRole) apiHandlerSubnets() usecase.Interactor {
	type subnet struct {
		CIDR         string `json:"cidr"`
		DiscoveryTTL int    `json:"discoveryTTL"`
	}
	type subnetsOutput struct {
		Subnets []subnet `json:"subnets"`
	}
	u := usecase.NewIOI(new(struct{}), new(subnetsOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*subnetsOutput)
		)
		prefix := r.i.KV().Key(types.KeyRole, types.KeySubnets).Prefix(true).String()
		subnets, err := r.i.KV().Get(ctx, prefix, clientv3.WithPrefix())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, rawSub := range subnets.Kvs {
			sub, err := r.subnetFromKV(rawSub)
			if err != nil {
				r.log.WithError(err).Warning("failed to parse subnet")
				continue
			}
			out.Subnets = append(out.Subnets, subnet{
				CIDR:         sub.CIDR,
				DiscoveryTTL: sub.DiscoveryTTL,
			})
		}
		return nil
	})
	u.SetTitle("Discovery subnets")
	u.SetTags("discovery")
	u.SetDescription("List all Discovery subnets.")
	u.SetExpectedErrors(status.Internal)
	return u
}
