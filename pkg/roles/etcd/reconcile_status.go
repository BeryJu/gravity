package etcd

import (
	"context"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/storage"
	clientv3 "go.etcd.io/etcd/client/v3"
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

func (lcr *LeaderClusterReconciler) ClusterStatus(ctx context.Context) (*ClusterStatus, error) {
	members, err := lcr.i.KV().MemberList(ctx, clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}
	cst := &ClusterStatus{
		MemberStatus: map[uint64]MemberStatus{},
	}
	for _, member := range members.Members {
		c := storage.NewClient(
			extconfig.Get().Etcd.Prefix,
			lcr.log.Named("etcd").Named(member.Name),
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
