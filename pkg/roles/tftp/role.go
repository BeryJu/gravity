package tftp

import (
	"context"
	"net/http"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"github.com/pin/tftp/v3"
	"go.uber.org/zap"
)

type Role struct {
	s      *tftp.Server
	log    *zap.Logger
	i      roles.Instance
	ctx    context.Context
	server *http.Server
}

func New(instance roles.Instance) *Role {
	r := &Role{
		log: instance.Log(),
		i:   instance,
		ctx: instance.Context(),
	}
	s := tftp.NewServer(r.readHandler, r.writeLogger)
	r.s = s
	s.SetTimeout(5 * time.Second) // optional
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	listen := extconfig.Get().Listen(69)

	r.log.Info("starting debug server", zap.String("listen", listen))
	go func() {
		err := r.s.ListenAndServe(listen)
		if err != nil && err != http.ErrServerClosed {
			r.log.Warn("failed to listen", zap.Error(err))
		}
	}()
	return nil
}

func (r *Role) Stop() {
	if r.server != nil {
		r.server.Shutdown(r.ctx)
	}
}
