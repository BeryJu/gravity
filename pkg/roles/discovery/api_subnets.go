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

		CIDR         string `json:"subnetCidr" required:"true"`
		DiscoveryTTL int    `json:"discoveryTTL" required:"true"`
	}
	type subnetsOutput struct {
		Subnets []subnet `json:"subnets"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *subnetsOutput) error {
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
			output.Subnets = append(output.Subnets, subnet{
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
		Name string `query:"identifier" required:"true" maxLength:"255"`

		SubnetCIDR   string `json:"subnetCidr" required:"true" maxLength:"40"`
		DiscoveryTTL int    `json:"discoveryTTL" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input subnetsInput, output *struct{}) error {
		s := r.NewSubnet(input.Name)
		s.CIDR = input.SubnetCIDR
		s.DiscoveryTTL = input.DiscoveryTTL
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

func (r *Role) apiHandlerSubnetsStart() usecase.Interactor {
	type subnetsInput struct {
		Name string `query:"identifier" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input subnetsInput, output *struct{}) error {
		rawSub, err := r.i.KV().Get(ctx, r.i.KV().Key(
			types.KeyRole,
			types.KeySubnets,
			input.Name,
		).String())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		if len(rawSub.Kvs) < 1 {
			return status.Wrap(err, status.NotFound)
		}
		s, err := r.subnetFromKV(rawSub.Kvs[0])
		if err != nil {
			r.log.WithError(err).Warning("failed to parse subnet from KV")
			return status.Wrap(err, status.Internal)
		}
		go s.RunDiscovery()
		return nil
	})
	u.SetName("discovery.subnet_start")
	u.SetTitle("Discovery Subnets")
	u.SetTags("roles/discovery")
	u.SetExpectedErrors(status.Internal, status.NotFound)
	return u
}

func (r *Role) apiHandlerSubnetsDelete() usecase.Interactor {
	type subnetsInput struct {
		Name string `query:"identifier"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input subnetsInput, output *struct{}) error {
		_, err := r.i.KV().Delete(ctx, r.i.KV().Key(
			types.KeyRole,
			types.KeySubnets,
			input.Name,
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
