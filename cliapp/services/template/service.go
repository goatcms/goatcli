package template

import (
	"fmt"
	"sync"
	"text/template"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/goattext/gtprovider"
)

// Service is global template provider
type Service struct {
	deps struct {
		TemplateCache string               `argument:"?template.cache"`
		Filespace     filesystem.Filespace `filespace:"root"`
	}
	providerMutex sync.Mutex
	funcs         template.FuncMap
	isUsed        bool
	cache         bool
}

// ProviderFactory create new template provider
func ProviderFactory(dp dependency.Provider) (interface{}, error) {
	s := &Service{
		funcs:  template.FuncMap{},
		isUsed: false,
	}
	if err := dp.InjectTo(&s.deps); err != nil {
		return nil, err
	}
	s.cache = s.deps.TemplateCache != "false"
	return services.TemplateService(s), nil
}

// AddFunc adds the elements of the argument map to the template's function map.
func (s *Service) AddFunc(name string, f interface{}) error {
	if s.isUsed {
		return fmt.Errorf("template.Service.AddFunc: Add functions to template after init template provider")
	}
	if _, ok := s.funcs[name]; ok {
		return fmt.Errorf("template.Service.AddFunc: Can not add function for the same name %s twice", name)
	}
	s.funcs[name] = f
	return nil
}

// Build create new template based on layout
func (s *Service) Build(fs filesystem.Filespace) (services.TemplateExecutor, error) {
	s.isUsed = true
	// prepare executor
	provider := gtprovider.NewProvider(fs, HelpersPath, LayoutPath, ViewPath, FileExtension, s.funcs, s.cache)
	return &Executor{
		provider: provider,
	}, nil
}
