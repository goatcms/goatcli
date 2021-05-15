package scripts

import (
	"bytes"
	"fmt"

	"github.com/goatcms/goatcli/cliapp/common"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/modules/commonm/commservices"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipservices"
	"github.com/goatcms/goatcore/app/modules/terminalm/termservices"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// Runner run application scripts
type Runner struct {
	deps struct {
		CWD       string                       `argument:"?cwd"`
		Template  gcliservices.TemplateService `dependency:"TemplateService"`
		Terminal  termservices.Terminal        `dependency:"TerminalService"`
		Runner    pipservices.Runner           `dependency:"PipRunner"`
		TasksUnit pipservices.TasksUnit        `dependency:"PipTasksUnit"`
	}
}

// RunnerFactory create new Runner instance
func RunnerFactory(dp app.DependencyProvider) (result interface{}, err error) {
	r := &Runner{}
	if err = dp.InjectTo(&r.deps); err != nil {
		return nil, err
	}
	return gcliservices.ScriptsRunner(r), nil
}

// Run script by name
func (runner *Runner) Run(ctx gcliservices.ScriptsContext, fs filesystem.Filespace, scriptName string, properties, secrets common.ElasticData, appData gcliservices.ApplicationData) (taskManager pipservices.TasksManager, err error) {
	var (
		executor gcliservices.TemplateExecutor
		isEmpty  bool
		path     = scriptsBasePath + scriptName
		buffer   = &bytes.Buffer{}
	)
	if !scriptNamePattern.MatchString(scriptName) {
		return nil, goaterr.Errorf("'%s' script name is incorrect", scriptName)
	}
	if executor, err = runner.deps.Template.TemplateExecutor(path); err != nil {
		return nil, err
	}
	if isEmpty, err = executor.IsEmpty(); err != nil {
		return nil, err
	}
	if isEmpty {
		return nil, goaterr.Errorf("The script %s is empty", path)
	}
	if err = executor.Execute(buffer, Context{
		AM:   appData.AM,
		Data: appData.ElasticData,
		Properties: TaskProperties{
			Project: properties,
			Secrets: secrets,
		},
	}); err != nil {
		return nil, err
	}
	childScope := scope.NewChild(ctx.Scope, scope.ChildParams{})
	if err = runner.deps.Runner.Run(pipservices.Pip{
		Context: pipservices.PipContext{
			In:    gio.NewInput(buffer),
			Out:   gio.NewNilOutput(),
			Err:   gio.NewNilOutput(),
			CWD:   ctx.CWD,
			Scope: childScope,
		},
		Name:        scriptName,
		Description: fmt.Sprintf("%s body", scriptName),
		Namespaces:  ctx.Namespaces,
		Sandbox:     "self",
		Lock:        commservices.LockMap{},
		Wait:        []string{},
	}); err != nil {
		return nil, err
	}
	go childScope.Close()
	return runner.deps.TasksUnit.FromScope(childScope)
}
