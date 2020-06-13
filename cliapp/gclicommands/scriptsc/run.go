package scriptsc

import (
	"time"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipcommands/pipc"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunScript run script by name
func RunScript(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Name            string                       `command:"?$1"`
			ScriptsRunner   gcliservices.ScriptsRunner   `dependency:"ScriptsRunner"`
			GCLIInputs      gcliservices.GCLIInputs      `dependency:"GCLIInputs"`
			GCLIEnvironment gcliservices.GCLIEnvironment `dependency:"GCLIEnvironment"`
		}
		propertiesData map[string]string
		secretsData    map[string]string
		appData        gcliservices.ApplicationData
		ctxScope       = ctx.Scope()
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctxScope.InjectTo(&deps); err != nil {
		return err
	}
	if deps.Name == "" {
		return goaterr.Errorf(FirstKeyParameterIsRequired)
	}
	// load variables
	if propertiesData, secretsData, appData, err = deps.GCLIInputs.Inputs(ctx); err != nil {
		return err
	}
	if err = deps.GCLIEnvironment.LoadEnvs(ctxScope, propertiesData, secretsData); err != nil {
		return err
	}
	// run script
	if err = deps.ScriptsRunner.Run(ctx, ctx.IO().CWD(), deps.Name, propertiesData, secretsData, appData); err != nil {
		return err
	}
	// wait
	return runScriptSaveLogs(deps.Name, a, ctx)
}

func runScriptSaveLogs(scriptName string, a app.App, ctx app.IOContext) (err error) {
	var (
		logsFileWriter    filesystem.Writer
		logsOutput        app.Output
		summaryFileWriter filesystem.Writer
		summaryOutput     app.Output
		now               = time.Now()
		basename          = ".goat/tmp/logs/scripts/" + now.Format("2006-01-02-15:04:05") + "-" + scriptName + "."
		cwd               = ctx.IO().CWD()
	)
	if err = cwd.MkdirAll(".goat/tmp/logs/scripts", filesystem.DefaultUnixDirMode); err != nil {
		return err
	}
	// write logs
	if logsFileWriter, err = cwd.Writer(basename + "logs.txt"); err != nil {
		return err
	}
	defer logsFileWriter.Close()
	logsOutput = gio.NewOutput(logsFileWriter)
	logsCtx := gio.NewChildIOContext(ctx, gio.ChildIOContextParams{
		Scope: scope.ChildParams{},
		IO: gio.IOParams{
			Out: logsOutput,
			Err: logsOutput,
		},
	})
	defer logsCtx.Close()
	if err = pipc.Logs(a, logsCtx); err != nil {
		return err
	}
	// write summary
	if summaryFileWriter, err = cwd.Writer(basename + "summary.txt"); err != nil {
		return err
	}
	defer summaryFileWriter.Close()
	summaryOutput = gio.NewOutput(summaryFileWriter)
	summaryCtx := gio.NewChildIOContext(ctx, gio.ChildIOContextParams{
		Scope: scope.ChildParams{},
		IO: gio.IOParams{
			Out: summaryOutput,
			Err: summaryOutput,
		},
	})
	defer summaryCtx.Close()
	return pipc.Summary(a, summaryCtx)
}
