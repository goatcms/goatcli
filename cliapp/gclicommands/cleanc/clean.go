package cleanc

import (
	"github.com/goatcms/goatcore/app"
)

// RunClean run clean command
func RunClean(a app.App, ctx app.IOContext) (err error) {
	if err = RunCleanBuild(a, ctx); err != nil {
		return err
	}
	return RunCleanDependencies(a, ctx)
}
