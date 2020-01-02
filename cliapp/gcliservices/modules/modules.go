package modules

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
)

var (
	errInitRequired = fmt.Errorf("ModulesService: Init is required before use dependency")
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
	config []*config.Module
}

// Factory create new repositories instance
func Factory(dp dependency.Provider) (interface{}, error) {
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
