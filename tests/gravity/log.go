package gravity

import (
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

// StdoutLogConsumer is a LogConsumer that prints the log to stdout
type StdoutLogConsumer struct {
	T      *testing.T
	Prefix string
}

// Accept prints the log to stdout
func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	lc.T.Logf("%s: %s", lc.Prefix, string(l.Content))
}
