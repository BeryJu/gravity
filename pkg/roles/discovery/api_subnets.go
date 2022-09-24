package discovery

import (
	"context"

	"beryju.io/gravity/pkg/roles/discovery/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *Role) apiHandlerSubnets() usecase.Interactor {
	type subnet struct {
		Name string `json:"name" required:"true"`

		CIDR         string `json:"cidr" required:"true"`
		DiscoveryTTL int    `json:"discoveryTTL" required:"true"`
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
				Name:         sub.Identifier,
				CIDR:         sub.CIDR,
				DiscoveryTTL: sub.DiscoveryTTL,
			})
		}
		return nil
	})
	u.SetName("discovery.get_subnets")
	u.SetTitle("Discovery subnets")
	u.SetTags("roles/discovery")
	u.SetExpectedErrors(status.Internal)
	return u
}

func (r *Role) apiHandlerSubnetsPut() usecase.Interactor {
	type subnetsInput struct {
		Name string `query:"identifier" required:"true"`

		SubnetCIDR   string `json:"subnetCidr" required:"true"`
		DiscoveryTTL int    `json:"discoveryTTL" required:"true"`
	}
	u := usecase.NewIOI(new(subnetsInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*subnetsInput)
		)
		s := r.newSubnet(in.Name)
		s.CIDR = in.SubnetCIDR
		s.DiscoveryTTL = in.DiscoveryTTL
		err := s.put()
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("discovery.put_subnets")
	u.SetTitle("Discovery Subnets")
	u.SetTags("roles/discovery")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}

func (r *Role) apiHandlerSubnetsDelete() usecase.Interactor {
	type subnetsInput struct {
		Name string `query:"identifier"`
	}
	u := usecase.NewIOI(new(subnetsInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*subnetsInput)
		)
		_, err := r.i.KV().Delete(ctx, r.i.KV().Key(
			types.KeyRole,
			types.KeySubnets,
			in.Name,
		).String())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("discovery.delete_subnets")
	u.SetTitle("Discovery Subnets")
	u.SetTags("roles/discovery")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}
