package etcd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/avast/retry-go/v4"
	"github.com/pkg/errors"
	"github.com/swaggest/rest/web"
	"go.uber.org/zap"

	"go.etcd.io/etcd/server/v3/embed"
)

const (
	relInstCertPath = "/instance.pem"
	relInstKeyPath  = "/instance_key.pem"
)

func init() {
	roles.Register("etcd", func(i roles.Instance) roles.Role {
		return New(i)
	})
}

type Role struct {
	i roles.Instance

	e       *embed.Etcd
	cfg     *embed.Config
	log     *zap.Logger
	etcdDir string
	certDir string

	lcr     *LeaderClusterReconciler
	joining bool
}

func urlMustParse(raw string) url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return *u
}

var errRetry = errors.New("retry")

func New(instance roles.Instance) *Role {
	dirs := extconfig.Get().Dirs()
	cfg := embed.NewConfig()
	r := &Role{
		cfg:     cfg,
		log:     instance.Log(),
		i:       instance,
		etcdDir: dirs.EtcdDir,
		certDir: dirs.CertDir,
	}
	r.lcr = NewLeaderClusterConciler(instance)
	cfg.Dir = dirs.EtcdDir
	cfg.ZapLoggerBuilder = embed.NewZapLoggerBuilder(instance.Log())
	cfg.AutoCompactionMode = "periodic"
	cfg.AutoCompactionRetention = "60m"
	cfg.ListenClientUrls = []url.URL{
		urlMustParse("http://localhost:2379"),
		urlMustParse(fmt.Sprintf("http://%s", extconfig.Get().Listen(int32(extconfig.Get().Etcd.ClientPort)))),
	}
	cfg.AdvertiseClientUrls = []url.URL{
		urlMustParse("http://localhost:2379"),
		urlMustParse(fmt.Sprintf("http://%s", extconfig.Listen(extconfig.Get().Instance.IP, extconfig.Get().Etcd.ClientPort))),
	}
	cfg.ListenPeerUrls = []url.URL{
		urlMustParse(fmt.Sprintf("https://%s", extconfig.Get().Listen(int32(extconfig.Get().Etcd.PeerPort)))),
	}
	cfg.AdvertisePeerUrls = []url.URL{
		urlMustParse(fmt.Sprintf("https://%s", extconfig.Listen(extconfig.Get().Instance.IP, extconfig.Get().Etcd.PeerPort))),
	}
	cfg.Name = extconfig.Get().Instance.Identifier
	cfg.InitialCluster = fmt.Sprintf("%s=https://%s", cfg.Name, extconfig.Listen(extconfig.Get().Instance.IP, extconfig.Get().Etcd.PeerPort))
	cfg.PeerAutoTLS = true
	cfg.PeerTLSInfo.ClientCertFile = path.Join(r.certDir, "peer", relInstCertPath)
	cfg.PeerTLSInfo.ClientKeyFile = path.Join(r.certDir, "peer", relInstKeyPath)
	cfg.PeerTLSInfo.ClientCertAuth = true
	cfg.SelfSignedCertValidity = 1
	cfg.MaxRequestBytes = 10 * 1024 * 1024 // 10 MB
	err := r.prepareJoin(cfg)
	if err != nil {
		instance.Log().Warn("failed to join cluster", zap.Error(err))
		err = os.RemoveAll(path.Join(r.etcdDir, "member"))
		if err != nil {
			r.log.Warn("failed to remove etcd data", zap.Error(err))
		}
		return nil
	}
	r.i.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/etcd/members", r.APIClusterMembers())
		svc.Post("/api/v1/etcd/join", r.APIClusterJoin())
		svc.Post("/api/v1/etcd/move_leader", r.APIClusterMoveLeader())
		svc.Delete("/api/v1/etcd/remove", r.APIClusterRemove())
	})
	return r
}

func (ee *Role) prepareJoin(cfg *embed.Config) error {
	join := extconfig.Get().Etcd.JoinCluster
	if join == "" {
		return nil
	}

	// Don't attempt to join if we have an etcd directory already
	if _, err := os.Stat(path.Join(ee.etcdDir, "member")); err == nil {
		return nil
	}

	joinParts := strings.SplitN(join, ",", 2)
	if len(joinParts) < 2 {
		return fmt.Errorf("join string must consist of two parts: <token>,<api url>")
	}
	token := joinParts[0]
	apiUrl := joinParts[1]

	ee.log.Info("joining etcd cluster", zap.String("peer", apiUrl))

	u, err := url.Parse(apiUrl)
	if err != nil {
		return errors.Wrap(err, "failed to parse API url")
	}

	config := api.NewConfiguration()
	config.Host = u.Host
	config.Scheme = u.Scheme
	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	apiClient := api.NewAPIClient(config)

	res, err := retry.DoWithData(
		func() (*api.EtcdAPIMemberJoinOutput, error) {
			res, hr, err := apiClient.RolesEtcdAPI.EtcdJoinMember(context.Background()).EtcdAPIMemberJoinInput(
				api.EtcdAPIMemberJoinInput{
					Peer: api.PtrString(
						fmt.Sprintf(
							"https://%s:%d",
							extconfig.Get().Instance.IP,
							extconfig.Get().Etcd.PeerPort,
						),
					),
					Identifier: &extconfig.Get().Instance.Identifier,
					Roles:      &extconfig.Get().BootstrapRoles,
				},
			).Execute()
			if hr != nil && hr.StatusCode == 500 {
				return nil, errRetry
			}
			if err != nil || res.EtcdInitialCluster == nil {
				return nil, err
			}
			return res, nil
		},
		retry.DelayType(retry.BackOffDelay),
		retry.AttemptsForError(20, errRetry),
		retry.OnRetry(func(attempt uint, err error) {
			ee.log.Info("Join attempt", zap.Uint("attempt", attempt), zap.Error(err))
		}),
	)
	if err != nil {
		ee.log.Info("Failed to join", zap.Error(err))
		return err
	}
	cfg.ClusterState = embed.ClusterStateFlagExisting
	cfg.InitialCluster = res.GetEtcdInitialCluster()
	ee.log.Info("joining etcd cluster", zap.String("initialCluster", cfg.InitialCluster))
	ee.joining = true
	return nil
}

func (ee *Role) Config() *embed.Config {
	return ee.cfg
}

func (ee *Role) Start(ctx context.Context, cfg []byte) error {
	start := time.Now()
	ee.log.Info("starting embedded etcd")
	e, err := embed.StartEtcd(ee.cfg)
	if err != nil {
		return err
	}
	ee.e = e
	go func() {
		err := <-e.Err()
		if err != nil {
			ee.log.Warn("failed to start/stop etcd", zap.Error(err))
		}
	}()
	go ee.watchLeader()
	<-e.Server.ReadyNotify()
	ee.log.Info("embedded etcd Ready!", zap.Duration("runtime", time.Since(start)))
	if ee.joining {
		return ee.waitForPromotion()
	}
	return nil
}

func (ee *Role) watchLeader() {
	for {
		<-ee.e.Server.LeaderChangedNotify()
		ee.log.Info("etcd Leader changed", zap.Uint64("new-id", ee.e.Server.Lead()))
		if ee.e.Server.Leader() == ee.e.Server.MemberID() {
			ee.log.Info("We're the leader now.")
			ee.lcr.Start()
		} else {
			ee.log.Info("We're no longer the leader.")
			ee.lcr.Stop()
		}
	}
}

func (ee *Role) waitForPromotion() error {
	return retry.Do(
		func() error {
			l := ee.e.Server.IsLearner()
			if l {
				return errors.New("etcd is not promoted yet")
			}
			return nil
		},
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(50),
		retry.OnRetry(func(attempt uint, err error) {
			ee.log.Info("Checking if we're still learner", zap.Uint("attempt", attempt), zap.Error(err))
		}),
	)
}

func (ee *Role) Stop() {
	if ee.e == nil {
		return
	}
	ee.log.Info("stopping etcd")
	ee.e.Server.Stop()
	ee.e.Close()
}
