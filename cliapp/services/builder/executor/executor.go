package executor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goatcms/goatcli/cliapp/common/cutil"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/workers"
)

// GeneratorExecutor is code generator task manager
type GeneratorExecutor struct {
	templatesExecutor services.TemplatesExecutor
	sharedData        SharedData
	limit             int64
	ch                chan Task
	scope             app.Scope
}

// NewGeneratorExecutor create new GeneratorExecutor instance
func NewGeneratorExecutor(ctxScope app.Scope, sharedData SharedData, limit int64, templatesExecutor services.TemplatesExecutor) (instance *GeneratorExecutor, err error) {
	instance = &GeneratorExecutor{
		templatesExecutor: templatesExecutor,
		sharedData:        sharedData,
		limit:             limit,
		ch:                make(chan Task, limit),
		scope:             ctxScope,
	}
	for i := workers.MaxJob; i > 0; i-- {
		go instance.consumer()
	}
	return instance, nil
}

// Scope run executor scope
func (e *GeneratorExecutor) Scope() app.Scope {
	return e.scope
}

// ExecuteView run single executor template
func (e *GeneratorExecutor) ExecuteView(layout, viewPath string, properties map[string]string, dotData interface{}) (err error) {
	var list []string
	// Execute all single template
	if list, err = e.templatesExecutor.Templates(layout, viewPath); err != nil {
		return err
	}
	for _, name := range list {
		var (
			extension, destPath string
		)
		if e.scope.IsKilled() {
			return nil
		}
		extension = filepath.Ext(name)
		destPath = name[0 : len(name)-len(extension)]
		if strings.HasSuffix(name, ".once") {
			if e.sharedData.FS.IsExist(destPath) {
				continue
			}
			e.ExecuteTask(Task{
				Template: TemplateHandler{
					Layout: layout,
					Path:   viewPath,
					Name:   name,
				},
				DotData:         dotData,
				BuildProperties: properties,
				FSPath:          destPath,
			})
		} else if strings.HasSuffix(name, ".render") {
			e.sharedData.FS.Remove(destPath)
			e.ExecuteTask(Task{
				Template: TemplateHandler{
					Layout: layout,
					Path:   viewPath,
					Name:   name,
				},
				DotData:         dotData,
				BuildProperties: properties,
				FSPath:          destPath,
			})
		} else if strings.HasSuffix(name, ".ctrl") {
			e.ExecuteTask(Task{
				Template: TemplateHandler{
					Layout: layout,
					Path:   viewPath,
					Name:   name,
				},
				DotData:         dotData,
				BuildProperties: properties,
			})
		}
	}
	return nil
}

// ExecuteTask run single executor template
func (e *GeneratorExecutor) ExecuteTask(task Task) (err error) {
	if len(e.ch) == cap(e.ch) {
		return fmt.Errorf("GeneratorExecutor: channel is full")
	}
	e.scope.AddTasks(1)
	e.ch <- task
	return nil
}

// ExecuteHook run single hook templates
func (e *GeneratorExecutor) ExecuteHook(name string, data interface{}) (err error) {
	var nodes []os.FileInfo
	fs := e.sharedData.FS
	path := ".goat/build/templates/hook/" + name
	if !fs.IsDir(path) {
		return nil
	}
	if nodes, err = fs.ReadDir(path); err != nil {
		return err
	}
	for _, node := range nodes {
		if err = e.ExecuteView("", "hook/"+name+"/"+node.Name(), map[string]string{}, data); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteTask run single executor template
func (e *GeneratorExecutor) consumer() (err error) {
	generatedFileds := e.sharedData.VCSData.VCSGeneratedFiles()
	for {
		select {
		case task, more := <-e.ch:
			if !more {
				return
			}
			if e.scope.IsKilled() {
				e.scope.DoneTask()
				continue
			}
			if err = e.run(task); err != nil {
				e.scope.AppendError(err)
			}
			e.scope.DoneTask()
			if task.FSPath != "" {
				generatedFileds.Add(&services.GeneratedFile{
					Path:    task.FSPath,
					ModTime: time.Now(),
				})
			}
		}
	}
}

func (e *GeneratorExecutor) executeToWriter(writer io.Writer, task Task) (err error) {
	if e.scope.IsKilled() {
		return fmt.Errorf("Context killed")
	}
	sharedData := e.sharedData
	ctx := &Context{
		AM:        sharedData.AM,
		Template:  task.Template,
		DotData:   task.DotData,
		PlainData: sharedData.PlainData,
		Properties: TaskProperties{
			Build:   task.BuildProperties,
			Project: sharedData.Properties.Project,
			Secrets: sharedData.Properties.Secrets,
		},
		fs:       sharedData.FS,
		executor: e,
	}
	tmpl := task.Template
	if tmpl.Name == "" {
		return e.templatesExecutor.Execute(tmpl.Layout, tmpl.Path, writer, ctx)
	}
	return e.templatesExecutor.ExecuteTemplate(tmpl.Layout, tmpl.Path, tmpl.Name, writer, ctx)
}

func (e *GeneratorExecutor) run(task Task) (err error) {
	var writer filesystem.Writer
	if e.scope.IsKilled() {
		return fmt.Errorf("Context is killed")
	}
	if task.FSPath == "" {
		writer = cutil.NewNilWriter()
	} else {
		if writer, err = e.sharedData.FS.Writer(task.FSPath); err != nil {
			return err
		}
		defer writer.Close()
	}
	return e.executeToWriter(writer, task)
}
