package tftp

import (
	"github.com/pin/tftp/v3"
	"go.uber.org/zap"
)

func (r *Role) OnSuccess(stats tftp.TransferStats) {}

func (r *Role) OnFailure(stats tftp.TransferStats, err error) {
	r.log.Info("TFTP Error",
		zap.String("filename", stats.Filename),
		zap.String("client", stats.RemoteAddr.String()),
		zap.Error(err))
}
