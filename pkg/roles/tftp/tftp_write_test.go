package tftp_test

import (
	"bytes"
	"fmt"
	"net"
	"net/netip"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/tftp"
	"beryju.io/gravity/pkg/roles/tftp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

type IncomingTransfer struct {
	*bytes.Buffer
	IP string
}

func (it IncomingTransfer) RemoteAddr() net.UDPAddr {
	return *net.UDPAddrFromAddrPort(netip.MustParseAddrPort(fmt.Sprintf("%s:69", it.IP)))
}
func (it IncomingTransfer) Size() (n int64, ok bool) {
	return int64(it.Len()), true
}

func TestTFTP_Write(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("tftp", ctx)
	role := tftp.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte{}))
	defer role.Stop()

	buff := bytes.NewBufferString("foobar")
	err := role.Writer("filename", IncomingTransfer{buff, "1.2.3.4"})
	assert.NoError(t, err)

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyFiles,
			"1.2.3.4",
			"filename",
		),
		[]byte("foobar"),
	)
}
