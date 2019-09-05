package vcsc

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunIgnoredAdd run vcs:ignored:add command
func RunIgnoredAdd(a app.App, ctxScope app.Scope) (err error) {
	var (
		deps struct {
			Path string `argument:"?$2" ,command:"?$2"`

			CurrentFS filesystem.Filespace `filespace:"current"`

			VCSService services.VCSService `dependency:"VCSService"`
			Input      app.Input           `dependency:"InputService"`
			Output     app.Output          `dependency:"OutputService"`
		}
		vcsData services.VCSData
	)
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctxScope.InjectTo(&deps))); err != nil {
		return err
	}
	if deps.Path == "" {
		return goaterr.NewError("Path is required")
	}
	if vcsData, err = deps.VCSService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	vcsData.VCSIgnoredFiles().AddPath(deps.Path)
	if err = deps.VCSService.WriteDataToFS(deps.CurrentFS, vcsData); err != nil {
		return err
	}
	deps.Output.Printf("Path added\n")
	return nil
}
