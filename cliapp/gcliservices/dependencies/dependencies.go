package dependencies

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/fshelper"
	"github.com/goatcms/goatcore/repositories"
	"github.com/goatcms/goatcore/varutil"
)

// Dependencies provide project dependencies data
type Dependencies struct {
	deps struct {
		FS           filesystem.Filespace             `filespace:"current"`
		Repositories gcliservices.RepositoriesService `dependency:"RepositoriesService"`
	}
}

// Factory create new repositories instance
func Factory(dp app.DependencyProvider) (interface{}, error) {
	var err error
	instance := &Dependencies{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return gcliservices.DependenciesService(instance), nil
}

// ReadDefFromFS read dependencies definitions from filesystem
func (m *Dependencies) ReadDefFromFS(fs filesystem.Filespace) (dependencies []*config.Dependency, err error) {
	var json []byte
	if !fs.IsFile(DependenciesDefPath) {
		return make([]*config.Dependency, 0), nil
	}
	if json, err = fs.ReadFile(DependenciesDefPath); err != nil {
		return nil, err
	}
	if dependencies, err = config.NewDependencies(json); err != nil {
		return nil, err
	}
	return dependencies, nil
}

// WriteDefToFS write dependencies definitions to filespace
func (m *Dependencies) WriteDefToFS(fs filesystem.Filespace, deps []*config.Dependency) (err error) {
	var (
		json string
	)
	if json, err = varutil.ObjectToJSON(deps); err != nil {
		return err
	}
	return fs.WriteFile(DependenciesDefPath, []byte(json), filesystem.DefaultUnixFileMode)
}

// CloneDependencies download project dependencies
func (m *Dependencies) CloneDependencies(fs filesystem.Filespace, deps []*config.Dependency) (err error) {
	var (
		repoFS filesystem.Filespace
		dFS    filesystem.Filespace
	)
	for _, dep := range deps {
		if fs.IsExist(dep.Dest) {
			continue
		}
		if repoFS, err = m.deps.Repositories.Filespace(dep.Repo, repositories.Version{
			Branch:   dep.Branch,
			Revision: dep.Rev,
		}); err != nil {
			return err
		}
		if err = fs.MkdirAll(dep.Dest, filesystem.DefaultUnixDirMode); err != nil {
			return err
		}
		if dFS, err = fs.Filespace(dep.Dest); err != nil {
			return err
		}
		if err = fshelper.Copy(repoFS, dFS, func(fs filesystem.Filespace, subPath string) bool {
			return subPath != ".git"
		}); err != nil {
			return nil
		}
	}
	return nil
}
