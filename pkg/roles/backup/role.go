package backup

import (
	"context"
	"net/url"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
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
		Secure:    strings.EqualFold(endpoint.Scheme, "https"),
		Transport: extconfig.Transport(),
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

func (r *BackupRole) Stop() {
	if r.c != nil {
		<-r.c.Stop().Done()
	}
}
