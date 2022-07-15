package instance

import (
	"context"
	"encoding/base64"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ExportEntry struct {
	Key   string
	Value string
}

func (i *Instance) Export() ([]ExportEntry, error) {
	exps, err := i.kv.Get(context.Background(), "/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	entries := []ExportEntry{}
	for idx, exp := range exps.Kvs {
		entries[idx] = ExportEntry{
			Key:   strings.TrimPrefix(string(exp.Key), extconfig.Get().Etcd.Prefix),
			Value: base64.StdEncoding.EncodeToString(exp.Value),
		}
	}
	return entries, nil
}

func (i *Instance) Import(entries []ExportEntry) error {
	for _, entry := range entries {
		val, err := base64.StdEncoding.DecodeString(entry.Value)
		if err != nil {
			i.log.WithField("key", entry.Key).WithError(err).Warning("failed to decode value")
			continue
		}
		_, err = i.kv.Put(context.Background(), entry.Key, string(val))
		if err != nil {
			i.log.WithField("key", entry.Key).WithError(err).Warning("failed to put value")
			continue
		}
	}
	return nil
}
