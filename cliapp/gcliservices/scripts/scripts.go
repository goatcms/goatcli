package scripts

import (
	"bytes"
	"fmt"

	"github.com/goatcms/goatcore/varutil"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/modules"
	"github.com/goatcms/goatcore/app/modules/commonm/commservices"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipservices"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipservices/namespaces"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// Runner run application scripts
type Runner struct {
	deps struct {
		CWD       string                       `argument:"?cwd"`
		Template  gcliservices.TemplateService `dependency:"TemplateService"`
		Terminal  modules.Terminal             `dependency:"TerminalService"`
		Runner    pipservices.Runner           `dependency:"PipRunner"`
		TasksUnit pipservices.TasksUnit        `dependency:"PipTasksUnit"`
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
		executor    gcliservices.TemplateExecutor
		isEmpty     bool
		path        = scriptsBasePath + scriptName
		buffer      = &bytes.Buffer{}
		taskManager pipservices.TasksManager
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
	ctxIO := ctx.IO()
	if err = runner.deps.Runner.Run(pipservices.Pip{
		Context: pipservices.PipContext{
			In:    gio.NewInput(buffer),
			Out:   gio.NewNilOutput(),
			Err:   gio.NewNilOutput(),
			CWD:   ctxIO.CWD(),
			Scope: ctx.Scope(),
		},
		Name:        "MAIN",
		Description: fmt.Sprintf("%s body", scriptName),
		Namespaces: namespaces.NewNamespaces(pipservices.NamasepacesParams{
			Task: "",
			Lock: "",
		}),
		Sandbox: "self",
		Lock:    commservices.LockMap{},
		Wait:    []string{},
	}); err != nil {
		return err
	}
	if taskManager, err = runner.deps.TasksUnit.FromScope(ctx.Scope()); err != nil {
		return err
	}
	return waitForTasks(taskManager, ctx)
}

func waitForTasks(taskManager pipservices.TasksManager, ctx app.IOContext) (err error) {
	var (
		ended = []string{}
		names = taskManager.Names()
		out   = ctx.IO().Out()
		task  pipservices.Task
		ok    bool
	)
	for len(names) != len(ended) {
		for _, name := range names {
			if varutil.IsArrContainStr(ended, name) {
				continue
			}
			ended = append(ended, name)
			if task, ok = taskManager.Get(name); !ok {
				return goaterr.Errorf("Unknow task %s", name)
			}
			out.Printf("\n Wait for %s: %s ", name, task.Description())
			if err = task.Wait(); err != nil {
				return err
			}
			out.Printf("\n [%s] %s... %s", name, task.Description(), task.Status())
			for _, err = range task.Errors() {
				var (
					merr goaterr.MessageError
					ok   bool
				)
				if merr, ok = err.(goaterr.MessageError); ok {
					out.Printf("\n - %s", merr.Message())
				} else {
					out.Printf("\n - %s", err.Error())
				}
			}
		}
		names = taskManager.Names()
	}
	return nil
}
