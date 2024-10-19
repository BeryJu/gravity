package instance

import (
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"github.com/dop251/goja"
	"go.uber.org/zap"
)

func (ri *RoleInstance) ExecuteHook(options roles.HookOptions, args ...interface{}) {
	if options.Source == "" {
		return
	}
	log := ri.log.With(zap.String("hook", options.Method))

	vm := goja.New()
	for k, v := range options.Env {
		err := vm.Set(k, v)
		if err != nil {
			log.Warn("failed to set environment", zap.Error(err))
		}
	}
	err := vm.Set("gravity", map[string]interface{}{
		"log": func(msg goja.Value) {
			log.Info(msg.String())
		},
		"node":    extconfig.Get().Instance.Identifier,
		"version": extconfig.Version,
		"role":    ri,
	})
	if err != nil {
		log.Warn("failed to set environment", zap.Error(err))
	}

	_, err = vm.RunString(options.Source)
	if err != nil {
		log.Warn("failed to run scope hook", zap.Error(err))
		return
	}
	hookMeth, ok := goja.AssertFunction(vm.Get(options.Method))
	if !ok {
		log.Warn("hook not a function", zap.String("meth", options.Method))
		return
	}
	convertedArgs := []goja.Value{}
	for _, rv := range args {
		convertedArgs = append(convertedArgs, vm.ToValue(rv))
	}
	before := time.Now()
	_, err = hookMeth(vm.ToValue(ri), convertedArgs...)
	duration := time.Since(before)
	instanceRoleHooks.WithLabelValues(ri.roleId, options.Method).Observe(duration.Seconds())
	if err != nil {
		log.Warn("failed to call hook function", zap.Error(err))
	}
}
