package clonec

import (
	"fmt"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/result"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/repositories"
)

// Run run command in app.App context
func Run(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Interactive string `argument:"?interactive" ,command:"?interactive"`

			RepositoryURL      string `command:"?$1"`
			DestPath           string `command:"?$2"`
			RepositoryBranch   string `command:"?branch"`
			RepositoryRevision string `command:"?rev"`

			RootFilespace filesystem.Filespace `filespace:"root"`

			RepositoriesService gcliservices.RepositoriesService `dependency:"RepositoriesService"`
			PropertiesService   gcliservices.PropertiesService   `dependency:"PropertiesService"`
			CloneService        gcliservices.ClonerService       `dependency:"ClonerService"`
		}
		repofs         filesystem.Filespace
		propertiesDef  []*config.Property
		propertiesData map[string]string
		isChanged      bool
		destfs         filesystem.Filespace
		interactive    bool
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctx.Scope().InjectTo(&deps); err != nil {
		return err
	}
	interactive = strings.ToLower(deps.Interactive) != "false"
	if deps.RepositoryURL == "" {
		return fmt.Errorf("First argument repository url is required")
	}
	if deps.DestPath == "" {
		return fmt.Errorf("Second argument destination path is required")
	}
	version := repositories.Version{
		Branch:   deps.RepositoryBranch,
		Revision: deps.RepositoryRevision,
	}
	if repofs, err = deps.RepositoriesService.Filespace(deps.RepositoryURL, version); err != nil {
		return err
	}
	if propertiesDef, err = deps.PropertiesService.ReadDefFromFS(repofs); err != nil {
		return err
	}
	if propertiesData, err = deps.PropertiesService.ReadDataFromFS(repofs); err != nil {
		return err
	}
	if isChanged, err = deps.PropertiesService.FillData(propertiesDef, propertiesData, map[string]string{}, interactive); err != nil {
		return err
	}
	if err = deps.RootFilespace.MkdirAll(deps.DestPath, 0766); err != nil {
		return err
	}
	if destfs, err = deps.RootFilespace.Filespace(deps.DestPath); err != nil {
		return err
	}
	if isChanged {
		if err = deps.PropertiesService.WriteDataToFS(destfs, propertiesData); err != nil {
			return err
		}
	}
	propertiesResult := result.NewPropertiesResult(propertiesData)
	if err = deps.CloneService.Clone(deps.RepositoryURL, version, destfs, propertiesResult); err != nil {
		return err
	}
	return ctx.IO().Out().Printf("cloned")
}
