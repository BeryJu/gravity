package handlers

import (
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.uber.org/zap"
)

type Handler interface {
	Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg
	Identifier() string
}

type HandlerZoneContext interface {
	Log() *zap.Logger
	// // TODO Rename to Zone()
	// GetZone() *types.Zone
	RecordFromKV(kv *mvccpb.KeyValue) (HandlerRecord, error)
	RoleInstance() roles.Instance
	EtcdKey() string
}

type HandlerRecord interface {
	ToDNS(string) dns.RR
}
