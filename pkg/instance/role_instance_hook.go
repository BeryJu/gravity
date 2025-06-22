package instance

import (
	"net"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"github.com/dop251/goja"
	"go.uber.org/zap"
)

func (ri *RoleInstance) HookEnvironment(options roles.HookOptions) map[string]interface{} {
	return map[string]interface{}{
		"gravity": map[string]interface{}{
			"log": func(msg goja.Value) {
				ri.log.With(zap.String("hook", options.Method)).Info(msg.String())
			},
			"node":    extconfig.Get().Instance.Identifier,
			"version": extconfig.Version,
			"role":    ri,
		},
		"net": map[string]interface{}{
			"parseIP": func(input goja.Value, family goja.Value) []byte {
				parsed := net.ParseIP(input.String())
				if family.String() == "v4" {
					return parsed[12:]
				}
				return parsed
			},
		},
		"strconv": map[string]interface{}{
			"toBytes": func(input goja.Value) []byte {
				return []byte(input.String())
			},
		},
	}
}

func (ri *RoleInstance) ExecuteHook(options roles.HookOptions, args ...interface{}) interface{} {
	if options.Source == "" {
		return nil
	}
	log := ri.log.With(zap.String("hook", options.Method))

	vm := goja.New()
	for k, v := range options.Env {
		err := vm.Set(k, v)
		if err != nil {
			log.Warn("failed to set environment", zap.Error(err))
		}
	}
	for k, v := range ri.HookEnvironment(options) {
		err := vm.Set(k, v)
		if err != nil {
			log.Warn("failed to set environment", zap.Error(err))
		}
	}

	_, err := vm.RunString(options.Source)
	if err != nil {
		log.Warn("failed to run scope hook", zap.Error(err))
		return nil
	}
	m := vm.Get(options.Method)
	hookMeth, ok := goja.AssertFunction(m)
	if !ok {
		return nil
	}
	convertedArgs := []goja.Value{}
	for _, rv := range args {
		convertedArgs = append(convertedArgs, vm.ToValue(rv))
	}
	before := time.Now()
	v, err := hookMeth(vm.ToValue(ri), convertedArgs...)
	duration := time.Since(before)
	instanceRoleHooks.WithLabelValues(ri.roleId, options.Method).Observe(duration.Seconds())
	if err != nil {
		log.Warn("failed to call hook function", zap.Error(err))
		return nil
	}
	return v.Export()
}
