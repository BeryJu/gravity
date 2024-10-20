package tftp

import (
	"bytes"
	"context"
	"errors"
	"io"

	"beryju.io/gravity/pkg/roles/tftp/types"
	"github.com/getsentry/sentry-go"
	"github.com/pin/tftp/v3"
	"go.uber.org/zap"
)

func (r *Role) writeLogger(filename string, wt io.WriterTo) error {
	it := wt.(tftp.IncomingTransfer)
	r.log.Info("TFTP Write request", zap.String("filename", filename), zap.String("client", it.RemoteAddr().IP.String()))
	return r.writeHandler(filename, wt)
}

const maxSize = 1.5 * 1024

func (r *Role) writeHandler(filename string, wt io.WriterTo) error {
	it := wt.(tftp.IncomingTransfer)
	ctx, canc := context.WithCancel(context.Background())
	defer canc()
	span := sentry.StartTransaction(ctx, "gravity.tftp.request")
	defer span.Finish()
	hub := sentry.GetHubFromContext(span.Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.Scope().SetUser(sentry.User{
		IPAddress: it.RemoteAddr().IP.String(),
	})

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
		span.Context(),
		r.i.KV().Key(
			types.KeyRole,
			types.KeyFiles,
			it.RemoteAddr().IP.String(),
			filename,
		).String(),
		buf.String(),
	)
	return nil
}