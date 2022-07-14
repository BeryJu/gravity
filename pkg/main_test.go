package instance

import (
	"os"
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
)

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "gravity-etcd-test-*")
	if err != nil {
		panic(err)
	}
	extconfig.Get().DataPath = tmpDir
	inst := instance.NewInstance()
	inst.ForRole("tests").AddEventListener(instance.EventTopicInstanceBootstrapped, func(ev *roles.Event) {
		api := api.New(inst.ForRole("api"))
		err = api.CreateUser("foo", string("bar"))
		if err != nil {
			panic(err)
		}
		m.Run()
		defer func() {
			inst.Stop()
			os.RemoveAll(tmpDir)
		}()
	})
	inst.Start()
}
