package gravity

import (
	"fmt"

	"github.com/testcontainers/testcontainers-go"
)

// StdoutLogConsumer is a LogConsumer that prints the log to stdout
type StdoutLogConsumer struct {
	Prefix string
}

// Accept prints the log to stdout
func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	fmt.Printf("%s: %s", lc.Prefix, string(l.Content))
}
