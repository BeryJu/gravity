package dhcp

import (
	"context"
	"encoding/json"
	"strings"

	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/dhcp/types"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Lease struct {
	Identifier string `json:"-"`

	Address          string `json:"address"`
	Hostname         string `json:"hostname"`
	AddressLeaseTime string `json:"addressLeaseTime"`
	ScopeKey         string `json:"scopeKey"`

	scope   *Scope
	etcdKey string
	inst    roles.Instance
	log     *log.Entry
}

func (r *DHCPRole) leaseFromKV(raw *mvccpb.KeyValue) (*Lease, error) {
	s := &Lease{
		inst: r.i,
	}
	err := json.Unmarshal(raw.Value, &s)
	if err != nil {
		return nil, err
	}
	prefix := r.i.GetKV().Key(
		types.KeyRole,
		types.KeyScopes,
		// l.Scope.Name,
		types.KeyLeases,
		"",
	)
	s.Identifier = strings.TrimPrefix(string(raw.Key), prefix)
	// Get full etcd key without leading slash since this usually gets passed to Instance Key()
	s.etcdKey = string(raw.Key)[1:]

	s.log = log.WithField("lease", prefix)
	return s, nil
}

func (l *Lease) put(expiry int64) error {
	raw, err := json.Marshal(&l)
	if err != nil {
		return err
	}

	exp, err := l.inst.GetKV().Lease.Grant(context.TODO(), expiry)
	if err != nil {
		return err
	}

	leaseKey := l.inst.GetKV().Key(
		types.KeyRole,
		types.KeyScopes,
		l.scope.Name,
		types.KeyLeases,
		l.Identifier,
	)
	_, err = l.inst.GetKV().Put(
		context.TODO(),
		leaseKey,
		string(raw),
		clientv3.WithLease(exp.ID),
	)
	if err != nil {
		return err
	}
	ev := roles.NewEvent(
		map[string]interface{}{
			"hostname": l.Hostname,
			"address":  l.Address,
		},
	)
	ev.Payload.RelatedObjectKey = leaseKey
	ev.Payload.RelatedObjectOptions = []clientv3.OpOption{clientv3.WithLease(exp.ID)}
	l.inst.DispatchEvent(types.EventTopicDHCPLeaseGiven, ev)

	l.log.WithField("expiry", expiry).Debug("put lease")
	return nil
}
