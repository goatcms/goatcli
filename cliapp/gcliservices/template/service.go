package template

import (
	"text/template"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/gtprovider"
	"github.com/goatcms/goatcore/dependency"
)

// Service is global template provider
type Service struct {
	deps struct {
		AssetsProvider gcliservices.TemplateAssetsProvider `dependency:"TemplateAssetsProvider"`
		Config         gcliservices.TemplateConfig         `dependency:"TemplateConfig"`
	}
}

// ProviderFactory create new template provider
func ProviderFactory(dp dependency.Provider) (interface{}, error) {
	s := &Service{}
	if err := dp.InjectTo(&s.deps); err != nil {
		return nil, err
	}
	return gcliservices.TemplateService(s), nil
}

// TemplatesExecutor return view tree executor
func (s *Service) TemplatesExecutor() (gcliservices.TemplatesExecutor, error) {
	provider := gtprovider.NewTemplatesProvider(s.deps.AssetsProvider, s.deps.Config.FS(), ViewPath, s.deps.Config.IsCached())
	return NewTemplatesExecutor(provider), nil
}

// TemplateExecutor return single template executor
func (s *Service) TemplateExecutor(path string) (exeutor gcliservices.TemplateExecutor, err error) {
	var tmpl *template.Template
	if tmpl, err = s.deps.AssetsProvider.Base(); err != nil {
		return nil, err
	}
	provider := gtprovider.NewTemplateProvider(tmpl, s.deps.Config.FS(), path, s.deps.Config.IsCached())
	return NewTemplateExecutor(provider), nil
}
