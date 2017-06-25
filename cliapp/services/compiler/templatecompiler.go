package compiler

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
)

// TemplateCompiler clone and process new project
type TemplateCompiler struct {
	deps struct {
		Template services.Template `dependency:"TemplateService"`
	}
}

// TemplateCompilerFactory create new FileCompiler instance
func TemplateCompilerFactory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &TemplateCompiler{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return services.TemplateCompiler(instance), nil
}
