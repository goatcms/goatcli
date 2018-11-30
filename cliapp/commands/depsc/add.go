package depsc

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
)

// RunAddDep run deps:add command
func RunAddDep(a app.App) (err error) {
	var (
		deps struct {
			Path   string `argument:"?$2"`
			Repo   string `argument:"?$3"`
			Branch string `argument:"?$4"`
			Rev    string `argument:"?$5"`
			CWD    string `argument:"?cwd"`

			Dependencies services.DependenciesService `dependency:"DependenciesService"`
		}
		fs   filesystem.Filespace
		list []*config.Dependency
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if deps.Path == "" {
		return fmt.Errorf("First argument: destination path (like vendor/github.com/goatcms/goatcore) is required")
	}
	if deps.Repo == "" {
		return fmt.Errorf("Second argument: repository url (like http://github.com/goatcms/goatcore.git) is required")
	}
	if deps.CWD == "" {
		deps.CWD = "./"
	}
	if fs, err = diskfs.NewFilespace(deps.CWD); err != nil {
		return err
	}
	if list, err = deps.Dependencies.ReadDefFromFS(fs); err != nil {
		return err
	}
	list = append(list, &config.Dependency{
		Repo:   deps.Repo,
		Branch: deps.Branch,
		Rev:    deps.Rev,
		Dest:   deps.Path,
	})
	if err = deps.Dependencies.WriteDefToFS(fs, list); err != nil {
		return err
	}
	return nil
}
