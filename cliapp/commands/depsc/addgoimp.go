package depsc

import (
	"github.com/goatcms/goatcli/cliapp/commands/depsc/godependencies"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
)

// RunAddGOImportsDep run deps:add:go:imports command
func RunAddGOImportsDep(a app.App, ctxScope app.Scope) (err error) {
	var (
		deps struct {
			CWD    string `argument:"?cwd",command:"?cwd"`
			LogLvl string `argument:"?loglvl",command:"?loglvl"`

			Dependencies services.DependenciesService `dependency:"DependenciesService"`
			Input        app.Input                    `dependency:"InputService"`
			Output       app.Output                   `dependency:"OutputService"`
		}
		importer *godependencies.Importer
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctxScope.InjectTo(&deps); err != nil {
		return err
	}
	if deps.CWD == "" {
		deps.CWD = "./"
	}
	if importer, err = godependencies.NewImporter(deps.CWD, godependencies.ImporterLogs{
		GOPath: func(path string) {
			deps.Output.Printf("GOPATH: %s\n", path)
		},
		OnNewSource: func(path string) {
			deps.Output.Printf("New source: %s\n", path)
		},
	}, godependencies.ImporterOptions{
		MaxDep:  godependencies.MaxImportDepth,
		DevLogs: deps.LogLvl == "dev",
	}, deps.Dependencies); err != nil {
		return err
	}
	if err = importer.Import(); err != nil {
		return err
	}
	return importer.WriteDef()
}