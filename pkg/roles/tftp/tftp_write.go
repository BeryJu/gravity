package tftp

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/getsentry/sentry-go"
	"github.com/pin/tftp/v3"
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
	span := sentry.StartTransaction(ctx, filename)
	span.Op = "gravity.tftp.request"
	span.SetTag("http.request.method", "PUT")
	defer span.Finish()
	hub := sentry.GetHubFromContext(span.Context())
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
		span.Context(),
		r.getPath(filename, it.RemoteAddr()).String(),
		buf.String(),
	)
	return err
}
