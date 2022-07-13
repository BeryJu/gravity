package discovery

import (
	"context"
	"errors"
	"strings"

	"beryju.io/gravity/pkg/roles/discovery/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *DiscoveryRole) apiHandlerDevices() usecase.Interactor {
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
	u.SetTitle("Discovery devices")
	u.SetTags("discovery")
	u.SetDescription("List all discovered devices.")
	u.SetExpectedErrors(status.Internal)
	return u
}

func (r *DiscoveryRole) apiHandlerDeviceApply() usecase.Interactor {
	type deviceApplyInput struct {
		RelKey    string `query:"relKey"`
		DHCPScope string `query:"dhcpScope"`
		DNSZone   string `query:"dnsZone"`
	}
	u := usecase.NewIOI(new(deviceApplyInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*deviceApplyInput)
		)
		by := strings.SplitN(in.RelKey, "/", 1)[0]
		if by != types.KeyDevicesByMAC && by != types.KeyDevicesByIP {
			return status.Wrap(errors.New("invalid key"), status.InvalidArgument)
		}
		rawDevice, err := r.i.KV().Get(ctx, r.i.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			in.RelKey,
		).String())
		if err != nil {
			return status.Wrap(errors.New("invalid key"), status.InvalidArgument)
		}
		if len(rawDevice.Kvs) < 1 {
			return status.Wrap(errors.New("not found"), status.NotFound)
		}

		device := r.deviceFromKV(rawDevice.Kvs[0])
		if by == types.KeyDevicesByIP {
			err = device.toDHCP(in.DHCPScope)
		} else {
			err = device.toDNS(in.DNSZone)
		}
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetTitle("Apply Discovered devices")
	u.SetTags("discovery")
	u.SetDescription("Convert discovered device into DHCP lease/DNS record.")
	u.SetExpectedErrors(status.InvalidArgument, status.NotFound, status.Internal)
	return u
}
