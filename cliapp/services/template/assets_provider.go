package template

import (
	"text/template"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/template/gtprovider"
	"github.com/goatcms/goatcore/dependency"
)

// AssetsProvider is global template provider
type AssetsProvider struct {
	deps struct {
		Config services.TemplateConfig `dependency:"TemplateConfig"`
	}
	plainProvider *gtprovider.AssetsProvider
}

// AssetsProviderFactory create new AssetsProviderFactory instance
func AssetsProviderFactory(dp dependency.Provider) (i interface{}, err error) {
	provider := &AssetsProvider{}
	if err = dp.InjectTo(&provider.deps); err != nil {
		return nil, err
	}
	return services.TemplateAssetsProvider(provider), nil
}

// Build create new template based on layout
func (provider *AssetsProvider) init() {
	if provider.plainProvider != nil {
		return
	}
	provider.plainProvider = gtprovider.NewAssetsProvider(provider.deps.Config.FS(), HelpersPath, LayoutPath, provider.deps.Config.Func(), provider.deps.Config.IsCached())
}

// Layout return leyout template by name
func (provider *AssetsProvider) Layout(name string) (tmpl *template.Template, err error) {
	provider.init()
	return provider.plainProvider.Layout(name)
}

// Base return base template
func (provider *AssetsProvider) Base() (*template.Template, error) {
	provider.init()
	return provider.plainProvider.Base()
}
