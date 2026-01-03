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
	tests.Setup(t)
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

	var output tftp.APIFilesGetOutput
	assert.NoError(t, role.APIFilesGet().Interact(ctx, struct{}{}, &output))
	assert.NotNil(t, output)
	assert.Len(t, output.Files, 1)
	assert.Equal(t, output.Files[0].Host, "1.2.3.4")
	assert.Equal(t, output.Files[0].Name, "file")
	assert.Equal(t, output.Files[0].SizeBytes, len(data))
}

func TestAPIFilesPut(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("tftp", ctx)
	role := tftp.New(inst)

	data := securecookie.GenerateRandomKey(32)
	assert.NoError(t, role.APIFilesPut().Interact(ctx, tftp.APIFilesPutInput{
		Name: "foo",
		Host: "bar",
		Data: data,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyFiles,
			"bar",
			"foo",
		),
		data,
	)
}

func TestAPIFilesDownload(t *testing.T) {
	tests.Setup(t)
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

	var output tftp.APIFilesDownloadOutput
	assert.NoError(t, role.APIFilesDownload().Interact(ctx, tftp.APIFilesDownloadInput{
		Host: "1.2.3.4",
		Name: "file",
	}, &output))
	assert.NotNil(t, output)
	assert.Len(t, output.Data, len(data))
}

func TestAPIFilesDelete(t *testing.T) {
	tests.Setup(t)
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

	assert.NoError(t, role.APIFilesDelete().Interact(ctx, tftp.APIFilesDeleteInput{
		Host: "1.2.3.4",
		Name: "file",
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyFiles,
			"1.2.3.4",
			"file",
		),
	)
}
