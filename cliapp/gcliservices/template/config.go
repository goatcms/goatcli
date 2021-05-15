package template

import (
	"fmt"
	"text/template"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// Config is template provider
type Config struct {
	deps struct {
		Cache string               `argument:"?template.cache"`
		FS    filesystem.Filespace `filespace:"current"`
	}
	funcs  template.FuncMap
	cached bool
	locked bool
}

// ConfigFactory create new Config instance
func ConfigFactory(dp app.DependencyProvider) (i interface{}, err error) {
	config := &Config{
		funcs:  template.FuncMap{},
		locked: false,
	}
	if err = dp.InjectTo(&config.deps); err != nil {
		return nil, err
	}
	config.cached = config.deps.Cache != "false"
	return gcliservices.TemplateConfig(config), nil
}

// AddFunc adds the elements of the argument map to the template's function map.
func (config *Config) AddFunc(name string, f interface{}) error {
	if config.locked {
		return fmt.Errorf("template.Config.AddFunc: Add functions to template config is denied after read configuration")
	}
	if _, ok := config.funcs[name]; ok {
		return fmt.Errorf("template.Config.AddFunc: Can not add function for the same name %s twice", name)
	}
	config.funcs[name] = f
	return nil
}

// Func return assets functions
func (config *Config) Func() template.FuncMap {
	config.locked = true
	return config.funcs
}

// IsCached return true if cache flag is set to true
func (config *Config) IsCached() bool {
	return config.cached
}

// FS return template filesystem
func (config *Config) FS() filesystem.Filespace {
	return config.deps.FS
}
