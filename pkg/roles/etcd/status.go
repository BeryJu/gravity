package etcd

import (
	"context"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/storage"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func IsLearnerReady(leaderStatus, learnerStatus *clientv3.StatusResponse) bool {
	leaderRev := leaderStatus.Header.Revision
	learnerRev := learnerStatus.Header.Revision

	learnerReadyPercent := float64(learnerRev) / float64(leaderRev)
	return learnerReadyPercent >= 0.9
}

type ClusterStatus struct {
	Healthy      error
	MemberStatus map[uint64]MemberStatus
}
type MemberStatus struct {
	Raw *clientv3.StatusResponse
}

func (r *Role) clusterStatus(ctx context.Context) (*ClusterStatus, error) {
	members, err := r.i.KV().MemberList(ctx, clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}
	cst := &ClusterStatus{
		MemberStatus: map[uint64]MemberStatus{},
	}
	for _, member := range members.Members {
		c := storage.NewClient(
			extconfig.Get().Etcd.Prefix,
			r.log.Named("etcd").Named(member.Name),
			extconfig.Get().Debug,
			member.ClientURLs...,
		)
		st, err := c.Status(ctx, member.ClientURLs[0])
		if err != nil {
			cst.Healthy = err
			continue
		}
		cst.MemberStatus[member.ID] = MemberStatus{
			Raw: st,
		}
	}
	return cst, nil
}

func (cst *ClusterStatus) FindLeaderStatus() (uint64, *clientv3.StatusResponse) {
	for i := range cst.MemberStatus {
		status := cst.MemberStatus[i].Raw
		if status.Leader == status.Header.MemberId {
			return status.Header.MemberId, status
		}
	}
	return 0, nil
}

func (cst *ClusterStatus) FindLearnerStatus() (uint64, *clientv3.StatusResponse) {
	for i := range cst.MemberStatus {
		status := cst.MemberStatus[i].Raw
		if status.IsLearner {
			return status.Header.MemberId, status
		}
	}
	return 0, nil
}

func (r *Role) clusterCanJoin(ctx context.Context) bool {
	st, err := r.clusterStatus(ctx)
	if err != nil {
		r.log.Warn("failed to check cluster status", zap.Error(err))
		return false
	}
	if st.Healthy != nil {
		r.log.Warn("cluster is not healthy", zap.Error(err))
		return false
	}
	_, lds := st.FindLeaderStatus()
	if id, st := st.FindLearnerStatus(); id > 0 {
		r.log.Info("Found learner")
		if IsLearnerReady(lds, st) {
			r.log.Info("Learner is ready to be promoted")
			r.i.KV().MemberPromote(ctx, id)
		} else {
			r.log.Info("Learner is not ready yet")
			return false
		}
	}
	return true
}
