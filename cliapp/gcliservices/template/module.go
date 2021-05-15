package template

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp app.DependencyProvider) (err error) {
	return goaterr.ToError(goaterr.AppendError(nil,
		dp.AddDefaultFactory("TemplateAssetsProvider", AssetsProviderFactory),
		dp.AddDefaultFactory("TemplateConfig", ConfigFactory),
		dp.AddDefaultFactory("TemplateService", ProviderFactory),
	))
}
