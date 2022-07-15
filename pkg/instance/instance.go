package instance

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/backup"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/etcd"
	"beryju.io/gravity/pkg/roles/monitoring"
	"beryju.io/gravity/pkg/storage"
)

const (
	KeyInstance = "instance"
	KeyRole     = "role"
)

const (
	EventTopicInstanceBootstrapped = "instance.root.bootstrapped"
)

type RoleContext struct {
	Role              roles.Role
	Context           context.Context
	ContextCancelFunc context.CancelFunc
}

type Instance struct {
	roles      map[string]RoleContext
	kv         *storage.Client
	log        *log.Entry
	identifier string

	eventHandlersM sync.RWMutex
	eventHandlers  map[string]map[string][]roles.EventHandler

	etcd *etcd.Role
}

func NewInstance() *Instance {
	extCfg := extconfig.Get()
	return &Instance{
		log:            log.WithField("instance", extCfg.Instance.Identifier),
		roles:          make(map[string]RoleContext),
		identifier:     extCfg.Instance.Identifier,
		eventHandlersM: sync.RWMutex{},
		eventHandlers:  make(map[string]map[string][]roles.EventHandler),
		kv:             extCfg.EtcdClient(),
	}
}

func (i *Instance) Start() {
	i.log.WithField("version", extconfig.FullVersion()).Info("Gravity starting")
	go i.startSentry()
	if strings.Contains(extconfig.Get().BootstrapRoles, "etcd") {
		i.log.Info("'etcd' in bootstrap roles, starting embedded etcd")
		i.etcd = etcd.New(i.ForRole("etcd"))
		err := i.etcd.Start(func() {
			i.bootstrap()
		})
		if err != nil {
			i.log.WithError(err).Warning("failed to start etcd")
		}
	} else {
		i.bootstrap()
	}
}

func (i *Instance) startSentry() {
	transport := sentry.NewHTTPTransport()
	transport.Configure(sentry.ClientOptions{
		HTTPTransport: extconfig.Transport(),
	})
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://a6690a50e8924263bd6f82fe3a1a2386@sentry.beryju.org/17",
		Environment:      "",
		Release:          fmt.Sprintf("gravity@%s", extconfig.FullVersion()),
		TracesSampleRate: 1.0,
		Transport:        transport,
	})
	if err != nil {
		i.log.WithError(err).Warning("failed to init sentry")
	}
}

func (i *Instance) Log() *log.Entry {
	return i.log
}

func (i *Instance) getRoles() []string {
	rr, err := i.kv.Get(context.TODO(), i.kv.Key(KeyInstance, i.identifier, "roles").String())
	roles := extconfig.Get().BootstrapRoles
	if err == nil && len(rr.Kvs) > 0 {
		roles = rr.Kvs[0].String()
	} else {
		i.log.Info("defaulting to bootstrap roles")
	}
	return strings.Split(roles, ";")
}

func (i *Instance) bootstrap() {
	i.log.Trace("bootstrapping instance")
	for _, roleId := range i.getRoles() {
		instanceRoles.WithLabelValues(roleId).Add(1)
		ctx, cancel := context.WithCancel(context.Background())
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
		case "etcd":
			// Special case
			continue
		default:
			i.log.WithField("roleId", roleId).Info("Invalid role, skipping")
			continue
		}
		i.roles[roleId] = rc
	}
	i.ForRole("root").DispatchEvent(EventTopicInstanceBootstrapped, roles.NewEvent(map[string]interface{}{}))
	wg := sync.WaitGroup{}
	for roleId := range i.roles {
		wg.Add(1)
		go func(id string) {
			defer func() {
				err := recover()
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
			i.log.WithField("roleId", id).Info("starting role")
			role := i.roles[id]
			config, err := i.kv.Get(context.TODO(), i.kv.Key(KeyInstance, KeyRole, id).String())
			rawConfig := []byte{}
			if err == nil && len(config.Kvs) > 0 {
				rawConfig = config.Kvs[0].Value
			}
			role.Role.Start(role.Context, rawConfig)
			wg.Done()
		}(roleId)
	}
	wg.Wait()
}

func (i *Instance) Stop() {
	i.log.Info("Stopping")
	for _, role := range i.roles {
		role.ContextCancelFunc()
		go role.Role.Stop()
	}
	if i.etcd != nil {
		i.etcd.Stop()
	}
	sentry.Flush(2 * time.Second)
}
