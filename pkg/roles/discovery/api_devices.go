package discovery

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/roles/discovery/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *Role) apiHandlerDevices() usecase.Interactor {
	type device struct {
		Hostname string `json:"hostname"`
		IP       string `json:"ip"`
		MAC      string `json:"mac"`
	}
	type devicesOutput struct {
		Devices []device `json:"devices"`
	}
	u := usecase.NewIOI(new(struct{}), new(devicesOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*devicesOutput)
		)
		rawDevices, err := r.i.KV().Get(ctx, r.i.KV().Key(
			types.KeyRole,
			types.KeyDevices,
		).Prefix(true).String(), clientv3.WithPrefix())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		out.Devices = make([]device, 0)
		for _, rawDev := range rawDevices.Kvs {
			dev := r.deviceFromKV(rawDev)
			out.Devices = append(out.Devices, device{
				Hostname: dev.Hostname,
				IP:       dev.IP,
				MAC:      dev.MAC,
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

func (r *Role) apiHandlerDeviceApply() usecase.Interactor {
	type deviceApplyInput struct {
		Identifier string `path:"identifier"`
		To         string `json:"to" enum:"dhcp,dns"`
		DHCPScope  string `json:"dhcpScope"`
		DNSZone    string `json:"dnsZone"`
	}
	u := usecase.NewIOI(new(deviceApplyInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*deviceApplyInput)
		)
		rawDevice, err := r.i.KV().Get(ctx, r.i.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			in.Identifier,
		).String())
		if err != nil {
			return status.Wrap(errors.New("invalid key"), status.InvalidArgument)
		}
		if len(rawDevice.Kvs) < 1 {
			return status.Wrap(errors.New("not found"), status.NotFound)
		}

		device := r.deviceFromKV(rawDevice.Kvs[0])
		if in.To == "dhcp" {
			err = device.toDHCP(in.DHCPScope)
		} else {
			err = device.toDNS(in.DNSZone)
		}
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

func (r *Role) apiHandlerDevicesDelete() usecase.Interactor {
	type devicesInput struct {
		Name string `path:"identifier"`
	}
	u := usecase.NewIOI(new(devicesInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*devicesInput)
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
	u.SetName("discovery.delete_devices")
	u.SetTitle("Discovery devices")
	u.SetTags("roles/discovery")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}
