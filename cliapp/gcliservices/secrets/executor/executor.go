package executor

import (
	"fmt"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/cutil"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/workers"
)

// SecretsExecutor is code generator task manager
type SecretsExecutor struct {
	templateExecutor gcliservices.TemplateExecutor
	sharedData       SharedData
	limit            int64
	ch               chan Task
	scope            app.Scope
	context          *Context
}

// NewSecretsExecutor create new SecretsExecutor instance
func NewSecretsExecutor(ctxScope app.Scope, sharedData SharedData, limit int64, templateExecutor gcliservices.TemplateExecutor) (instance *SecretsExecutor, err error) {
	instance = &SecretsExecutor{
		templateExecutor: templateExecutor,
		sharedData:       sharedData,
		limit:            limit,
		ch:               make(chan Task, limit),
		scope:            ctxScope,
	}
	instance.context = newContext(instance, sharedData)
	for i := workers.MaxJob; i > 0; i-- {
		go instance.consumer()
	}
	return instance, nil
}

// Scope run executor scope
func (e *SecretsExecutor) Scope() app.Scope {
	return e.scope
}

// Secrets return context secrets
func (e *SecretsExecutor) Secrets() (secrets []*config.Property, err error) {
	return e.context.secrets, nil
}

// Execute run single executor template
func (e *SecretsExecutor) Execute() (err error) {
	var list []string
	// Execute all single template
	if list, err = e.templateExecutor.Templates(); err != nil {
		return err
	}
	for _, name := range list {
		if e.scope.IsDone() {
			return nil
		}
		if strings.HasSuffix(name, ".ctrl") {
			e.ExecuteTask(Task{
				TemplateName: name,
			})
		}
	}
	return nil
}

// ExecuteTask run single executor template
func (e *SecretsExecutor) ExecuteTask(task Task) (err error) {
	if len(e.ch) == cap(e.ch) {
		return fmt.Errorf("SecretsExecutor: channel is full")
	}
	e.scope.AddTasks(1)
	e.ch <- task
	return nil
}

// ExecuteTask run single executor template
func (e *SecretsExecutor) consumer() (err error) {
	for {
		select {
		case task, more := <-e.ch:
			if !more {
				return
			}
			if e.scope.IsDone() {
				e.scope.DoneTask()
				continue
			}
			if err = e.run(task); err != nil {
				e.scope.AppendError(err)
			}
			e.scope.DoneTask()
		}
	}
}

func (e *SecretsExecutor) run(task Task) (err error) {
	if e.scope.IsDone() {
		return fmt.Errorf("Context killed")
	}
	writer := cutil.NewNilWriter()
	if task.TemplateName == "" {
		if isEmpty, err := e.templateExecutor.IsEmpty(); isEmpty || err != nil {
			return err
		}
		return e.templateExecutor.Execute(writer, e.context)
	}
	return e.templateExecutor.ExecuteTemplate(task.TemplateName, writer, e.context)
}
