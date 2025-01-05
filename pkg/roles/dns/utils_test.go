package dns_test

import (
	"errors"
	"net"
	"testing"

	"beryju.io/gravity/pkg/roles/dns"
	d "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

type NullDNSWriter struct {
	closed bool
	msg    *d.Msg
}

func (nw *NullDNSWriter) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 53}
}

func (nw *NullDNSWriter) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1053}
}
func (nw *NullDNSWriter) WriteMsg(msg *d.Msg) error { nw.msg = msg; return nil }
func (nw *NullDNSWriter) Write([]byte) (int, error) { return 0, errors.New("Write not supported") }
func (nw *NullDNSWriter) Close() error              { nw.closed = true; return nil }
func (nw *NullDNSWriter) TsigStatus() error         { return nil }
func (nw *NullDNSWriter) TsigTimersOnly(v bool)     {}
func (nw *NullDNSWriter) Hijack()                   {}
func (nw *NullDNSWriter) Msg() *d.Msg               { return nw.msg }

func NewNullDNSWriter() *NullDNSWriter {
	return &NullDNSWriter{}
}

func AssertDNS(t *testing.T, role *dns.Role, q []d.Question, expected ...string) {
	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: q,
	})
	// Convert given answers and expected answers to strings
	givenAnswersStr := []string{}
	for _, a := range fw.Msg().Answer {
		givenAnswersStr = append(givenAnswersStr, a.String())
	}
	expectedAnswersStr := []string{}
	for _, a := range expected {
		r, err := d.NewRR(a)
		assert.NoError(t, err)
		assert.NotNil(t, r)
		expectedAnswersStr = append(expectedAnswersStr, r.String())
	}
	assert.Len(t, givenAnswersStr, len(expectedAnswersStr), "Count of answers is mismatched")
	assert.ElementsMatch(t, givenAnswersStr, expectedAnswersStr, "Given and expected answers mismatch")
}
