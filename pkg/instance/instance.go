package instance

import (
	"context"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/backup"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/etcd"
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

	etcd *etcd.EmbeddedEtcd
}

func NewInstance() *Instance {
	extCfg := extconfig.Get()
	return &Instance{
		log:            log.WithField("instance", extCfg.Instance.Identifier),
		roles:          make(map[string]RoleContext),
		identifier:     extCfg.Instance.Identifier,
		eventHandlersM: sync.RWMutex{},
		eventHandlers:  make(map[string]map[string][]roles.EventHandler),
	}
}

func (i *Instance) Start() {
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

func (i *Instance) Log() *log.Entry {
	return i.log
}

func (i *Instance) getRoles() []string {
	rr, err := i.kv.Get(context.TODO(), i.kv.Key(KeyInstance, i.identifier, "roles"))
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
	i.kv = storage.NewClient(
		extconfig.Get().Etcd.Endpoint,
		extconfig.Get().Etcd.Prefix,
	)
	for _, roleId := range i.getRoles() {
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
			i.log.WithField("roleId", id).Info("starting role")
			role := i.roles[id]
			config, err := i.kv.Get(context.TODO(), i.kv.Key(KeyInstance, KeyRole, id))
			rawConfig := []byte{}
			if err == nil && len(config.Kvs) > 0 {
				rawConfig = config.Kvs[0].Value
			}
			role.Role.Start(role.Context, rawConfig)
			go func() {
				err := recover()
				if err != nil {
					i.log.WithField("roleId", id).WithError(err.(error)).Error("Panic in role")
				}
			}()
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
}
