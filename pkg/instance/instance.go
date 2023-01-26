package instance

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/backup"
	"beryju.io/gravity/pkg/roles/debug"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/etcd"
	"beryju.io/gravity/pkg/roles/monitoring"
	"beryju.io/gravity/pkg/roles/tsdb"
	"beryju.io/gravity/pkg/storage"
)

type ClusterInfo struct {
	Setup bool `json:"setup"`
}

type RoleContext struct {
	Role              roles.Role
	Context           context.Context
	ContextCancelFunc context.CancelFunc
}

type Instance struct {
	roles      map[string]RoleContext
	rolesM     sync.Mutex
	kv         *storage.Client
	log        *zap.Logger
	identifier string

	eventHandlers  map[string]map[string][]roles.EventHandler
	eventHandlersM sync.RWMutex

	etcd *etcd.Role

	rootContext       context.Context
	rootContextCancel context.CancelFunc

	instanceInfoLease *clientv3.LeaseID
}

func New() *Instance {
	extCfg := extconfig.Get()
	ctx, canc := context.WithCancel(context.Background())
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
	go i.startSentry()
	if strings.Contains(extconfig.Get().BootstrapRoles, "etcd") {
		i.log.Info("'etcd' in bootstrap roles, starting embedded etcd")
		i.etcd = etcd.New(i.ForRole("etcd"))
		i.etcd.Start(i.rootContext)
	}
	i.bootstrap()
}

func (i *Instance) Log() *zap.Logger {
	return i.log
}

func (i *Instance) getRoles() []string {
	rr, err := i.kv.Get(
		i.rootContext,
		i.kv.Key(
			types.KeyInstance,
			i.identifier,
			types.KeyRoles,
		).String(),
	)
	roles := extconfig.Get().BootstrapRoles
	if err == nil && len(rr.Kvs) > 0 {
		roles = string(rr.Kvs[0].Value)
		i.log.Info("roles configured for instance", zap.Strings("roles", strings.Split(roles, ";")))
	} else {
		i.log.Info("defaulting to bootstrap roles", zap.Strings("roles", strings.Split(roles, ";")))
	}
	return strings.Split(roles, ";")
}

func (i *Instance) bootstrap() {
	i.log.Debug("bootstrapping instance")
	i.keepAliveInstanceInfo()
	i.putInstanceInfo()
	i.setupInstanceAPI()
	for _, roleId := range i.getRoles() {
		instanceRoles.WithLabelValues(roleId).Add(1)
		ctx, cancel := context.WithCancel(i.rootContext)
		rc := RoleContext{
			Context:           ctx,
			ContextCancelFunc: cancel,
		}
		roleInst := i.ForRole(roleId)
		switch roleId {
		case "dhcp":
			rc.Role = dhcp.New(roleInst)
		case "dns":
			rc.Role = dns.New(roleInst)
		case "api":
			rc.Role = api.New(roleInst)
		case "discovery":
			rc.Role = discovery.New(roleInst)
		case "backup":
			rc.Role = backup.New(roleInst)
		case "monitoring":
			rc.Role = monitoring.New(roleInst)
		case "debug":
			rc.Role = debug.New(roleInst)
		case "tsdb":
			rc.Role = tsdb.New(roleInst)
		case "etcd":
			// Special case
			continue
		default:
			i.log.Info("Invalid role, skipping", zap.String("roleId", roleId))
			continue
		}
		i.rolesM.Lock()
		i.roles[roleId] = rc
		i.rolesM.Unlock()
	}
	i.ForRole("root").AddEventListener(types.EventTopicRoleRestart, i.eventRoleRestart)
	i.ForRole("root").DispatchEvent(
		types.EventTopicInstanceBootstrapped,
		roles.NewEvent(i.rootContext, map[string]interface{}{}),
	)
	i.checkFirstStart()
	for roleId := range i.roles {
		go i.startWatchRole(roleId)
	}
	<-i.rootContext.Done()
}

func (i *Instance) eventRoleRestart(ev *roles.Event) {
	id := ev.Payload.Data["id"].(string)
	config := ev.Payload.Data["config"].([]byte)
	i.stopRole(id)
	i.startRole(id, config)
}

func (i *Instance) checkFirstStart() {
	inst := i.ForRole("root")
	cluster, err := inst.KV().Get(
		i.rootContext,
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
	inst.DispatchEvent(
		types.EventTopicInstanceFirstStart,
		roles.NewEvent(i.rootContext, map[string]interface{}{}),
	)

	clusterJson, err := json.Marshal(&ClusterInfo{
		Setup: true,
	})
	if err != nil {
		i.log.Warn("failed to marshall cluster info", zap.Error(err))
		return
	}

	_, err = inst.KV().Put(
		i.rootContext,
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

func (i *Instance) startWatchRole(id string) {
	defer func() {
		err := extconfig.RecoverWrapper(recover())
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
		i.roles[id].Context,
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
	i.startRole(id, rawConfig)
	for resp := range i.kv.Watch(
		i.rootContext,
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
			i.log.Info("stopping role due to config change", zap.String("roleId", id), zap.String("key", string(ev.Kv.Key)))
			i.DispatchEvent(types.EventTopicRoleRestart, roles.NewEvent(
				i.rootContext,
				map[string]interface{}{
					"id":     id,
					"config": rawConfig,
				},
			))
		}
	}
}

func (i *Instance) startRole(id string, rawConfig []byte) bool {
	defer i.putInstanceInfo()
	instanceRoleStarted.WithLabelValues(id).SetToCurrentTime()
	err := i.roles[id].Role.Start(i.roles[id].Context, rawConfig)
	if err == roles.ErrRoleNotConfigured {
		i.log.Info("role not configured", zap.String("roleId", id))
	} else if err != nil {
		i.log.Warn("failed to start role", zap.String("roleId", id), zap.Error(err))
		return false
	}
	i.log.Debug("started role", zap.String("roleId", id))
	i.DispatchEvent(types.EventTopicRoleStarted, roles.NewEvent(
		i.rootContext,
		map[string]interface{}{
			"role":   id,
			"config": rawConfig,
		},
	))
	return true
}

func (i *Instance) stopRole(id string) {
	i.log.Info("stopping role", zap.String("roleId", id))
	i.roles[id].Role.Stop()
	// Cancel context and re-create the context
	i.roles[id].ContextCancelFunc()
	ctx, cancel := context.WithCancel(i.rootContext)
	i.rolesM.Lock()
	i.roles[id] = RoleContext{
		Role:              i.roles[id].Role,
		Context:           ctx,
		ContextCancelFunc: cancel,
	}
	i.rolesM.Unlock()
}

func (i *Instance) Stop() {
	i.log.Info("stopping")
	for id, role := range i.roles {
		i.log.Info("stopping role", zap.String("roleId", id))
		role.ContextCancelFunc()
		role.Role.Stop()
	}
	if i.etcd != nil {
		i.log.Info("stopping role", zap.String("roleId", "etcd"))
		i.etcd.Stop()
	}
	i.rootContextCancel()
	sentry.Flush(2 * time.Second)
}
