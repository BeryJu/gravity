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
	Message string                 `json:"message" required:"true"`
	Time    time.Time              `json:"time" required:"true"`
	Level   string                 `json:"level" required:"true"`
	Logger  string                 `json:"logger" required:"true"`
	Fields  map[string]interface{} `json:"fields" required:"true"`
	Node    string                 `json:"node" required:"true"`
}

type APILogMessages struct {
	Messages []APILogMessage `json:"messages"`
}

func (r *Role) APIClusterNodeLogMessages() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APILogMessages) error {
		for _, lm := range log_iml.Get().Messages() {
			output.Messages = append(output.Messages, APILogMessage{
				Message: lm.Entry.Message,
				Level:   lm.Entry.Level.CapitalString(),
				Time:    lm.Entry.Time,
				Logger:  lm.Entry.LoggerName,
				Fields:  lm.FieldsToMap(),
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
