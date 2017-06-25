package compiler

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
)

// Compiler compile file
type Compiler struct {
	deps struct{}
}

// CompilerFactory create new repositories instance
func CompilerFactory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &Compiler{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return services.Compiler(instance), nil
}
