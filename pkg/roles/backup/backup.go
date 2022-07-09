package backup

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"beryju.io/ddet/pkg/roles"
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

	log *log.Entry
	i   roles.Instance
}

func New(instance roles.Instance) *BackupRole {
	return &BackupRole{
		log: instance.GetLogger().WithField("role", KeyRole),
		i:   instance,
	}
}

func (r *BackupRole) Start(config []byte) error {
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
	c := cron.New()
	ei, err := c.AddFunc(r.cfg.CronExpr, func() {
		r.saveSnapshot()
	})
	if err != nil {
		return err
	}
	r.log.WithField("next", c.Entry(ei).Next).Debug("next backup run")
	c.Start()
	go r.saveSnapshot()
	return nil
}

func (r *BackupRole) saveSnapshot() {
	read, err := r.i.KV().Snapshot(context.Background())
	if err != nil {
		r.log.WithError(err).Warning("failed to snapshot")
		return
	}
	now := time.Now()
	fileName := fmt.Sprintf("ddet-snapshot-%d-%d-%d", now.Year(), now.Month(), now.Day())
	i, err := r.mc.PutObject(context.Background(), r.cfg.Bucket, fileName, read, -1, minio.PutObjectOptions{})
	if err != nil {
		r.log.WithError(err).Warning("failed to upload snapshot")
		return
	}
	r.log.WithField("size", i.Size).Info("Uploaded snapshot")
}

func (r *BackupRole) Stop() {}
