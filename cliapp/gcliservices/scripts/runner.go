package scripts

import (
	"bytes"
	"sync"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/modules"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// Runner run application scripts
type Runner struct {
	deps struct {
		CWD      string                       `argument:"?cwd"`
		Template gcliservices.TemplateService `dependency:"TemplateService"`
		Terminal modules.Terminal             `dependency:"TerminalService"`
	}
	runned   map[string]bool
	runnedMU sync.RWMutex
}

// RunnerFactory create new Runner instance
func RunnerFactory(dp dependency.Provider) (result interface{}, err error) {
	r := &Runner{
		runned: make(map[string]bool),
	}
	if err = dp.InjectTo(&r.deps); err != nil {
		return nil, err
	}
	return gcliservices.ScriptsRunner(r), nil
}

// RunByName script by name
func (runner *Runner) RunByName(ctx app.IOContext, params gcliservices.ScriptsRunnerParams) (err error) {
	var (
		executor gcliservices.TemplateExecutor
		isEmpty  bool
		path     = scriptsBasePath + params.ScriptName
		buffer   = &bytes.Buffer{}
	)
	if !scriptNamePattern.MatchString(params.ScriptName) {
		return goaterr.Errorf("'%s' script name is incorrect", params.ScriptName)
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
		AM:        params.Data.AM,
		PlainData: params.Data.Plain,
		Properties: TaskProperties{
			Project: params.Properties,
			Secrets: params.Secrets,
		},
	}); err != nil {
		return err
	}
	parentScope := ctx.Scope()
	childCtx := gio.NewChildIOContext(ctx, gio.ChildIOContextParams{
		Scope: scope.ChildParams{
			DataScope:  parentScope,
			EventScope: parentScope,
		},
		IO: gio.IOParams{
			In: gio.NewInput(buffer),
		},
	})
	defer childCtx.Scope().Close()
	childCtx.Scope().On(app.CommitEvent, func(interface{}) (err error) {
		return runner.markAsExecutedPermanently(params.ScriptName)
	})
	return runner.deps.Terminal.RunLoop(childCtx)
}

func (runner *Runner) markAsExecutedPermanently(name string) (err error) {
	//
}
