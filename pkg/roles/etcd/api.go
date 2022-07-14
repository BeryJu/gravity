package etcd

import (
	"context"
	"fmt"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *EmbeddedEtcd) apiHandlerMembers() usecase.Interactor {
	type member struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}
	type membersOutput struct {
		Members []member `json:"members"`
	}
	u := usecase.NewIOI(new(struct{}), new(membersOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*membersOutput)
		)
		members, err := r.i.KV().MemberList(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, mem := range members.Members {
			out.Members = append(out.Members, member{
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

func (r *EmbeddedEtcd) apiHandlerJoin() usecase.Interactor {
	type etcdJoinInput struct {
		Peer string `query:"peer"`
	}
	type etcdJoinOutput struct {
		Env string `json:"env"`
	}
	u := usecase.NewIOI(new(etcdJoinInput), new(etcdJoinOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*etcdJoinInput)
			out = output.(*etcdJoinOutput)
		)
		_, err := r.i.KV().MemberAdd(ctx, []string{in.Peer})
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		env := fmt.Sprintf(
			"ETCD_JOIN_CLUSTER='%s;%s'",
			extconfig.Get().Instance.Identifier,
			extconfig.Get().Instance.IP,
		)
		out.Env = env
		return nil
	})
	u.SetName("etcd.join_member")
	u.SetTitle("Etcd join")
	u.SetTags("roles/etcd")
	u.SetExpectedErrors(status.Internal)
	return u
}
