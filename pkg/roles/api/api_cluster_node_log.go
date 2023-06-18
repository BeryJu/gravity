package api

import (
	"context"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/extconfig/log_iml"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type APILogMessage struct {
	Message string `json:"message"`
	Node    string `json:"node"`
}

type APILogMessages struct {
	IsJSON   bool            `json:"isJSON"`
	Messages []APILogMessage `json:"messages"`
}

func (r *Role) APIClusterNodeLogMessages() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APILogMessages) error {
		output.IsJSON = !extconfig.Get().Debug
		for _, lm := range log_iml.Get().Messages() {
			output.Messages = append(output.Messages, APILogMessage{
				Message: lm,
				Node:    extconfig.Get().Instance.Identifier,
			})
		}
		return nil
	})
	u.SetName("api.get_log_messages")
	u.SetTitle("Log messages")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}
