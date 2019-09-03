package executor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/goatcms/goatcore/filesystem"
)

// Context contains template dot object with data and APIs
type Context struct {
	Template   TemplateHandler
	DotData    interface{}
	PlainData  map[string]string
	Properties TaskProperties
	AM         interface{}

	// internal variables
	fs       filesystem.Filespace
	executor *GeneratorExecutor
}

func (c *Context) Error(msg string) (err error) {
	err = fmt.Errorf(msg)
	c.executor.Scope().AppendError(err)
	return err
}

// RenderOnce render a template to new file if the file does't exist
func (c *Context) RenderOnce(destPath, layout, path, name string, dotData interface{}) (err error) {
	if c.fs.IsExist(destPath) {
		return nil
	}
	return c.Render(destPath, layout, path, name, dotData)
}

// Render render a template to file
func (c *Context) Render(destPath, layout, path, name string, dotData interface{}) (err error) {
	var (
		task     Task
		ctxScope = c.executor.Scope()
	)
	if ctxScope.IsKilled() {
		return ctxScope.ToError()
	}
	destPath = filepath.Clean(destPath)
	if strings.HasPrefix(destPath, "hook/") || destPath == "hook" {
		return fmt.Errorf("'hook' derecotry is reserved and can not be run from template ")
	}
	task = Task{
		Template: TemplateHandler{
			Layout: layout,
			Path:   path,
			Name:   name,
		},
		DotData:         dotData,
		BuildProperties: c.Properties.Build,
		FSPath:          destPath,
	}
	if layout == "" {
		task.Template.Layout = c.Template.Layout
	}
	if path == "" {
		task.Template.Path = c.Template.Path
	}
	if name == "" {
		task.Template.Name = c.Template.Name
	}
	return c.executor.ExecuteTask(task)
}
