package tftp

import (
	"bytes"
	"context"
	"errors"
	"io"

	"beryju.io/gravity/pkg/o11y"
	"github.com/getsentry/sentry-go"
	"github.com/pin/tftp/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (r *Role) Writer(filename string, wt io.WriterTo) error {
	it := wt.(tftp.IncomingTransfer)
	r.log.Info("TFTP Write request", zap.String("filename", filename), zap.String("client", it.RemoteAddr().IP.String()))
	return r.writeHandler(filename, wt)
}

func (r *Role) writeHandler(filename string, wt io.WriterTo) error {
	it := wt.(tftp.IncomingTransfer)
	ctx, canc := context.WithCancel(context.Background())
	defer canc()
	sctx, span := o11y.Tracer.Start(ctx, filename)
	span.SetAttributes(attribute.String("http.request.method", "PUT"))
	defer span.End()
	hub := sentry.GetHubFromContext(sctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.Scope().SetUser(sentry.User{
		IPAddress: it.RemoteAddr().IP.String(),
	})

	s, ok := it.Size()
	if ok && s >= etcdMaxSize {
		return errors.New("file too big")
	}
	buf := bytes.NewBuffer([]byte{})
	s, err := wt.WriteTo(buf)
	if s >= etcdMaxSize {
		return errors.New("file too big")
	}
	if err != nil {
		return err
	}
	_, err = r.i.KV().Put(
		sctx,
		r.getPath(filename, it.RemoteAddr()).String(),
		buf.String(),
	)
	return err
}
