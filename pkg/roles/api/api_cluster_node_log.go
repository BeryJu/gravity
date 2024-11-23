package api

import (
	"context"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/extconfig/log_iml"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type APILogMessage struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Logger  string    `json:"logger"`

	Node string `json:"node"`
}

type APILogMessages struct {
	Messages []APILogMessage `json:"messages"`
}

func (r *Role) APIClusterNodeLogMessages() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APILogMessages) error {
		for _, lm := range log_iml.Get().Messages() {
			output.Messages = append(output.Messages, APILogMessage{
				Message: lm.Message,
				Level:   lm.Level.CapitalString(),
				Time:    lm.Time,
				Logger:  lm.LoggerName,
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
