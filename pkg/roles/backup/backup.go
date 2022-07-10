package backup

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"beryju.io/gravity/pkg/roles"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
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
	return &BackupRole{
		log: instance.GetLogger().WithField("role", KeyRole),
		i:   instance,
	}
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

func (r *BackupRole) saveSnapshot() {
	read, err := r.i.KV().Snapshot(r.ctx)
	if err != nil {
		r.log.WithError(err).Warning("failed to snapshot")
		return
	}
	now := time.Now()
	fileName := fmt.Sprintf("gravity-snapshot-%d-%d-%d", now.Year(), now.Month(), now.Day())
	i, err := r.mc.PutObject(r.ctx, r.cfg.Bucket, fileName, read, -1, minio.PutObjectOptions{})
	if err != nil {
		r.log.WithError(err).Warning("failed to upload snapshot")
		return
	}
	r.log.WithField("size", i.Size).Info("Uploaded snapshot")
}

func (r *BackupRole) Stop() {
	if r.c != nil {
		<-r.c.Stop().Done()
	}
}
