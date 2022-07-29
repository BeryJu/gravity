package backup

import (
	"errors"
	"fmt"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/minio/minio-go/v7"
)

type BackupStatus struct {
	Status   string `json:"status,omitempty"`
	Error    error  `json:"error,omitempty"`
	Filename string `json:"filename,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Duration int64  `json:"duration,omitempty"`
}

const BackupStatusSuccess = "success"
const BackupStatusStarted = "started"
const BackupStatusFailed = "failed"

func (r *Role) setStatus(status *BackupStatus) *BackupStatus {
	backupStatus.WithLabelValues(status.Status).SetToCurrentTime()
	if status.Status == BackupStatusSuccess {
		backupSize.Set(float64(status.Size))
		backupDuration.Set(float64(status.Duration))
	}
	return status
}

func (r *Role) saveSnapshot() *BackupStatus {
	start := time.Now()
	if r.mc == nil {
		return r.setStatus(&BackupStatus{
			Status: BackupStatusFailed,
			Error:  errors.New("Backup not configured"),
		})
	}
	// TODO: Only let the master do backups to prevent duplicates
	read, err := r.i.KV().Snapshot(r.ctx)
	if err != nil {
		r.log.WithError(err).Warning("failed to snapshot")
		return r.setStatus(&BackupStatus{
			Status: BackupStatusFailed,
			Error:  err,
		})
	}
	now := time.Now()
	fileName := fmt.Sprintf("gravity-snapshot-%s-%d_%d_%d", extconfig.FullVersion(), now.Year(), now.Month(), now.Day())
	if r.cfg.Path != "" {
		fileName = fmt.Sprintf("%s/%s", r.cfg.Path, fileName)
	}
	i, err := r.mc.PutObject(r.ctx, r.cfg.Bucket, fileName, read, -1, minio.PutObjectOptions{})
	if err != nil {
		r.log.WithError(err).Warning("failed to upload snapshot")
		return r.setStatus(&BackupStatus{
			Status: BackupStatusFailed,
			Error:  err,
		})
	}
	r.log.WithField("size", i.Size).Info("Uploaded snapshot")
	finish := time.Since(start)
	return r.setStatus(&BackupStatus{
		Status:   BackupStatusSuccess,
		Filename: fileName,
		Size:     i.Size,
		Duration: int64(finish.Seconds()),
	})
}
