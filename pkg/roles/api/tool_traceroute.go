package api

import (
	"context"
	"time"

	"github.com/aeden/traceroute"
	"github.com/pkg/errors"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *Role) Traceroute(addr string) (traceroute.TracerouteResult, error) {
	options := traceroute.TracerouteOptions{}
	options.SetRetries(5)
	options.SetMaxHops(64)
	options.SetFirstHop(1)
	res, err := traceroute.Traceroute(addr, &options)
	return res, err
}

type APIToolTracerouteInput struct {
	Host string `json:"host"`
}

type APIToolTracerouteOutputHop struct {
	Address     string        `json:"address"`
	Success     bool          `json:"success"`
	ElapsedTime time.Duration `json:"elapsedTime"`
}
type APIToolTracerouteOutput struct {
	Hops []APIToolTracerouteOutputHop `json:"hops"`
}

func (r *Role) APIToolTraceroute() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIToolTracerouteInput, output *APIToolTracerouteOutput) error {
		hops, err := r.Traceroute(input.Host)
		if err != nil {
			return errors.Wrap(err, "failed to traceroute")
		}
		for _, h := range hops.Hops {
			ah := APIToolTracerouteOutputHop{
				Success:     h.Success,
				ElapsedTime: h.ElapsedTime,
			}
			if h.Host != "" {
				ah.Address = h.Host
			} else {
				ah.Address = h.AddressString()
			}
			output.Hops = append(output.Hops, ah)
		}
		return nil
	})
	u.SetName("tools.traceroute")
	u.SetTitle("Traceroute tool")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}
