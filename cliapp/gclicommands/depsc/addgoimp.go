package depsc

import (
	"github.com/goatcms/goatcli/cliapp/gclicommands/depsc/godependencies"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
)

// RunAddGOImportsDep run deps:add:go:imports command
func RunAddGOImportsDep(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			CWD    string `argument:"?cwd" ,command:"?cwd"`
			LogLvl string `argument:"?loglvl" ,command:"?loglvl"`

			Dependencies gcliservices.DependenciesService `dependency:"DependenciesService"`
		}
		importer *godependencies.Importer
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctx.Scope().InjectTo(&deps); err != nil {
		return err
	}
	if deps.CWD == "" {
		deps.CWD = "./"
	}
	if importer, err = godependencies.NewImporter(deps.CWD, godependencies.ImporterLogs{
		GOPath: func(path string) {
			ctx.IO().Out().Printf("GOPATH: %s\n", path)
		},
		OnNewSource: func(path string) {
			ctx.IO().Out().Printf("New source: %s\n", path)
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
