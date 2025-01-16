package dns

import (
	"strings"

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

func init() {
	HandlerRegistry.Add(MemoryType, func(z *Zone, rawConfig map[string]interface{}) Handler {
		return NewMemoryHandler(z, rawConfig)
	})
}

func NewMemoryHandler(z *Zone, config map[string]interface{}) *MemoryHandler {
	mh := &MemoryHandler{
		EtcdHandler: &EtcdHandler{z: z},
		z:           z,
	}
	mh.lookupKey = func(k *storage.Key, qname string, r *utils.DNSRequest) []dns.RR {
		answers := []dns.RR{}
		var recs map[string]*Record = make(map[string]*Record)
		var ok bool
		if k.IsPrefix() {
			prefix := k.String()
			// If the key is a prefix, we can't just directly look it up in the map,
			// and have to fall back to a "slightly" slower method of iterating through the map
			for key, rr := range mh.z.records.Iter() {
				if !strings.HasPrefix(key, prefix) {
					continue
				}
				for ikey, r := range rr {
					recs[ikey] = r
				}
			}
		} else {
			recs, ok = mh.z.records.GetPrefix(k.String())
			if !ok {
				return answers
			}
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
