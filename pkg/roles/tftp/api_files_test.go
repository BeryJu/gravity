package tftp_test

import (
	"encoding/base64"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/tftp"
	"beryju.io/gravity/pkg/roles/tftp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
)

func TestAPIFilesGet(t *testing.T) {
	defer tests.Setup(t)()
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("tftp", ctx)
	role := tftp.New(inst)

	data := base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyFiles,
			"1.2.3.4",
			"file",
		).String(),
		data,
	))

	var output tftp.APIFilesOutput
	assert.NoError(t, role.APIFilesGet().Interact(ctx, struct{}{}, &output))
	assert.NotNil(t, output)
	assert.Len(t, output.Files, 1)
	assert.Equal(t, output.Files[0].Host, "1.2.3.4")
	assert.Equal(t, output.Files[0].Name, "file")
	assert.Equal(t, output.Files[0].SizeBytes, len(data))
}
