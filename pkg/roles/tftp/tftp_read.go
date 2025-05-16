package tftp

import (
	"bytes"
	"context"
	"io"
	"strings"

	"beryju.io/gravity/internal/resources"
	"github.com/getsentry/sentry-go"
	"github.com/pin/tftp/v3"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func (r *Role) Reader(filename string, rf io.ReaderFrom) error {
	ot := rf.(tftp.OutgoingTransfer)
	r.log.Info("TFTP Read request", zap.String("filename", filename), zap.String("client", ot.RemoteAddr().IP.String()))
	return r.readHandler(filename, rf)
}

func (r *Role) readHandler(filename string, rf io.ReaderFrom) error {
	ot := rf.(tftp.OutgoingTransfer)
	ctx, canc := context.WithCancel(context.Background())
	defer canc()
	span := sentry.StartTransaction(ctx, filename)
	span.Op = "gravity.tftp.request"
	span.SetTag("http.request.method", "GET")
	defer span.Finish()
	hub := sentry.GetHubFromContext(span.Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.Scope().SetUser(sentry.User{
		IPAddress: ot.RemoteAddr().IP.String(),
	})

	var f io.Reader
	var err error
	if strings.HasPrefix(filename, "bundled/") {
		f, err = resources.TFTPRoot.Open(strings.Replace(filename, "bundled/", "tftp/", 1))
	} else if strings.HasPrefix(filename, "local") && r.cfg.EnableLocal {
		f, err = r.localfs.Open(strings.Replace(filename, "local/", "", 1))
	} else {
		var re *clientv3.GetResponse
		re, err = r.i.KV().Get(
			span.Context(),
			r.getPath(filename, ot.RemoteAddr()).String(),
		)
		if err != nil || len(re.Kvs) < 1 {
			return err
		}
		f = bytes.NewBuffer(re.Kvs[0].Value)
	}
	if err != nil {
		return err
	}
	_, err = rf.ReadFrom(f)
	return err
}
