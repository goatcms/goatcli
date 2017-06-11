package modules

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
)

var (
	errInitRequired = fmt.Errorf("ModulesService: Init is required before use dependency")
)

const (
	modulesDefPath = ".goat/modules.def.json"
)

// Modules provide project modules data
type Modules struct {
	deps struct {
		FS filesystem.Filespace `filespace:"root"`
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
	return services.Modules(instance), nil
}

func (m *Modules) Init() (err error) {
	if m.deps.FS.IsFile(modulesDefPath) {
		var defJSON []byte
		defJSON, err = m.deps.FS.ReadFile(modulesDefPath)
		if err != nil {
			return err
		}
		m.config, err = config.NewModules(defJSON)
		if err != nil {
			return err
		}
	} else {
		m.config = make([]*config.Module, 0)
	}
	return nil
}

// ModulesConfig return modules config
func (m *Modules) ModulesConfig() ([]*config.Module, error) {
	if m.config == nil {
		return nil, errInitRequired
	}
	return m.config, nil
}
