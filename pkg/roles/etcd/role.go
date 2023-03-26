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
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.etcd.io/etcd/server/v3/embed"
)

const (
	relInstCertPath = "/instance.pem"
	relInstKeyPath  = "/instance_key.pem"
)

type Role struct {
	i roles.Instance

	e       *embed.Etcd
	cfg     *embed.Config
	log     *zap.Logger
	etcdDir string
	certDir string
}

func urlMustParse(raw string) url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return *u
}

func New(instance roles.Instance) *Role {
	dirs := extconfig.Get().Dirs()
	cfg := embed.NewConfig()
	ee := &Role{
		cfg:     cfg,
		log:     instance.Log(),
		i:       instance,
		etcdDir: dirs.EtcdDir,
		certDir: dirs.CertDir,
	}
	cfg.Dir = dirs.EtcdDir
	cfg.ZapLoggerBuilder = embed.NewZapCoreLoggerBuilder(
		extconfig.Get().BuildLoggerWithLevel(zapcore.WarnLevel).Named("role.etcd"),
		nil,
		nil,
	)
	cfg.AutoCompactionMode = "periodic"
	cfg.AutoCompactionRetention = "60m"
	// --listen-client-urls
	cfg.LCUrls = []url.URL{
		urlMustParse("http://localhost:2379"),
	}
	// --advertise-client-urls
	cfg.ACUrls = []url.URL{
		urlMustParse("http://localhost:2379"),
	}
	// --listen-peer-urls
	cfg.LPUrls = []url.URL{
		urlMustParse(fmt.Sprintf("https://%s", extconfig.Get().Listen(2380))),
	}
	// --initial-advertise-peer-urls
	cfg.APUrls = []url.URL{
		urlMustParse(fmt.Sprintf("https://%s:2380", extconfig.Get().Instance.IP)),
	}
	cfg.Name = extconfig.Get().Instance.Identifier
	cfg.InitialCluster = fmt.Sprintf("%s=https://%s:2380", cfg.Name, extconfig.Get().Instance.IP)
	cfg.PeerAutoTLS = true
	cfg.PeerTLSInfo.ClientCertFile = path.Join(ee.certDir, "peer", relInstCertPath)
	cfg.PeerTLSInfo.ClientKeyFile = path.Join(ee.certDir, "peer", relInstKeyPath)
	cfg.PeerTLSInfo.ClientCertAuth = true
	cfg.SelfSignedCertValidity = 1
	err := ee.prepareJoin(cfg)
	if err != nil {
		instance.Log().Warn("failed to join cluster", zap.Error(err))
		return nil
	}
	return ee
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

	res, _, err := apiClient.RolesEtcdApi.EtcdJoinMember(context.Background()).ApiAPIMemberJoinInput(
		api.ApiAPIMemberJoinInput{
			Peer: api.PtrString(
				fmt.Sprintf(
					"https://%s:2380",
					extconfig.Get().Instance.IP,
				),
			),
			Identifier: &extconfig.Get().Instance.Identifier,
			Roles:      &extconfig.Get().BootstrapRoles,
		},
	).Execute()
	if err != nil || res.EtcdInitialCluster == nil {
		return errors.Wrap(err, "failed to send api request to join")
	}
	cfg.ClusterState = embed.ClusterStateFlagExisting
	cfg.InitialCluster = res.GetEtcdInitialCluster()
	ee.log.Info("joining etcd cluster", zap.String("initialCluster", cfg.InitialCluster))
	return nil
}

func (ee *Role) Start(ctx context.Context) error {
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
	<-e.Server.ReadyNotify()
	ee.log.Info("embedded etcd Ready!", zap.Duration("runtime", time.Since(start)))
	return nil
}

func (ee *Role) Stop() {
	if ee.e == nil {
		return
	}
	ee.log.Info("stopping etcd")
	ee.e.Close()
}
