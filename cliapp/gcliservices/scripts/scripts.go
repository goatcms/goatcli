package scripts

import (
	"bytes"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/modules"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// Runner run application scripts
type Runner struct {
	deps struct {
		CWD      string                       `argument:"?cwd"`
		Template gcliservices.TemplateService `dependency:"TemplateService"`
		Terminal modules.Terminal             `dependency:"TerminalService"`
	}
}

// RunnerFactory create new Runner instance
func RunnerFactory(dp dependency.Provider) (result interface{}, err error) {
	r := &Runner{}
	if err = dp.InjectTo(&r.deps); err != nil {
		return nil, err
	}
	return gcliservices.ScriptsRunner(r), nil
}

// Run script by name
func (runner *Runner) Run(ctx app.IOContext, fs filesystem.Filespace, scriptName string, properties, secrets map[string]string, appData gcliservices.ApplicationData) (err error) {
	var (
		executor gcliservices.TemplateExecutor
		isEmpty  bool
		path     = scriptsBasePath + scriptName
		buffer   = &bytes.Buffer{}
	)
	if !scriptNamePattern.MatchString(scriptName) {
		return goaterr.Errorf("'%s' script name is incorrect", scriptName)
	}
	if executor, err = runner.deps.Template.TemplateExecutor(path); err != nil {
		return err
	}
	if isEmpty, err = executor.IsEmpty(); err != nil {
		return err
	}
	if isEmpty {
		return goaterr.Errorf("The script %s is empty", path)
	}
	if err = executor.Execute(buffer, Context{
		AM:        appData.AM,
		PlainData: appData.Plain,
		Properties: TaskProperties{
			Project: properties,
			Secrets: secrets,
		},
	}); err != nil {
		return err
	}
	childCtx := gio.NewChildIOContext(ctx, gio.NewInput(buffer), nil, nil, nil)
	defer childCtx.Scope().Close()
	return runner.deps.Terminal.RunLoop(childCtx)
}
