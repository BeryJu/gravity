package api

import (
	"context"
	"fmt"
	"strings"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/backup"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"go.uber.org/zap"
)

type APIMember struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
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
				ID:   mem.ID,
				Name: mem.Name,
			})
		}
		return nil
	})
	u.SetName("api.get_members")
	u.SetTitle("Etcd members")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APIMemberJoinInput struct {
	Peer string `json:"peer" maxLength:"255"`
}
type APIMemberJoinOutput struct {
	Env string `json:"env"`
}

func (r *Role) APIClusterJoin() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIMemberJoinInput, output *APIMemberJoinOutput) error {
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

		go func() {
			_, err = r.i.KV().MemberAdd(ctx, []string{input.Peer})
			if err != nil {
				r.log.Warn("failed to add member", zap.Error(err))
			}
		}()

		output.Env = strings.Join(initialCluster, ",")
		return nil
	})
	u.SetName("etcd.join_member")
	u.SetTitle("Etcd join")
	u.SetTags("roles/etcd")
	u.SetExpectedErrors(status.Internal)
	return u
}
