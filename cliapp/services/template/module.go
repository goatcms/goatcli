package template

import (
	"github.com/goatcms/goatcli/cliapp/services/template/tfunc"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/dependency"
)

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp dependency.Provider) error {
	return dp.AddDefaultFactory("TemplateService", ProviderFactory)
}

// InitDependencies is init callback to inject dependencies inside module
func InitDependencies(a app.App) (err error) {
	return tfunc.RegisterFunctions(a.DependencyProvider())
}
