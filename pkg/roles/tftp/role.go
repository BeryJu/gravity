package tftp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/tftp/types"
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
	s := tftp.NewServer(nil, r.writeLogger)
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

func (r *Role) writeLogger(filename string, wt io.WriterTo) error {
	it := wt.(tftp.IncomingTransfer)
	r.log.Info("TFTP Write request", zap.String("filename", filename), zap.String("client", it.RemoteAddr().IP.String()))
	return r.writeHandler(filename, wt)
}

const maxSize = 1.5 * 1024

func (r *Role) writeHandler(filename string, wt io.WriterTo) error {
	it := wt.(tftp.IncomingTransfer)
	s, ok := it.Size()
	if ok && s >= maxSize {
		return errors.New("file too big")
	}
	buf := bytes.NewBuffer([]byte{})
	s, err := wt.WriteTo(buf)
	if s >= maxSize {
		return errors.New("file too big")
	}
	if err != nil {
		return err
	}
	r.i.KV().Put(
		context.Background(),
		r.i.KV().Key(
			types.KeyRole,
			types.KeyFiles,
			filename,
		).String(),
		buf.String(),
	)
	return nil
}
