package api

import (
	"context"

	"beryju.io/gravity/pkg/extconfig/log_iml"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type APILogMessages struct {
	Messages []string `json:"messages"`
}

func (r *Role) APIClusterNodeLogMessages() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APILogMessages) error {
		output.Messages = log_iml.Get().Messages()
		return nil
	})
	u.SetName("api.get_log_messages")
	u.SetTitle("Log messages")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}
