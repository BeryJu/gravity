package backup

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/swaggest/rest/web"
)

const (
	KeyRole = "backup"
)

type BackupRole struct {
	mc  *minio.Client
	cfg *BackupRoleConfig
	c   *cron.Cron

	log *log.Entry
	i   roles.Instance
	ctx context.Context
}

func New(instance roles.Instance) *BackupRole {
	r := &BackupRole{
		log: instance.GetLogger().WithField("role", KeyRole),
		i:   instance,
	}
	r.i.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Post("/api/v1/backup/start", r.apiHandlerBackupStart())
	})
	return r
}

func (r *BackupRole) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeBackupRoleConfig(config)
	if !r.cfg.Enabled {
		return nil
	}
	endpoint, err := url.Parse(r.cfg.Endpoint)
	if err != nil {
		return err
	}
	opts := &minio.Options{
		Secure: strings.EqualFold(endpoint.Scheme, "https"),
	}
	if r.cfg.AccessKey != "" {
		opts.Creds = credentials.NewStaticV4(r.cfg.AccessKey, r.cfg.SecretKey, "")
	} else {
		opts.Creds = credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.EnvAWS{},
				&credentials.EnvMinio{},
				&credentials.FileAWSCredentials{},
				&credentials.FileMinioClient{},
				&credentials.IAM{},
			},
		)
	}
	minioClient, err := minio.New(endpoint.Host, opts)
	if err != nil {
		return err
	}
	r.mc = minioClient
	r.c = cron.New()
	ei, err := r.c.AddFunc(r.cfg.CronExpr, func() {
		r.saveSnapshot()
	})
	if err != nil {
		return err
	}
	r.log.WithField("next", r.c.Entry(ei).Next).Debug("next backup run")
	r.c.Start()
	go r.saveSnapshot()
	return nil
}

type BackupStatus struct {
	Status   string `json:"status,omitempty"`
	Error    error  `json:"error,omitempty"`
	Filename string `json:"filename,omitempty"`
	Size     int64  `json:"size,omitempty"`
}

const BackupStatusSuccess = "success"
const BackupStatusStarted = "started"
const BackupStatusFailed = "failed"

func (r *BackupRole) saveSnapshot() *BackupStatus {
	if r.mc == nil {
		return &BackupStatus{
			Status: BackupStatusFailed,
			Error:  errors.New("Backup not configured"),
		}
	}
	// TODO: Only let the master do backups to prevent duplicates
	read, err := r.i.KV().Snapshot(r.ctx)
	if err != nil {
		r.log.WithError(err).Warning("failed to snapshot")
		return &BackupStatus{
			Status: BackupStatusFailed,
			Error:  err,
		}
	}
	now := time.Now()
	fileName := fmt.Sprintf("gravity-snapshot-%d-%d-%d", now.Year(), now.Month(), now.Day())
	i, err := r.mc.PutObject(r.ctx, r.cfg.Bucket, fileName, read, -1, minio.PutObjectOptions{})
	if err != nil {
		r.log.WithError(err).Warning("failed to upload snapshot")
		return &BackupStatus{
			Status: BackupStatusFailed,
			Error:  err,
		}
	}
	r.log.WithField("size", i.Size).Info("Uploaded snapshot")
	return &BackupStatus{
		Status:   BackupStatusSuccess,
		Filename: fileName,
		Size:     i.Size,
	}
}

func (r *BackupRole) Stop() {
	if r.c != nil {
		<-r.c.Stop().Done()
	}
}
