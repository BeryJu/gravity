package etcd

import (
	"context"
	"sync"
	"time"

	"beryju.io/gravity/pkg/roles"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.uber.org/zap"
)

type LeaderClusterReconciler struct {
	i      roles.Instance
	log    *zap.Logger
	ctx    context.Context
	cancel context.CancelFunc
	st     map[uint64]struct{}
	mutex  sync.Mutex
}

func NewLeaderClusterConciler(i roles.Instance) *LeaderClusterReconciler {
	return &LeaderClusterReconciler{
		i:     i,
		log:   i.Log(),
		mutex: sync.Mutex{},
	}
}

func (lcr *LeaderClusterReconciler) Start() {
	lcr.mutex.Lock()
	defer lcr.mutex.Unlock()
	lcr.ctx, lcr.cancel = context.WithCancel(context.Background())
	lcr.st = map[uint64]struct{}{}
	for {
		select {
		case <-lcr.ctx.Done():
			return
		default:
		}
		resp, err := lcr.i.KV().MemberList(lcr.ctx)
		if err != nil {
			lcr.log.Warn("member list error", zap.Error(err))
			time.Sleep(3 * time.Second)
			continue
		}
		current := map[uint64]*etcdserverpb.Member{}
		for _, m := range resp.Members {
			current[m.ID] = m
			if _, ok := lcr.st[m.ID]; !ok {
				lcr.MemberAdded(m)
			}
		}
		for id := range lcr.st {
			if _, ok := current[id]; !ok {
				lcr.MemberRemoved(id)
			}
		}
		lcr.st = map[uint64]struct{}{}
		for id := range current {
			lcr.st[id] = struct{}{}
		}
		lcr.Reconcile()

		time.Sleep(5 * time.Second)
	}
}

func (lcr *LeaderClusterReconciler) Stop() {
	lcr.mutex.TryLock()
	if lcr.cancel != nil {
		lcr.cancel()
	}
}

func (lcr *LeaderClusterReconciler) MemberAdded(m *etcdserverpb.Member) {
	lcr.log.Debug("New member added", zap.Uint64("id", m.ID), zap.String("name", m.Name))
}

func (lcr *LeaderClusterReconciler) MemberRemoved(id uint64) {
	lcr.log.Debug("Member removed", zap.Uint64("id", id))
}

func (lcr *LeaderClusterReconciler) Reconcile() {
	lcr.log.Debug("Reconciling cluster state")
	st, err := lcr.ClusterStatus(lcr.ctx)
	if err != nil {
		lcr.log.Warn("failed to check cluster status", zap.Error(err))
		return
	}
	if st.Healthy != nil {
		lcr.log.Warn("cluster is not healthy", zap.Error(err))
		return
	}
	_, lds := st.FindLeaderStatus()
	if id, st := st.FindLearnerStatus(); id > 0 {
		lcr.log.Info("Found learner")
		if IsLearnerReady(lds, st) {
			lcr.log.Info("Learner is ready to be promoted")
			_, err := lcr.i.KV().MemberPromote(lcr.ctx, id)
			if err != nil {
				lcr.log.Info("Failed to promote member", zap.Error(err))
				return
			}
		} else {
			lcr.log.Info("Learner is not ready yet")
		}
		return
	}
}
