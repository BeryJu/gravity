package dhcp

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type APILeasesGetInput struct {
	ScopeName  string `query:"scope"`
	Identifier string `query:"identifier" description:"Optional identifier of a lease to get"`
}
type APILeaseInfo struct {
	Vendor string `json:"vendor"`
}
type APILease struct {
	Info             *APILeaseInfo `json:"info"`
	Identifier       string        `json:"identifier" required:"true"`
	Address          string        `json:"address" required:"true"`
	Hostname         string        `json:"hostname" required:"true"`
	AddressLeaseTime string        `json:"addressLeaseTime" required:"true"`
	ScopeKey         string        `json:"scopeKey" required:"true"`
	DNSZone          string        `json:"dnsZone"`
	Expiry           int64         `json:"expiry"`
	Description      string        `json:"description" required:"true"`
}
type APILeasesGetOutput struct {
	Leases []*APILease `json:"leases" required:"true"`
}

func (r *Role) APILeasesGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APILeasesGetInput, output *APILeasesGetOutput) error {
		// Validate that the scope name exists
		rawScope, err := r.i.KV().Get(
			ctx,
			r.i.KV().Key(
				types.KeyRole,
				types.KeyScopes,
				input.ScopeName,
			).String(),
		)
		if err != nil || len(rawScope.Kvs) < 1 {
			r.log.Warn("failed to get scope", zap.Error(err))
			return status.Wrap(errors.New("failed to get scope"), status.Internal)
		}

		leaseKey := r.i.KV().Key(
			types.KeyRole,
			types.KeyLeases,
		)
		if input.Identifier == "" {
			leaseKey = leaseKey.Prefix(true)
		} else {
			leaseKey = leaseKey.Add(input.Identifier)
		}
		rawLeases, err := r.i.KV().Get(ctx, leaseKey.String(), clientv3.WithPrefix())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, lease := range rawLeases.Kvs {
			l, err := r.leaseFromKV(lease)
			if err != nil {
				r.log.Warn("failed to parse lease", zap.Error(err))
				continue
			}
			if l.ScopeKey != input.ScopeName {
				continue
			}
			al := &APILease{
				Identifier:       l.Identifier,
				Address:          l.Address,
				Hostname:         l.Hostname,
				AddressLeaseTime: l.AddressLeaseTime,
				ScopeKey:         l.ScopeKey,
				DNSZone:          l.DNSZone,
				Expiry:           l.Expiry,
				Description:      l.Description,
			}
			if r.oui != nil {
				entry, err := r.oui.LookupString(l.Identifier)
				if err == nil {
					al.Info = &APILeaseInfo{
						Vendor: entry.Organization,
					}
				}
			}
			output.Leases = append(output.Leases, al)
		}
		return nil
	})
	u.SetName("dhcp.get_leases")
	u.SetTitle("DHCP Leases")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APILeasesPutInput struct {
	Identifier string `query:"identifier" required:"true" maxLength:"255"`
	Scope      string `query:"scope" required:"true" maxLength:"255"`

	Address          string `json:"address" required:"true" maxLength:"40"`
	Hostname         string `json:"hostname" required:"true" maxLength:"255"`
	AddressLeaseTime string `json:"addressLeaseTime" required:"true" maxLength:"40"`
	DNSZone          string `json:"dnsZone" maxLength:"255"`
	Expiry           int64  `json:"expiry"`
	Description      string `json:"description"`
}

func (r *Role) APILeasesPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APILeasesPutInput, output *struct{}) error {
		rawScope, err := r.i.KV().Get(
			ctx,
			r.i.KV().Key(
				types.KeyRole,
				types.KeyScopes,
				input.Scope,
			).String(),
		)
		if err != nil || len(rawScope.Kvs) < 1 {
			r.log.Warn("failed to get scope", zap.Error(err))
			return status.Wrap(errors.New("failed to get scope"), status.Internal)
		}
		scope, err := r.scopeFromKV(rawScope.Kvs[0])
		if err != nil {
			r.log.Warn("failed to construct scope", zap.Error(err))
			return status.Wrap(errors.New("failed to construct scope"), status.Internal)
		}

		l := r.NewLease(input.Identifier)
		l.Address = input.Address
		l.Hostname = input.Hostname
		l.AddressLeaseTime = input.AddressLeaseTime
		l.ScopeKey = input.Scope
		l.DNSZone = input.DNSZone
		l.Expiry = input.Expiry
		l.Description = input.Description

		l.scope = scope
		err = l.Put(ctx, -1)
		if err != nil {
			r.log.Warn("failed to put lease", zap.Error(err))
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dhcp.put_leases")
	u.SetTitle("DHCP Leases")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}

type APILeasesWOLInput struct {
	Identifier string `query:"identifier" required:"true"`
	Scope      string `query:"scope" required:"true"`
}

func (r *Role) APILeasesWOL() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APILeasesWOLInput, output *struct{}) error {
		l, ok := r.leases.GetPrefix(input.Identifier)
		if !ok {
			return status.InvalidArgument
		}
		err := l.sendWOL()
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dhcp.wol_leases")
	u.SetTitle("DHCP Leases")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}

type APILeasesDeleteInput struct {
	Identifier string `query:"identifier"`
	Scope      string `query:"scope"`
}

func (r *Role) APILeasesDelete() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APILeasesDeleteInput, output *struct{}) error {
		key := r.i.KV().Key(
			types.KeyRole,
			types.KeyLeases,
			input.Identifier,
		)
		_, err := r.i.KV().Delete(
			ctx,
			key.String(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dhcp.delete_leases")
	u.SetTitle("DHCP Leases")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}
