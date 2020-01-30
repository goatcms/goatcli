package vcsc

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunClean run vcs:clean command. It remove path from ignored files foreach doeas't exist file
func RunClean(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			CurrentFS  filesystem.Filespace    `filespace:"current"`
			VCSService gcliservices.VCSService `dependency:"VCSService"`
		}
		vcsData gcliservices.VCSData
		ignored = vcs.NewIgnoredFiles(true)
	)
	if err = goaterr.ToError(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctx.Scope().InjectTo(&deps))); err != nil {
		return err
	}
	if vcsData, err = deps.VCSService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	for _, path := range vcsData.VCSIgnoredFiles().All() {
		if deps.CurrentFS.IsFile(path) {
			ignored.AddPath(path)
		} else {
			ctx.IO().Out().Printf(" Deleted from ignored: %s\n", path)
		}
	}
	if len(ignored.All()) == len(vcsData.VCSIgnoredFiles().All()) {
		ctx.IO().Out().Printf("ignored files are clean\n")
		return nil
	}
	vcsData = vcs.NewData(vcsData.VCSGeneratedFiles(), ignored)
	if err = deps.VCSService.WriteDataToFS(deps.CurrentFS, vcsData); err != nil {
		return err
	}
	return ctx.IO().Out().Printf("cleaned\n")
}
