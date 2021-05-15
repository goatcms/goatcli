package template

import (
	"text/template"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/gtprovider"
	"github.com/goatcms/goatcore/app"
)

// AssetsProvider is global template provider
type AssetsProvider struct {
	deps struct {
		Config gcliservices.TemplateConfig `dependency:"TemplateConfig"`
	}
	plainProvider *gtprovider.AssetsProvider
}

// AssetsProviderFactory create new AssetsProviderFactory instance
func AssetsProviderFactory(dp app.DependencyProvider) (i interface{}, err error) {
	provider := &AssetsProvider{}
	if err = dp.InjectTo(&provider.deps); err != nil {
		return nil, err
	}
	return gcliservices.TemplateAssetsProvider(provider), nil
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
