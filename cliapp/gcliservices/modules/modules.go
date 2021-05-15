package modules

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

const (
	// ModulesDefPath is path to modules definition
	ModulesDefPath = ".goat/modules.def.json"
)

// Modules provide project modules data
type Modules struct {
	deps struct {
		FS filesystem.Filespace `filespace:"current"`
	}
}

// Factory create new repositories instance
func Factory(dp app.DependencyProvider) (interface{}, error) {
	var err error
	instance := &Modules{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return gcliservices.ModulesService(instance), nil
}

// ReadDefFromFS read modules definitions from filesystem
func (m *Modules) ReadDefFromFS(fs filesystem.Filespace) (modules []*config.Module, err error) {
	var json []byte
	if !fs.IsFile(ModulesDefPath) {
		return make([]*config.Module, 0), nil
	}
	if json, err = fs.ReadFile(ModulesDefPath); err != nil {
		return nil, err
	}
	if modules, err = config.NewModules(json); err != nil {
		return nil, err
	}
	return modules, nil
}
