package cleanc

import (
	"github.com/goatcms/goatcore/app"
)

// RunClean run clean command
func RunClean(a app.App, ctxScope app.Scope) (err error) {
	if err = RunCleanBuild(a, ctxScope); err != nil {
		return err
	}
	return RunCleanDependencies(a, ctxScope)
}
