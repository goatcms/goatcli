package cleanc

import (
	"github.com/goatcms/goatcli/cliapp/commands/vcsc"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunCleanBuild run clean:build command
func RunCleanBuild(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			CurrentFS  filesystem.Filespace `filespace:"current"`
			VCSService services.VCSService  `dependency:"VCSService"`
		}
		vcsData services.VCSData
	)
	if err = vcsc.RunScan(a, ctx); err != nil {
		return nil
	}
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctx.Scope().InjectTo(&deps))); err != nil {
		return err
	}
	// load data
	if vcsData, err = deps.VCSService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	ignored := vcsData.VCSIgnoredFiles()
	for _, row := range vcsData.VCSGeneratedFiles().All() {
		if !deps.CurrentFS.IsFile(row.Path) || ignored.ContainsPath(row.Path) {
			continue
		}
		if err = deps.CurrentFS.Remove(row.Path); err != nil {
			return err
		}
	}
	return ctx.IO().Out().Printf("cleaned\n")
}
