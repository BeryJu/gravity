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

const (
	EventTopicBackupRun = "roles.backup.run"
)

type Role struct {
	mc  *minio.Client
	cfg *RoleConfig
	c   *cron.Cron

	log *log.Entry
	i   roles.Instance
	ctx context.Context
}

func New(instance roles.Instance) *Role {
	r := &Role{
		log: instance.Log(),
		i:   instance,
	}
	r.i.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Post("/api/v1/backup/start", r.apiHandlerBackupStart())
		svc.Get("/api/v1/roles/backup", r.apiHandlerRoleConfigGet())
		svc.Post("/api/v1/roles/backup", r.apiHandlerRoleConfigPut())
	})
	r.i.AddEventListener(EventTopicBackupRun, func(ev *roles.Event) {
		r.saveSnapshot()
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
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

func (r *Role) Stop() {
	if r.c != nil {
		<-r.c.Stop().Done()
	}
}
