package template

import (
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp dependency.Provider) (err error) {
	return goaterr.ToErrors(goaterr.AppendError(nil,
		dp.AddDefaultFactory("TemplateAssetsProvider", AssetsProviderFactory),
		dp.AddDefaultFactory("TemplateConfig", ConfigFactory),
		dp.AddDefaultFactory("TemplateService", ProviderFactory),
	))
}
