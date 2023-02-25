package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles/backup/types"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type BackupStatus struct {
	Status   string    `json:"status,omitempty" required:"true"`
	Error    string    `json:"error,omitempty" required:"true"`
	Filename string    `json:"filename,omitempty" required:"true"`
	Size     int64     `json:"size,omitempty" required:"true"`
	Duration int64     `json:"duration,omitempty" required:"true"`
	Time     time.Time `json:"time" required:"true"`
}

const (
	BackupStatusSuccess = "success"
	BackupStatusStarted = "started"
	BackupStatusFailed  = "failed"
)

func (r *Role) setStatus(ctx context.Context, status *BackupStatus) {
	// Reset all backup statuses to 0
	for _, st := range []string{BackupStatusSuccess, BackupStatusStarted, BackupStatusFailed} {
		backupStatus.WithLabelValues(st).Set(0)
	}
	backupStatus.WithLabelValues(status.Status).Set(1)

	if status.Status == BackupStatusSuccess {
		backupSize.Set(float64(status.Size))
		backupDuration.Set(float64(status.Duration))
	}
	jstatus, err := json.Marshal(status)
	if err != nil {
		r.log.Warn("failed to marshal status", zap.Error(err))
		return
	}
	_, err = r.i.KV().Put(
		ctx,
		r.i.KV().Key(
			types.KeyRole,
			types.KeyStatus,
			extconfig.Get().Instance.Identifier,
		).String(),
		string(jstatus),
	)
	if err != nil {
		r.log.Warn("failed to save status", zap.Error(err))
		return
	}
}

func (r *Role) GetBackupName() string {
	now := time.Now()
	path := fmt.Sprintf("%s/gravity/%s/%s",
		r.cfg.Path,
		extconfig.FullVersion(),
		extconfig.Get().Instance.Identifier)
	fileName := fmt.Sprintf(
		"%s/gravity-snapshot-%d_%d_%d",
		path,
		now.Year(),
		now.Month(),
		now.Day(),
	)
	return fileName
}

func (r *Role) snapshotToFile(ctx context.Context) (*os.File, error) {
	reader, err := r.i.KV().Snapshot(ctx)
	if err != nil {
		r.log.Warn("failed to snapshot", zap.Error(err))
		return nil, err
	}
	file, err := os.CreateTemp(os.TempDir(), "gravity-snapshot.*.etcd")
	if err != nil {
		return nil, err
	}
	r.log.Info("creating local snapshot", zap.String("path", file.Name()))

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := file.Write(buf[:n]); err != nil {
			return nil, err
		}
	}
	return file, nil
}

func (r *Role) SaveSnapshot(ctx context.Context) *BackupStatus {
	start := time.Now()
	status := &BackupStatus{
		Status: BackupStatusFailed,
		Time:   time.Now(),
	}
	defer r.setStatus(ctx, status)
	file, err := r.snapshotToFile(ctx)
	if err != nil {
		status.Error = err.Error()
		return status
	}
	newPath := path.Join(extconfig.Get().Dirs().BackupDir, "gravity-local-snapshot")
	os.Rename(file.Name(), newPath)
	if r.mc == nil {
		status.Status = BackupStatusSuccess
		return status
	}
	file.Seek(0, io.SeekStart)
	stat, err := file.Stat()
	if err != nil {
		status.Error = err.Error()
		return status
	}
	fileName := r.GetBackupName()
	i, err := r.mc.PutObject(r.ctx, r.cfg.Bucket, fileName, file, stat.Size(), minio.PutObjectOptions{})
	if err != nil {
		r.log.Warn("failed to upload snapshot", zap.Error(err))
		status.Error = err.Error()
		return status
	}
	r.log.Info("Uploaded snapshot", zap.Int64("size", i.Size))
	finish := time.Since(start)
	status.Status = BackupStatusSuccess
	status.Size = i.Size
	status.Filename = fileName
	status.Duration = int64(finish.Seconds())
	status.Time = time.Now()
	return status
}

func (r *Role) ensureBucket() {
	exists, err := r.mc.BucketExists(r.ctx, r.cfg.Bucket)
	if err == nil && exists {
		return
	}
	err = r.mc.MakeBucket(r.ctx, r.cfg.Bucket, minio.MakeBucketOptions{})
	if err != nil {
		r.log.Warn("failed to create bucket", zap.Error(err))
	}
}
