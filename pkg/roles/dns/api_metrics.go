package dns

import (
	"context"
	"strconv"
	"strings"

	"beryju.io/gravity/pkg/roles/dns/types"
	tsdbTypes "beryju.io/gravity/pkg/roles/tsdb/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func (r *Role) APIMetrics() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *tsdbTypes.APIMetricsGetOutput) error {
		prefix := r.i.KV().Key(
			tsdbTypes.KeyRole,
			types.KeyRole,
		).Prefix(true).String()
		rawMetrics, err := r.i.KV().Get(
			ctx,
			prefix,
			clientv3.WithPrefix(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, kv := range rawMetrics.Kvs {
			value, err := strconv.ParseInt(string(kv.Value), 10, 0)
			if err != nil {
				r.log.Warn("failed to parse metric value", zap.Error(err))
				continue
			}
			keyParts := strings.Split(strings.TrimPrefix(string(kv.Key), prefix), "/")
			output.Records = append(output.Records, tsdbTypes.APIMetricsRecord{
				Time:    keyParts[2],
				Handler: keyParts[0],
				Node:    keyParts[1],
				Value:   value,
			})
		}
		return nil
	})
	u.SetName("dns.get_metrics")
	u.SetTitle("DNS Metrics")
	u.SetTags("roles/dns")
	u.SetExpectedErrors(status.Internal)
	return u
}
