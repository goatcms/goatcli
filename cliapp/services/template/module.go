package template

import (
	"github.com/goatcms/goatcli/cliapp/services/template/tfunc"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/dependency"
)

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp dependency.Provider) error {
	if err := dp.AddDefaultFactory("TemplateService", ProviderFactory); err != nil {
		return err
	}
	return nil
}

// InitDependencies is init callback to inject dependencies inside module
func InitDependencies(a app.App) (err error) {
	return tfunc.RegisterFunctions(a.DependencyProvider())
}
