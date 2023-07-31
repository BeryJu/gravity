package tftp

import (
	"bytes"
	"context"
	"io"
	"strings"

	"beryju.io/gravity/internal/resources"
	"beryju.io/gravity/pkg/roles/tftp/types"
	"github.com/pin/tftp/v3"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func (r *Role) readLogger(filename string, rf io.ReaderFrom) error {
	ot := rf.(tftp.OutgoingTransfer)
	r.log.Info("TFTP Read request", zap.String("filename", filename), zap.String("client", ot.RemoteAddr().IP.String()))
	return r.readHandler(filename, rf)
}

func (r *Role) readHandler(filename string, rf io.ReaderFrom) error {
	ot := rf.(tftp.OutgoingTransfer)
	var f io.Reader
	var err error
	if strings.HasPrefix(filename, "bundled/") {
		f, err = resources.TFTPRoot.Open(strings.Replace(filename, "bundled/", "tftp/", 1))
	} else if strings.HasPrefix(filename, "local") {
		f, err = r.localfs.Open(strings.Replace(filename, "local/", "", 1))
	} else {
		var re *clientv3.GetResponse
		re, err = r.i.KV().Get(context.Background(),
			r.i.KV().Key(
				types.KeyRole,
				types.KeyFiles,
				ot.RemoteAddr().IP.String(),
				filename,
			).String())
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
