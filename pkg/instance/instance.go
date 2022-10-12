package instance

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"

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

type RoleContext struct {
	Role              roles.Role
	Context           context.Context
	ContextCancelFunc context.CancelFunc
}

type Instance struct {
	roles      map[string]RoleContext
	rolesM     sync.Mutex
	kv         *storage.Client
	log        *log.Entry
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
	return &Instance{
		roles:             make(map[string]RoleContext),
		rolesM:            sync.Mutex{},
		log:               log.WithField("instance", extCfg.Instance.Identifier).WithField("forRole", "root"),
		identifier:        extCfg.Instance.Identifier,
		eventHandlers:     make(map[string]map[string][]roles.EventHandler),
		eventHandlersM:    sync.RWMutex{},
		kv:                extCfg.EtcdClient(),
		rootContext:       ctx,
		rootContextCancel: canc,
	}
}

func (i *Instance) Start() {
	i.log.WithField("version", extconfig.FullVersion()).Info("Gravity starting")
	go i.startSentry()
	if strings.Contains(extconfig.Get().BootstrapRoles, "etcd") {
		i.log.Info("'etcd' in bootstrap roles, starting embedded etcd")
		i.etcd = etcd.New(i.ForRole("etcd"))
		i.etcd.Start(i.rootContext, []byte{})
	}
	i.bootstrap()
}

func (i *Instance) startSentry() {
	transport := sentry.NewHTTPTransport()
	transport.Configure(sentry.ClientOptions{
		HTTPTransport: extconfig.Transport(),
	})
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://ccd520a9a2b8458ca1e82108a8afb801@sentry.beryju.org/17",
		Environment:      "",
		Release:          fmt.Sprintf("gravity@%s", extconfig.FullVersion()),
		TracesSampleRate: 0.5,
		Transport:        transport,
		Debug:            extconfig.Get().Debug,
		DebugWriter:      NewSentryWriter(i.log.WithField("forRole", "sentry")),
	})
	if err != nil {
		i.log.WithError(err).Warning("failed to init sentry")
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("gravity.instance", extconfig.Get().Instance.Identifier)
		scope.SetTag("gravity.version", extconfig.Version)
		scope.SetTag("gravity.hash", extconfig.BuildHash)
	})
}

func (i *Instance) Log() *log.Entry {
	return i.log
}

func (i *Instance) getRoles() []string {
	rr, err := i.kv.Get(
		i.rootContext,
		i.kv.Key(
			types.KeyInstance,
			i.identifier,
			"roles",
		).String(),
	)
	roles := extconfig.Get().BootstrapRoles
	if err == nil && len(rr.Kvs) > 0 {
		roles = rr.Kvs[0].String()
	} else {
		i.log.WithField("roles", roles).Info("defaulting to bootstrap roles")
	}
	return strings.Split(roles, ";")
}

func (i *Instance) bootstrap() {
	i.log.Trace("bootstrapping instance")
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
			i.log.WithField("roleId", roleId).Info("Invalid role, skipping")
			continue
		}
		i.rolesM.Lock()
		i.roles[roleId] = rc
		i.rolesM.Unlock()
	}
	i.ForRole("root").DispatchEvent(
		types.EventTopicInstanceBootstrapped,
		roles.NewEvent(i.rootContext, map[string]interface{}{}),
	)
	for roleId := range i.roles {
		go i.startWatchRole(roleId)
	}
	<-i.rootContext.Done()
}

func (i *Instance) startWatchRole(id string) {
	defer func() {
		err := extconfig.RecoverWrapper(recover())
		if err == nil {
			return
		}
		if e, ok := err.(error); ok {
			i.log.WithError(e).Warning("recover in role")
			sentry.CaptureException(e)
		} else {
			i.log.WithField("panic", err).Warning("recover in role")
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
	started := i.startRole(id, rawConfig)
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
			if started {
				i.log.WithField("roleId", id).WithField("key", string(ev.Kv.Key)).Info("stopping role due to config change")
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
			started = i.startRole(id, rawConfig)
		}
	}
}

func (i *Instance) startRole(id string, rawConfig []byte) bool {
	defer i.putInstanceInfo()
	instanceRoleStarted.WithLabelValues(id).SetToCurrentTime()
	err := i.roles[id].Role.Start(i.roles[id].Context, rawConfig)
	if err == roles.ErrRoleNotConfigured {
		i.log.WithField("roleId", id).Info("role not configured")
	} else if err != nil {
		i.log.WithField("roleId", id).WithError(err).Warning("failed to start role")
		return false
	}
	i.log.WithField("roleId", id).Debug("started role")
	return true
}

func (i *Instance) Stop() {
	i.log.Info("Stopping")
	for id, role := range i.roles {
		i.log.WithField("roleId", id).Debug("stopping role")
		role.ContextCancelFunc()
		role.Role.Stop()
	}
	if i.etcd != nil {
		i.etcd.Stop()
	}
	i.rootContextCancel()
	sentry.Flush(2 * time.Second)
}
