package depsc

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
)

// RunAddGODep run deps:add:go command
func RunAddGODep(a app.App) (err error) {
	var (
		deps struct {
			GORepo string `argument:"?$2"`
			Branch string `argument:"?$3"`
			Rev    string `argument:"?$3"`
			CWD    string `argument:"?cwd"`

			Dependencies services.DependenciesService `dependency:"DependenciesService"`
		}
		fs   filesystem.Filespace
		list []*config.Dependency
		dest string
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if deps.GORepo == "" {
		return fmt.Errorf("First argument: golang repository path (like github.com/goatcms/goatcore) is required")
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
	dest = "vendor/" + deps.GORepo
	list = append(list, &config.Dependency{
		Repo:   "git://" + deps.GORepo,
		Branch: deps.Branch,
		Rev:    deps.Rev,
		Dest:   dest,
	})
	if err = deps.Dependencies.WriteDefToFS(fs, list); err != nil {
		return err
	}
	return nil
}
