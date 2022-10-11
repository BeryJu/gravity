package backup

import (
	"encoding/json"
	"fmt"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles/backup/types"
	"github.com/minio/minio-go/v7"
)

type BackupStatus struct {
	Status   string    `json:"status,omitempty"`
	Error    string    `json:"error,omitempty"`
	Filename string    `json:"filename,omitempty"`
	Size     int64     `json:"size,omitempty"`
	Duration int64     `json:"duration,omitempty"`
	Time     time.Time `json:"time"`
}

const BackupStatusSuccess = "success"
const BackupStatusStarted = "started"
const BackupStatusFailed = "failed"

func (r *Role) setStatus(status *BackupStatus) {
	backupStatus.WithLabelValues(status.Status).SetToCurrentTime()
	if status.Status == BackupStatusSuccess {
		backupSize.Set(float64(status.Size))
		backupDuration.Set(float64(status.Duration))
	}
	jstatus, err := json.Marshal(status)
	if err != nil {
		r.log.WithError(err).Warning("failed to marshal status")
		return
	}
	_, err = r.i.KV().Put(
		r.ctx,
		r.i.KV().Key(
			types.KeyRole,
			types.KeyStatus,
			extconfig.Get().Instance.Identifier,
		).String(),
		string(jstatus),
	)
	if err != nil {
		r.log.WithError(err).Warning("failed to save status")
		return
	}
}

func (r *Role) GetBackupName() string {
	now := time.Now()
	fileName := fmt.Sprintf(
		"gravity-snapshot-%s-%d_%d_%d",
		extconfig.FullVersion(),
		now.Year(),
		now.Month(),
		now.Day(),
	)
	if r.cfg.Path != "" {
		fileName = fmt.Sprintf("%s/%s", r.cfg.Path, fileName)
	}
	return fileName
}

func (r *Role) SaveSnapshot() *BackupStatus {
	start := time.Now()
	status := &BackupStatus{
		Status: BackupStatusFailed,
		Time:   time.Now(),
	}
	defer r.setStatus(status)
	if r.mc == nil {
		status.Error = "backup not configured"
		return status
	}
	// TODO: Only let the master do backups to prevent duplicates
	read, err := r.i.KV().Snapshot(r.ctx)
	if err != nil {
		r.log.WithError(err).Warning("failed to snapshot")
		status.Error = err.Error()
		return status
	}
	fileName := r.GetBackupName()
	i, err := r.mc.PutObject(r.ctx, r.cfg.Bucket, fileName, read, -1, minio.PutObjectOptions{})
	if err != nil {
		r.log.WithError(err).Warning("failed to upload snapshot")
		status.Error = err.Error()
		return status
	}
	r.log.WithField("size", i.Size).Info("Uploaded snapshot")
	finish := time.Since(start)
	status.Status = BackupStatusSuccess
	status.Size = i.Size
	status.Filename = fileName
	status.Duration = int64(finish.Seconds())
	return status
}

func (r *Role) ensureBucket() {
	exists, err := r.mc.BucketExists(r.ctx, r.cfg.Bucket)
	if err == nil && exists {
		return
	}
	err = r.mc.MakeBucket(r.ctx, r.cfg.Bucket, minio.MakeBucketOptions{})
	if err != nil {
		r.log.WithError(err).Warning("failed to create bucket")
	}
}
