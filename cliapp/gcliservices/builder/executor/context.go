package executor

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common"

	"github.com/goatcms/goatcore/filesystem"
)

// Context contains template dot object with data and APIs
type Context struct {
	Template   TemplateHandler
	DotData    interface{}
	Data       common.ElasticData
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

// ExecTemplate run template by name and return result as string
func (c *Context) ExecTemplate(name string, dotData interface{}) (result string, err error) {
	var buf = &bytes.Buffer{}
	if err = c.executor.templatesExecutor.ExecuteTemplate(
		c.Template.Layout,
		c.Template.Path,
		name,
		buf,
		dotData); err != nil {
		return "", err
	}
	return buf.String(), nil
}
