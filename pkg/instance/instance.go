package instance

import (
	"context"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/api"
	"beryju.io/ddet/pkg/roles/dhcp"
	"beryju.io/ddet/pkg/roles/discovery"
	"beryju.io/ddet/pkg/roles/dns"
	"beryju.io/ddet/pkg/roles/etcd"
	"beryju.io/ddet/pkg/storage"
)

const (
	KeyInstance = "instance"
	KeyRole     = "role"
)

const (
	EventTopicInstanceBootstrapped = "instance.root.bootstrapped"
)

type Instance struct {
	roles      map[string]roles.Role
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
		roles:          make(map[string]roles.Role),
		identifier:     extCfg.Instance.Identifier,
		eventHandlersM: sync.RWMutex{},
		eventHandlers:  make(map[string]map[string][]roles.EventHandler),
	}
}

func (i *Instance) Start() {
	if strings.Contains(extconfig.Get().BootstrapRoles, "etcd") {
		i.log.Info("'etcd' in bootstrap roles, starting embedded etcd")
		// TODO: join existing cluster?
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

func (i *Instance) GetLogger() *log.Entry {
	return i.log
}

func (i *Instance) bootstrap() {
	i.log.Trace("bootstrapping instance")
	i.kv = storage.NewClient(
		extconfig.Get().Etcd.Endpoint,
		extconfig.Get().Etcd.Prefix,
	)
	for _, roleId := range i.getRoles() {
		roleInst := i.ForRole(roleId)
		switch roleId {
		case "dhcp":
			i.roles[roleId] = dhcp.New(roleInst)
		case "dns":
			i.roles[roleId] = dns.New(roleInst)
		case "api":
			i.roles[roleId] = api.New(roleInst)
		case "discovery":
			i.roles[roleId] = discovery.New(roleInst)
		case "etcd":
			// Special case
			continue
		default:
			i.log.WithField("roleId", roleId).Info("Invalid role, skipping")
		}
	}
	i.ForRole("root").DispatchEvent(EventTopicInstanceBootstrapped, roles.NewEvent(map[string]interface{}{}))
	wg := sync.WaitGroup{}
	for roleId := range i.roles {
		wg.Add(1)
		go i.startRole(roleId, wg)
	}
	wg.Wait()
}

func (i *Instance) startRole(id string, wg sync.WaitGroup) {
	i.log.WithField("roleId", id).Info("starting role")
	role := i.roles[id]
	config, err := i.kv.Get(context.TODO(), i.kv.Key(KeyInstance, KeyRole, id))
	rawConfig := []byte{}
	if err == nil && len(config.Kvs) > 0 {
		rawConfig = config.Kvs[0].Value
	}
	role.Start(rawConfig)
	wg.Done()
}

func (i *Instance) Stop() {
	i.log.Info("Stopping")
	for _, role := range i.roles {
		go role.Stop()
	}
	if i.etcd != nil {
		i.etcd.Stop()
	}
}
