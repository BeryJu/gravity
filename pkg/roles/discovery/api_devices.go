package discovery

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/roles/discovery/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type APIDevicesGetInput struct {
	Identifier string `query:"identifier"  description:"Optionally get device by identifier"`
}
type APIDevice struct {
	Identifier string `json:"identifier" required:"true"`
	Hostname   string `json:"hostname" required:"true"`
	IP         string `json:"ip" required:"true"`
	MAC        string `json:"mac" required:"true"`
}
type APIDevicesGetOutput struct {
	Devices []APIDevice `json:"devices"`
}

func (r *Role) APIDevicesGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIDevicesGetInput, output *APIDevicesGetOutput) error {
		key := r.i.KV().Key(types.KeyRole, types.KeyDevices)
		if input.Identifier == "" {
			key = key.Prefix(true)
		} else {
			key = key.Add(input.Identifier)
		}
		rawDevices, err := r.i.KV().Get(ctx, key.String(), clientv3.WithPrefix())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, rawDev := range rawDevices.Kvs {
			dev := r.deviceFromKV(rawDev)
			output.Devices = append(output.Devices, APIDevice{
				Identifier: dev.Identifier,
				Hostname:   dev.Hostname,
				IP:         dev.IP,
				MAC:        dev.MAC,
			})
		}
		return nil
	})
	u.SetName("discovery.get_devices")
	u.SetTitle("Discovery devices")
	u.SetTags("roles/discovery")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APIDevicesApplyInput struct {
	Identifier string `query:"identifier" required:"true"`
	To         string `json:"to" enum:"dhcp,dns" required:"true"`
	DHCPScope  string `json:"dhcpScope" required:"true"`
	DNSZone    string `json:"dnsZone" required:"true"`
}

func (r *Role) APIDevicesApply() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIDevicesApplyInput, output *struct{}) error {
		rawDevice, err := r.i.KV().Get(ctx, r.i.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			input.Identifier,
		).String())
		if err != nil {
			return status.Wrap(errors.New("invalid key"), status.InvalidArgument)
		}
		if len(rawDevice.Kvs) < 1 {
			return status.Wrap(errors.New("not found"), status.NotFound)
		}

		device := r.deviceFromKV(rawDevice.Kvs[0])
		if input.To == "dhcp" {
			err = device.toDHCP(ctx, input.DHCPScope)
		} else {
			err = device.toDNS(ctx, input.DNSZone)
		}
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		_, err = r.i.KV().Delete(ctx, r.i.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			device.Identifier,
		).String())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("discovery.apply_device")
	u.SetTitle("Apply Discovered devices")
	u.SetTags("roles/discovery")
	u.SetExpectedErrors(status.InvalidArgument, status.NotFound, status.Internal)
	return u
}

type APIDevicesDeleteInput struct {
	Name string `query:"identifier"`
}

func (r *Role) APIDevicesDelete() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIDevicesDeleteInput, output *struct{}) error {
		_, err := r.i.KV().Delete(ctx, r.i.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			input.Name,
		).String())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("discovery.delete_devices")
	u.SetTitle("Discovery devices")
	u.SetTags("roles/discovery")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}
