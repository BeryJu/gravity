package backup

import (
	"context"
	"net/url"
	"os"
	"strings"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/robfig/cron/v3"
	"github.com/swaggest/rest/web"
	"go.uber.org/zap"
)

const (
	EventTopicBackupRun = "roles.backup.run"
)

type Role struct {
	mc  *minio.Client
	cfg *RoleConfig
	c   *cron.Cron

	log *zap.Logger
	i   roles.Instance
	ctx context.Context
}

func init() {
	roles.Register("backup", func(i roles.Instance) roles.Role {
		return New(i)
	})
}

func New(instance roles.Instance) *Role {
	r := &Role{
		log: instance.Log(),
		i:   instance,
		ctx: instance.Context(),
	}
	r.i.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Post("/api/v1/backup/start", r.APIBackupStart())
		svc.Get("/api/v1/backup/status", r.APIBackupStatus())
		svc.Get("/api/v1/roles/backup", r.APIRoleConfigGet())
		svc.Post("/api/v1/roles/backup", r.APIRoleConfigPut())
	})
	r.i.AddEventListener(EventTopicBackupRun, func(ev *roles.Event) {
		r.SaveSnapshot(ev.Context)
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	err := os.MkdirAll(extconfig.Get().Dirs().BackupDir, os.ModeSticky|os.ModePerm)
	if err != nil {
		return err
	}
	r.cfg = r.decodeRoleConfig(config)
	if r.cfg.Endpoint == "" {
		return roles.ErrRoleNotConfigured
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
	r.c = cron.New()
	if r.cfg.CronExpr != "" {
		ei, err := r.c.AddFunc(r.cfg.CronExpr, func() {
			ctx, cancel := context.WithTimeout(r.ctx, 60*time.Minute)
			defer cancel()
			r.SaveSnapshot(ctx)
		})
		if err != nil {
			return err
		}
		r.log.Info("next backup run", zap.String("next", r.c.Entry(ei).Next.String()))
		r.c.Start()
	}
	minioClient, err := minio.New(endpoint.Host, opts)
	if err != nil {
		return err
	}
	r.mc = minioClient
	r.ensureBucket()
	return nil
}

func (r *Role) Stop() {
	if r.c != nil {
		<-r.c.Stop().Done()
	}
}
