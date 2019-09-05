package vcsc

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunClean run vcs:clean command. It remove path from ignored files foreach doeas't exist file
func RunClean(a app.App, ctxScope app.Scope) (err error) {
	var (
		deps struct {
			CurrentFS filesystem.Filespace `filespace:"current"`

			VCSService services.VCSService `dependency:"VCSService"`
			Input      app.Input           `dependency:"InputService"`
			Output     app.Output          `dependency:"OutputService"`
		}
		vcsData services.VCSData
		ignored = vcs.NewIgnoredFiles(true)
	)
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctxScope.InjectTo(&deps))); err != nil {
		return err
	}
	if vcsData, err = deps.VCSService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	for _, path := range vcsData.VCSIgnoredFiles().All() {
		if deps.CurrentFS.IsFile(path) {
			ignored.AddPath(path)
		} else {
			deps.Output.Printf(" Deleted from ignored: %s\n", path)
		}
	}
	if len(ignored.All()) == len(vcsData.VCSIgnoredFiles().All()) {
		deps.Output.Printf("ignored files are clean\n")
		return nil
	}
	vcsData = vcs.NewData(vcsData.VCSGeneratedFiles(), ignored)
	if err = deps.VCSService.WriteDataToFS(deps.CurrentFS, vcsData); err != nil {
		return err
	}
	deps.Output.Printf("cleaned\n")
	return nil
}
