package template

import (
	"io"
	"text/template"

	"github.com/goatcms/goatcore/goattext/gtprovider"
)

// Executor is global template provider
type Executor struct {
	provider *gtprovider.Provider
}

// Execute run main code of a view
func (executor *Executor) Execute(layoutName, viewName string, wr io.Writer, data interface{}) (err error) {
	var view *template.Template
	if view, err = executor.provider.View(layoutName, viewName, nil); err != nil {
		return err
	}
	return view.Execute(wr, data)
}

// ExecuteTemplate run template of a view
func (executor *Executor) ExecuteTemplate(layoutName, viewName, templateName string, wr io.Writer, data interface{}) (err error) {
	var view *template.Template
	if view, err = executor.provider.View(layoutName, viewName, nil); err != nil {
		return err
	}
	return view.ExecuteTemplate(wr, templateName, data)
}
