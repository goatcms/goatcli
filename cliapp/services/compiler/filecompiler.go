package compiler

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
)

// FileCompiler clone and process new project
type FileCompiler struct {
	deps struct{}
}

// FileCompilerFactory create new FileCompiler instance
func FileCompilerFactory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &FileCompiler{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return services.FileCompiler(instance), nil
}
