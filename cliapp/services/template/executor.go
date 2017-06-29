package template

import (
	"html/template"
	"io"

	"github.com/goatcms/goatcore/goathtml/ghprovider"
)

// Executor is global template provider
type Executor struct {
	provider *ghprovider.Provider
}

// Execute run template by name
func (executor *Executor) Execute(layoutName, viewName string, wr io.Writer, data interface{}) (err error) {
	var view *template.Template
	if view, err = executor.provider.View(layoutName, viewName, nil); err != nil {
		return err
	}
	return view.Execute(wr, data)
}
