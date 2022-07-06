package instance

import (
	"context"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/beryju/dns-dhcp-etcd-thingy/internal/extconfig"
	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles"
	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles/api"
	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles/dhcp"
	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles/dns"
	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles/etcd"
	"github.com/beryju/dns-dhcp-etcd-thingy/internal/storage"
)

const (
	InstancePrefix = "instance"
)

func NewInstance() {
	extCfg := extconfig.Get()
	instance := Instance{
		log:        log.WithField("instance", extCfg.InstanceIdentifier),
		roles:      make(map[string]roles.Role),
		identifier: extCfg.InstanceIdentifier,
	}
	defer instance.Stop()
	if strings.Contains(extCfg.BootstrapRoles, "etcd") {
		instance.log.Info("'etcd' in bootstrap roles, starting embedded etcd")
		// TODO: join existing cluster?
		instance.etcd = etcd.New(extCfg.Etcd.Prefix)
		err := instance.etcd.Start(func() {
			instance.bootstrap()
		})
		if err != nil {
			instance.log.WithError(err).Warning("failed to start etcd")
		}
	} else {
		instance.bootstrap()
	}
}

type Instance struct {
	roles      map[string]roles.Role
	kv         *storage.Client
	log        *log.Entry
	identifier string

	etcd *etcd.EmbeddedEtcd
}

func (i *Instance) getRoles() []string {
	rr, err := i.kv.Get(context.TODO(), i.kv.Key(InstancePrefix, i.identifier, "roles"))
	roles := extconfig.Get().BootstrapRoles
	if err == nil && len(rr.Kvs) > 0 {
		roles = rr.Kvs[0].String()
	} else {
		i.log.Info("defaulting to bootstrap roles")
	}
	return strings.Split(roles, ";")
}

func (i *Instance) bootstrap() {
	i.log.Debug("bootstrapping instance")
	i.kv = storage.NewClient(
		extconfig.Get().Etcd.Endpoint,
		extconfig.Get().Etcd.Prefix,
	)
	for _, roleId := range i.getRoles() {
		switch roleId {
		case "dhcp":
			i.roles[roleId] = dhcp.New(i)
		case "dns":
			i.roles[roleId] = dns.New(i)
		case "api":
			i.roles[roleId] = api.New(i)
		case "etcd":
			// Special case
			continue
		default:
			i.log.WithField("roleId", roleId).Info("Invalid role, skipping")
		}
	}
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
	config, err := i.kv.Get(context.TODO(), i.kv.Key(InstancePrefix, "role", id))
	rawConfig := []byte{}
	if err == nil && len(config.Kvs) > 0 {
		rawConfig = config.Kvs[0].Value
	}
	role.Start(rawConfig)
	wg.Done()
}

func (i *Instance) Stop() {
	if i.etcd != nil {
		i.log.WithField("role", "embedded-etcd")
		i.etcd.Stop()
	}
}
