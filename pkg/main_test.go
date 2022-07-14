package instance

import (
	"os"
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/api/auth"
)

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "gravity-etcd-test-*")
	if err != nil {
		panic(err)
	}
	extconfig.Get().DataPath = tmpDir
	rootInst := instance.NewInstance()
	inst := rootInst.ForRole("tests")
	inst.AddEventListener(instance.EventTopicInstanceBootstrapped, func(ev *roles.Event) {
		api := api.New(inst)
		am := auth.NewAuthProvider(api, inst)
		err = am.CreateUser("foo", string("bar"))
		if err != nil {
			panic(err)
		}
		m.Run()
		defer func() {
			rootInst.Stop()
			os.RemoveAll(tmpDir)
		}()
	})
	rootInst.Start()
}
