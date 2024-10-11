package api

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *Role) Ping(addr string) (*probing.Statistics, error) {
	pinger, err := probing.NewPinger(addr)
	if err != nil {
		return nil, err
	}
	if os.Getuid() < 1 {
		pinger.SetPrivileged(true)
	}
	pinger.Count = 5
	err = pinger.Run()
	if err != nil {
		return nil, err
	}
	stats := pinger.Statistics()
	return stats, nil
}

type APIToolPingInput struct {
	Host string `json:"host"`
}

type APIToolPingOutput struct {
	PacketsRecv           int           `json:"packetsRecv"`
	PacketsSent           int           `json:"packetsSent"`
	PacketsRecvDuplicates int           `json:"packetsRecvDuplicates"`
	PacketLoss            float64       `json:"packetLoss"`
	MinRtt                time.Duration `json:"minRtt"`
	MaxRtt                time.Duration `json:"maxRtt"`
	AvgRtt                time.Duration `json:"avgRtt"`
	StdDevRtt             time.Duration `json:"stdDevRtt"`
}

func (r *Role) APIToolPing() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIToolPingInput, output *APIToolPingOutput) error {
		out, err := r.Ping(input.Host)
		if err != nil {
			return errors.Wrap(err, "failed to ping")
		}
		output.PacketsRecv = out.PacketsRecv
		output.PacketsSent = out.PacketsSent
		output.PacketsRecvDuplicates = out.PacketsRecvDuplicates
		output.PacketLoss = out.PacketLoss
		output.MinRtt = out.MinRtt
		output.MaxRtt = out.MaxRtt
		output.AvgRtt = out.AvgRtt
		output.StdDevRtt = out.StdDevRtt
		return nil
	})
	u.SetName("tools.ping")
	u.SetTitle("Ping tool")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}
