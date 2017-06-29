package cloner

import (
	"os"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/fsloop"
)

// Cloner clone and process new project
type Cloner struct {
	deps struct {
		Repositories services.RepositoriesService `dependency:"RepositoriesService"`
		Modules      services.ModulesService      `dependency:"ModulesService"`
	}
}

// Factory create new repositories instance
func Factory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &Cloner{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return services.ClonerService(instance), nil
}

// Clone clone repository
func (cloner *Cloner) Clone(repository, rev string, destfs filesystem.Filespace, si common.StringInjector) (err error) {
	var (
		sourcefs filesystem.Filespace
		replaces []*config.Replace
		modules  []*config.Module
	)
	if sourcefs, err = cloner.deps.Repositories.Filespace(repository, rev); err != nil {
		return err
	}
	if modules, err = cloner.deps.Modules.ReadDefFromFS(sourcefs); err != nil {
		return err
	}
	for _, module := range modules {
		var modulefs filesystem.Filespace
		if err = destfs.MkdirAll(module.SourceDir, 0766); err != nil {
			return err
		}
		if modulefs, err = destfs.Filespace(module.SourceDir); err != nil {
			return err
		}
		if err = cloner.Clone(module.SourceURL, module.SourceRev, modulefs, si); err != nil {
			return err
		}
	}
	if sourcefs.IsFile(ReplaceConfigFile) {
		var json []byte
		if json, err = sourcefs.ReadFile(ReplaceConfigFile); err != nil {
			return err
		}
		if replaces, err = config.NewReplaces(json, si); err != nil {
			return err
		}
	}
	cleanRequired := false
	loop := fsloop.NewLoop(&fsloop.LoopData{
		Filespace: sourcefs,
		DirFilter: func(fs filesystem.Filespace, subPath string) bool {
			return subPath != "./.git"
		},
		OnDir: func(fs filesystem.Filespace, subPath string) error {
			return destfs.MkdirAll(subPath, 0777)
		},
		OnFile: func(fs filesystem.Filespace, subPath string) error {
			if err = copy(sourcefs, destfs, subPath, replaces); err != nil {
				cleanRequired = true
				return err
			}
			return nil
		},
		Consumers:  1,
		Producents: 1,
	}, nil)
	loop.Run("")
	loop.Wait()
	if cleanRequired {
		var dirnodes []os.FileInfo
		dirnodes, err = destfs.ReadDir("./")
		for _, dirnode := range dirnodes {
			if err = destfs.RemoveAll(dirnode.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}
