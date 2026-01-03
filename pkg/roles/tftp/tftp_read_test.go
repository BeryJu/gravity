package tftp_test

import (
	"bytes"
	"net"
	"net/netip"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/tftp"
	"beryju.io/gravity/pkg/roles/tftp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

type OutgoingTransfer struct {
	*bytes.Buffer
}

func (ot OutgoingTransfer) RemoteAddr() net.UDPAddr {
	return *net.UDPAddrFromAddrPort(netip.MustParseAddrPort("1.2.3.4:69"))
}
func (ot OutgoingTransfer) SetSize(n int64) {}

func TestTFTP_Read(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("tftp", ctx)
	role := tftp.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte{}))
	defer role.Stop()

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyFiles,
			"1.2.3.4",
			"foo",
		).String(),
		"foobar",
	))

	buff := bytes.NewBuffer([]byte{})
	err := role.Reader("foo", OutgoingTransfer{buff})
	assert.NoError(t, err)
	assert.Len(t, buff.Bytes(), 6)
}

func TestTFTP_Read_Bundled(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("tftp", ctx)
	role := tftp.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte{}))
	defer role.Stop()

	buff := bytes.NewBuffer([]byte{})
	err := role.Reader("bundled/ipxe.undionly.kpxe", OutgoingTransfer{buff})
	assert.NoError(t, err)
	assert.True(t, len(buff.Bytes()) > 0)
}

func TestTFTP_Read_NonExistent(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("tftp", ctx)
	role := tftp.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte{}))
	defer role.Stop()

	buff := bytes.NewBuffer([]byte{})
	err := role.Reader("bundled/foo", OutgoingTransfer{buff})
	assert.Error(t, err)
}
