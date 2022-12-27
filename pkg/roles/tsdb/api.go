package tsdb

import (
	"context"
	"strconv"
	"strings"

	"beryju.io/gravity/pkg/roles/tsdb/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func (r *Role) APIMetricsMemory() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *types.APIMetricsGetOutput) error {
		prefix := r.i.KV().Key(
			types.KeyRole,
			types.KeySystem,
			"memory",
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
			output.Records = append(output.Records, types.APIMetricsRecord{
				Time:  keyParts[1],
				Node:  keyParts[0],
				Value: value,
			})
		}
		return nil
	})
	u.SetName("api.get_metrics_memory")
	u.SetTitle("System Metrics")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}

func (r *Role) APIMetricsCPU() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *types.APIMetricsGetOutput) error {
		prefix := r.i.KV().Key(
			types.KeyRole,
			types.KeySystem,
			"cpu",
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
			output.Records = append(output.Records, types.APIMetricsRecord{
				Time:  keyParts[1],
				Node:  keyParts[0],
				Value: value,
			})
		}
		return nil
	})
	u.SetName("api.get_metrics_cpu")
	u.SetTitle("System Metrics")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}
