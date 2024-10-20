package tftp

import (
	"context"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	apiTypes "beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/roles/tftp/types"
	"beryju.io/gravity/pkg/storage"
	"github.com/pin/tftp/v3"
	"github.com/swaggest/rest/web"
	"go.uber.org/zap"
)

const etcdMaxSize = 10 * 1024 * 1024
const sharedNamespace = "shared"

type Role struct {
	localfs fs.FS
	s       *tftp.Server

	log *zap.Logger
	i   roles.Instance
	ctx context.Context
	cfg *RoleConfig
}

func New(instance roles.Instance) *Role {
	r := &Role{
		log:     instance.Log(),
		i:       instance,
		ctx:     instance.Context(),
		localfs: os.DirFS(extconfig.Get().Dirs().TFTPLocalDir),
	}
	s := tftp.NewServer(r.Reader, r.Writer)
	r.s = s
	s.SetTimeout(5 * time.Second)
	r.i.AddEventListener(apiTypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/tftp/files", r.APIFilesGet())
		svc.Get("/api/v1/tftp/files/download", r.APIFilesDownload())
		svc.Post("/api/v1/tftp/files", r.APIFilesPut())
		svc.Delete("/api/v1/tftp/files", r.APIFilesDelete())
		svc.Get("/api/v1/roles/tftp", r.APIRoleConfigGet())
		svc.Post("/api/v1/roles/tftp", r.APIRoleConfigPut())
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.cfg = r.decodeRoleConfig(config)

	listen := extconfig.Get().Listen(r.cfg.Port)

	r.log.Info("starting tftp server", zap.String("listen", listen))
	go func() {
		err := r.s.ListenAndServe(listen)
		if err != nil && err != http.ErrServerClosed {
			r.log.Warn("failed to listen", zap.Error(err))
		}
	}()
	return nil
}

func (r *Role) Stop() {
	if r.s != nil {
		r.s.Shutdown()
	}
}

func (r *Role) getPath(filename string, addr net.UDPAddr) *storage.Key {
	bk := r.i.KV().Key(
		types.KeyRole,
		types.KeyFiles,
	)
	if strings.HasPrefix(filename, sharedNamespace) {
		bk.Add(filename)
	} else {
		bk.Add(addr.IP.String(), filename)
	}
	return bk
}
