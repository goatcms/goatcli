package template

import (
	"io"
	"text/template"

	"github.com/goatcms/goatcli/cliapp/services/template/gtprovider"
)

// Executor is global template provider
type Executor struct {
	provider *gtprovider.Provider
}

// Templates return template list for a view
func (executor *Executor) Templates(layoutName, viewName string) (list []string, err error) {
	var view *template.Template
	if view, err = executor.provider.View(layoutName, viewName); err != nil {
		return nil, err
	}
	list = []string{}
	for _, t := range view.Templates() {
		list = append(list, t.Name())
	}
	return list, nil
}

// Execute run main code of a view
func (executor *Executor) Execute(layoutName, viewName string, wr io.Writer, data interface{}) (err error) {
	var view *template.Template
	if view, err = executor.provider.View(layoutName, viewName); err != nil {
		return err
	}
	return view.Execute(wr, data)
}

// ExecuteTemplate run template of a view
func (executor *Executor) ExecuteTemplate(layoutName, viewName, templateName string, wr io.Writer, data interface{}) (err error) {
	var view *template.Template
	if view, err = executor.provider.View(layoutName, viewName); err != nil {
		return err
	}
	return view.ExecuteTemplate(wr, templateName, data)
}
