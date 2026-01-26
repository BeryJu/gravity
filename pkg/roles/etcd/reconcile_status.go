package etcd

import (
	"context"
	"net/url"

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
	MemberStatus map[uint64]*clientv3.StatusResponse
}

func FindNonLocalhost(addrs []string) (string, error) {
	var lastError error
	for _, addr := range addrs {
		url, err := url.Parse(addr)
		if err != nil {
			lastError = err
			continue
		}
		if url.Hostname() == "localhost" {
			continue
		}
		return url.String(), nil
	}
	return "", lastError
}

func (lcr *LeaderClusterReconciler) ClusterStatus(ctx context.Context) (*ClusterStatus, error) {
	members, err := lcr.i.KV().MemberList(ctx, clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}
	cst := &ClusterStatus{
		MemberStatus: map[uint64]*clientv3.StatusResponse{},
	}
	for _, member := range members.Members {
		nonLocalhost, err := FindNonLocalhost(member.ClientURLs)
		if err != nil {
			lcr.log.Warn("failed to get member IP for non-localhost", zap.Error(err))
			continue
		}
		c := storage.NewClient(
			extconfig.Get().Etcd.Prefix,
			lcr.log.Named(member.Name),
			extconfig.Get().Debug,
			nonLocalhost,
		)
		defer func() {
			err := c.Close()
			if err != nil {
				lcr.log.Warn("failed to close etcd client", zap.Error(err))
			}
		}()
		st, err := c.Status(ctx, nonLocalhost)
		if err != nil {
			cst.Healthy = err
			continue
		}
		cst.MemberStatus[member.ID] = st
	}
	return cst, nil
}

func (cst *ClusterStatus) FindLeaderStatus() *clientv3.StatusResponse {
	for _, st := range cst.MemberStatus {
		if st.Leader == st.Header.MemberId {
			return st
		}
	}
	return nil
}

func (cst *ClusterStatus) FindLearnerStatus() *clientv3.StatusResponse {
	for _, st := range cst.MemberStatus {
		if st.IsLearner {
			return st
		}
	}
	return nil
}
