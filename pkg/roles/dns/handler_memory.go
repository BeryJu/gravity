package dns

import (
	"beryju.io/gravity/pkg/roles/dns/utils"
	"beryju.io/gravity/pkg/storage"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

const MemoryType = "memory"

type MemoryHandler struct {
	*EtcdHandler
	log *zap.Logger
	z   *Zone
}

func NewMemoryHandler(z *Zone, config map[string]string) *MemoryHandler {
	mh := &MemoryHandler{
		EtcdHandler: &EtcdHandler{z: z},
		z:           z,
	}
	mh.lookupKey = func(k *storage.Key, qname string, r *utils.DNSRequest) []dns.RR {
		mh.z.recordsSync.RLock()
		recs, ok := mh.z.records[k.String()]
		mh.z.recordsSync.RUnlock()
		answers := []dns.RR{}
		if !ok {
			return answers
		}
		for _, rec := range recs {
			ans := rec.ToDNS(qname)
			if ans != nil {
				answers = append(answers, ans)
			}
		}
		return answers
	}
	mh.log = z.log.With(zap.String("handler", mh.Identifier()))
	return mh
}

func (mh *MemoryHandler) Identifier() string {
	return MemoryType
}
