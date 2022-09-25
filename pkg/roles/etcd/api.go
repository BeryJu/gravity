package etcd

import (
	"context"
	"fmt"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/backup"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *Role) apiHandlerMembers() usecase.Interactor {
	type member struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}
	type membersOutput struct {
		Members []member `json:"members"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *membersOutput) error {
		members, err := r.i.KV().MemberList(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, mem := range members.Members {
			output.Members = append(output.Members, member{
				ID:   mem.ID,
				Name: mem.Name,
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

func (r *Role) apiHandlerJoin() usecase.Interactor {
	type etcdJoinInput struct {
		Peer string `json:"peer" maxLength:"255"`
	}
	type etcdJoinOutput struct {
		Env string `json:"env"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input etcdJoinInput, output *etcdJoinOutput) error {
		r.i.DispatchEvent(backup.EventTopicBackupRun, roles.NewEvent(ctx, map[string]interface{}{}))

		_, err := r.i.KV().MemberAdd(ctx, []string{input.Peer})
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		env := fmt.Sprintf(
			"ETCD_JOIN_CLUSTER='%s;%s'",
			extconfig.Get().Instance.Identifier,
			extconfig.Get().Instance.IP,
		)
		output.Env = env
		return nil
	})
	u.SetName("etcd.join_member")
	u.SetTitle("Etcd join")
	u.SetTags("roles/etcd")
	u.SetExpectedErrors(status.Internal)
	return u
}
