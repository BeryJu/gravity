package instance

import (
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/dop251/goja"
	"go.uber.org/zap"
)

func (ri *RoleInstance) HookMeth(src string, meth string, args ...interface{}) {
	if src == "" {
		return
	}
	log := ri.log.With(zap.String("hook", meth))
	vm := goja.New()
	vm.Set("gravity", map[string]interface{}{
		"log": func(msg goja.Value) {
			log.Info(msg.String())
		},
		"node":    extconfig.Get().Instance.Identifier,
		"version": extconfig.Version,
	})
	_, err := vm.RunString(src)
	if err != nil {
		log.Warn("failed to run scope hook", zap.Error(err))
		return
	}
	hookMeth, ok := goja.AssertFunction(vm.Get(meth))
	if !ok {
		log.Warn("hook not a function", zap.String("meth", meth))
		return
	}
	convertedArgs := []goja.Value{}
	for _, rv := range args {
		convertedArgs = append(convertedArgs, vm.ToValue(rv))
	}
	before := time.Now()
	_, err = hookMeth(vm.ToValue(ri), convertedArgs...)
	duration := time.Since(before)
	instanceRoleHooks.WithLabelValues(ri.roleId, meth).Observe(duration.Seconds())
	if err != nil {
		log.Warn("failed to call hook function", zap.Error(err))
	}
}
