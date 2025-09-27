package instance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/storage"
)

var (
	ErrRoleRestarting   = errors.New("Role restarting")
	ErrRoleStopping     = errors.New("Role stopping")
	ErrInstanceStopping = errors.New("Instance stopping")
)

type ClusterInfo struct {
	Setup bool `json:"setup"`
}

type RoleContext struct {
	Role              roles.Role
	RoleInstance      *RoleInstance
	ContextCancelFunc context.CancelCauseFunc
}

type Instance struct {
	rootContext context.Context
	roles       map[string]RoleContext
	rolesM      sync.Mutex
	kv          *storage.Client
	log         *zap.Logger

	eventHandlers     map[string]map[string][]roles.EventHandler
	eventHandlersM    sync.RWMutex
	rootContextCancel context.CancelCauseFunc

	identifier string

	instanceSession *concurrency.Session
	etcd            roles.Role
}

func New() *Instance {
	extCfg := extconfig.Get()
	ctx, canc := context.WithCancelCause(context.Background())
	log := extCfg.Logger().With(zap.String("instance", extCfg.Instance.Identifier)).Named("instance")
	return &Instance{
		roles:             make(map[string]RoleContext),
		rolesM:            sync.Mutex{},
		log:               log,
		identifier:        extCfg.Instance.Identifier,
		eventHandlers:     make(map[string]map[string][]roles.EventHandler),
		eventHandlersM:    sync.RWMutex{},
		kv:                extCfg.EtcdClient(),
		rootContext:       ctx,
		rootContextCancel: canc,
	}
}

func (i *Instance) Role(id string) roles.Role {
	role, ok := i.roles[id]
	if !ok {
		return nil
	}
	return role.Role
}

func (i *Instance) Start() {
	i.log.Info("Gravity starting", zap.String("version", extconfig.FullVersion()))
	i.startSentry()
	bs := sentry.StartTransaction(i.rootContext, "gravity.instance.bootstrap")
	if strings.Contains(extconfig.Get().BootstrapRoles, "etcd") {
		if !i.startEtcd(bs.Context()) {
			return
		}
	}
	i.bootstrap(bs.Context())
	<-i.rootContext.Done()
}

func (i *Instance) startEtcd(ctx context.Context) bool {
	i.log.Info("'etcd' in bootstrap roles, starting embedded etcd")
	es := sentry.TransactionFromContext(ctx).StartChild("gravity.instance.bootstrap_etcd")
	defer es.Finish()
	i.etcd = roles.GetRole("etcd")(i.ForRole("etcd", es.Context()))
	if i.etcd == nil {
		i.Stop()
		return false
	}
	err := i.etcd.Start(es.Context(), []byte{})
	if err != nil {
		i.log.Warn("failed to start etcd", zap.Error(err))
		i.Stop()
		return false
	}
	return true
}

func (i *Instance) startSentry() {
	if !extconfig.Get().Sentry.Enabled || extconfig.Get().CI {
		return
	}
	release := fmt.Sprintf("gravity@%s", extconfig.FullVersion())
	rate := 0.5
	if extconfig.Get().Debug {
		rate = 1
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              extconfig.Get().Sentry.DSN,
		Release:          release,
		EnableTracing:    true,
		TracesSampleRate: rate,
		HTTPTransport:    extconfig.NewUserAgentTransport(release, extconfig.Transport()),
		Debug:            extconfig.Get().Debug,
		DebugWriter:      NewSentryWriter(i.log.Named("sentry")),
	})
	if err != nil {
		i.log.Warn("failed to init sentry", zap.Error(err))
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("gravity.instance", extconfig.Get().Instance.Identifier)
		scope.SetTag("gravity.version", extconfig.Version)
		scope.SetTag("gravity.hash", extconfig.BuildHash)
	})
}

func (i *Instance) Log() *zap.Logger {
	return i.log
}

func (i *Instance) getRoles(ctx context.Context) []string {
	rr, err := i.kv.Get(
		ctx,
		i.kv.Key(
			types.KeyInstance,
			i.identifier,
			types.KeyRoles,
		).String(),
	)
	roles := extconfig.Get().BootstrapRoles
	if err == nil && len(rr.Kvs) > 0 {
		roles = string(rr.Kvs[0].Value)
		i.log.Info("roles configured for instance", zap.Strings("roles", strings.Split(roles, types.RoleSeparator)))
	} else {
		i.log.Info("defaulting to bootstrap roles", zap.Strings("roles", strings.Split(roles, types.RoleSeparator)))
	}
	return strings.Split(roles, types.RoleSeparator)
}

func (i *Instance) eventRoleRestart(ev *roles.Event) {
	id := ev.Payload.Data["id"].(string)
	config := ev.Payload.Data["config"].([]byte)
	ctx := context.Background()
	tx := sentry.StartTransaction(ctx, "gravity.instance.role.restart")
	tx.Description = id
	tx.SetTag("gravity.role", id)
	defer tx.Finish()
	i.stopRole(tx.Context(), ErrRoleRestarting, id)
	i.startRole(tx.Context(), id, config)
}

func (i *Instance) checkFirstStart(ctx context.Context) {
	inst := i.ForRole("root", ctx)
	cluster, err := inst.KV().Get(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyCluster,
		).String(),
	)
	if err != nil {
		return
	}
	if len(cluster.Kvs) > 0 {
		return
	}
	i.log.Info("Initial startup")
	i.autoImportConfig()
	inst.DispatchEvent(
		types.EventTopicInstanceFirstStart,
		roles.NewEvent(ctx, map[string]interface{}{}),
	)

	clusterJson, err := json.Marshal(&ClusterInfo{
		Setup: true,
	})
	if err != nil {
		i.log.Warn("failed to marshall cluster info", zap.Error(err))
		return
	}

	_, err = inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyCluster,
		).String(),
		string(clusterJson),
	)
	if err != nil {
		i.log.Warn("failed to put cluster info", zap.Error(err))
		return
	}
}

func (i *Instance) startWatchRole(ctx context.Context, id string, startCallback func()) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		if e, ok := err.(error); ok {
			i.log.Error("recover in role", zap.String("roleId", id), zap.Error(e))
			sentry.CaptureException(e)
		} else {
			i.log.Error("recover in role", zap.String("roleId", id), zap.Any("panic", err))
		}
	}()
	// Load current config
	config, err := i.kv.Get(
		ctx,
		i.kv.Key(
			types.KeyInstance,
			types.KeyRole,
			id,
		).String(),
	)
	rawConfig := []byte{}
	if err == nil && len(config.Kvs) > 0 {
		rawConfig = config.Kvs[0].Value
	}
	i.startRole(ctx, id, rawConfig)
	startCallback()
	for resp := range i.kv.Watch(
		ctx,
		i.kv.Key(
			types.KeyInstance,
			types.KeyRole,
			id,
		).String(),
	) {
		for _, ev := range resp.Events {
			rawConfig := []byte{}
			if ev.Type != clientv3.EventTypeDelete && len(ev.Kv.Value) > 0 {
				rawConfig = ev.Kv.Value
			}
			i.log.Info("restarting role due to config change", zap.String("roleId", id), zap.String("key", string(ev.Kv.Key)))
			i.DispatchEvent(types.EventTopicRoleRestart, roles.NewEvent(
				ctx,
				map[string]interface{}{
					"id":     id,
					"config": rawConfig,
				},
			))
		}
	}
}

func (i *Instance) startRole(ctx context.Context, id string, rawConfig []byte) bool {
	srs := sentry.TransactionFromContext(ctx).StartChild("gravity.instance.role.start")
	srs.Description = id
	srs.SetTag("gravity.role", id)
	defer srs.Finish()
	defer i.putInstanceInfo(srs.Context())
	instanceRoleStarted.WithLabelValues(id).SetToCurrentTime()
	client := i.roles[id].RoleInstance.kv
	if mr, ok := i.roles[id].Role.(roles.MigratableRole); ok {
		mr.RegisterMigrations()
		// Run migrations
		_client, err := i.roles[id].RoleInstance.Migrator().Run(srs.Context())
		if err != nil {
			i.log.Warn("failed to run migrations for role", zap.String("roleId", id))
			return false
		}
		client = _client
	}
	// Overwrite role's KV client with the potentially hooked client for migrations
	i.roles[id].RoleInstance.kv = client
	// Start role
	err := i.roles[id].Role.Start(srs.Context(), rawConfig)
	if err == roles.ErrRoleNotConfigured {
		i.log.Info("role not configured", zap.String("roleId", id))
	} else if err != nil {
		i.log.Warn("failed to start role", zap.String("roleId", id), zap.Error(err))
		return false
	}
	i.log.Info("Started role", zap.String("roleId", id))
	i.DispatchEvent(types.EventTopicRoleStarted, roles.NewEvent(
		srs.Context(),
		map[string]interface{}{
			"role":   id,
			"config": rawConfig,
		},
	))
	return true
}

func (i *Instance) stopRole(ctx context.Context, cause error, id string) {
	srs := sentry.TransactionFromContext(ctx).StartChild("gravity.instance.role.stop")
	srs.Description = id
	srs.SetTag("gravity.role", id)
	defer srs.Finish()
	i.log.Info("stopping role", zap.String("roleId", id))
	i.roles[id].Role.Stop()
	// Cancel context and re-create the context
	i.roles[id].ContextCancelFunc(cause)
	ctx, cancel := context.WithCancelCause(i.rootContext)
	i.rolesM.Lock()
	i.roles[id] = RoleContext{
		Role:              i.roles[id].Role,
		RoleInstance:      i.ForRole(id, ctx),
		ContextCancelFunc: cancel,
	}
	i.rolesM.Unlock()
}

func (i *Instance) Stop() {
	i.log.Info("stopping")
	for id, role := range i.roles {
		i.log.Info("stopping role", zap.String("roleId", id))
		role.ContextCancelFunc(ErrInstanceStopping)
		role.Role.Stop()
	}
	if i.etcd != nil {
		i.log.Info("stopping role", zap.String("roleId", "etcd"))
		i.etcd.Stop()
	}
	i.rootContextCancel(ErrInstanceStopping)
	sentry.Flush(2 * time.Second)
}
