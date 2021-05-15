package gclicore

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
)

// Core provider
type Core struct {
	deps struct {
		Arguments gcliservices.Arguments `dependency:"GCLICoreArguments"`
	}
}

// CoreFactory create new Core instance
func CoreFactory(dp app.DependencyProvider) (interface{}, error) {
	var err error
	instance := &Core{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return gcliservices.Core(instance), nil
}

func (core *Core) Arguments() gcliservices.Arguments {
	return core.deps.Arguments
}
