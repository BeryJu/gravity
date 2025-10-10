package etcd

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/backup"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"go.uber.org/zap"
)

type APIMember struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	IsLearner bool   `json:"isLearner"`
	IsLeader  bool   `json:"isLeader"`
}
type APIMembersOutput struct {
	Members []APIMember `json:"members"`
}

func (r *Role) APIClusterMembers() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIMembersOutput) error {
		members, err := r.i.KV().MemberList(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, mem := range members.Members {
			output.Members = append(output.Members, APIMember{
				ID:        fmt.Sprintf("%x", mem.ID),
				Name:      mem.Name,
				IsLearner: mem.IsLearner,
				IsLeader:  mem.ID == r.e.Server.Lead(),
			})
		}
		return nil
	})
	u.SetName("etcd.get_members")
	u.SetTitle("Etcd members")
	u.SetTags("roles/etcd")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APIMemberJoinInput struct {
	Peer       string   `json:"peer" maxLength:"255"`
	Roles      []string `json:"roles"`
	Identifier string   `json:"identifier"`
}
type APIMemberJoinOutput struct {
	EtcdInitialCluster string `json:"etcdInitialCluster"`
}

func (r *Role) APIClusterJoin() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIMemberJoinInput, output *APIMemberJoinOutput) error {
		if !r.clusterCanJoin(ctx) {
			return status.Unavailable
		}

		r.i.DispatchEvent(backup.EventTopicBackupRun, roles.NewEvent(ctx, map[string]interface{}{}))
		initialCluster := []string{}
		members, err := r.i.KV().MemberList(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, mem := range members.Members {
			for _, u := range mem.PeerURLs {
				initialCluster = append(initialCluster, fmt.Sprintf("%s=%s", mem.Name, u))
			}
		}
		initialCluster = append(initialCluster, fmt.Sprintf(
			"%s=%s", input.Identifier, input.Peer,
		))

		// Pre-configure roles for new node
		if len(input.Roles) == 0 {
			input.Roles = strings.Split(extconfig.Get().BootstrapRoles, ",")
			// If we're copying our roles, exclude backup
			input.Roles = slices.DeleteFunc(input.Roles, func(role string) bool {
				return role == "backup"
			})
		}
		_, err = r.i.KV().Put(
			ctx,
			r.i.KV().Key(
				types.KeyInstance,
				input.Identifier,
				types.KeyRoles,
			).String(),
			strings.Join(input.Roles, ","),
		)
		if err != nil {
			r.log.Warn("failed to put roles for node", zap.Error(err))
		}

		_, err = r.i.KV().MemberAddAsLearner(context.Background(), []string{input.Peer})
		if err != nil {
			r.log.Warn("failed to add member", zap.Error(err))
		}

		output.EtcdInitialCluster = strings.Join(initialCluster, ",")
		return nil
	})
	u.SetName("etcd.join_member")
	u.SetTitle("Etcd join")
	u.SetTags("roles/etcd")
	u.SetExpectedErrors(status.Internal, status.Unavailable)
	return u
}

type APIMemberInput struct {
	PeerID string `query:"peerID" required:"true"`
}

func (r *Role) APIClusterRemove() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIMemberInput, output *struct{}) error {
		r.i.DispatchEvent(backup.EventTopicBackupRun, roles.NewEvent(ctx, map[string]interface{}{}))
		iid, err := strconv.ParseUint(input.PeerID, 16, 64)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		r.i.Log().Debug("Removing instance", zap.Uint64("id", iid))
		go func() {
			_, err := r.i.KV().MemberRemove(context.Background(), iid)
			if err != nil {
				r.log.Warn("failed to remove member", zap.Error(err))
			}
		}()

		return nil
	})
	u.SetName("etcd.remove_member")
	u.SetTitle("Etcd remove")
	u.SetTags("roles/etcd")
	u.SetExpectedErrors(status.Internal)
	return u
}

func (r *Role) APIClusterMoveLeader() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIMemberInput, output *struct{}) error {
		iid, err := strconv.ParseUint(input.PeerID, 16, 64)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		r.i.Log().Debug("Moving leader", zap.Uint64("id", iid))

		_, err = r.i.KV().MoveLeader(ctx, iid)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("etcd.move_leader")
	u.SetTitle("Etcd move leader")
	u.SetTags("roles/etcd")
	u.SetExpectedErrors(status.Internal)
	return u
}

func (r *Role) clusterCanJoin(ctx context.Context) bool {
	st, err := r.lcr.ClusterStatus(ctx)
	if err != nil {
		r.log.Warn("failed to check cluster status", zap.Error(err))
		return false
	}
	if st.Healthy != nil {
		r.log.Warn("cluster is not healthy", zap.Error(err))
		return false
	}
	lds := st.FindLeaderStatus()
	if st := st.FindLearnerStatus(); st != nil {
		r.log.Info("Found learner")
		if IsLearnerReady(lds, st) {
			r.log.Info("Learner is ready, leader should promote it")
		} else {
			r.log.Info("Learner is not ready yet")
		}
		return false
	}
	return true
}
