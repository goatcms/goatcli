package template

import (
	"io"
	"text/template"

	"github.com/goatcms/goatcli/cliapp/services/template/gtprovider"
)

// TemplateExecutor is single template executor
type TemplateExecutor struct {
	provider *gtprovider.TemplateProvider
}

// NewTemplateExecutor create NewTemplateExecutor instance
func NewTemplateExecutor(provider *gtprovider.TemplateProvider) *TemplateExecutor {
	return &TemplateExecutor{
		provider: provider,
	}
}

// Execute run main code of a view
func (executor *TemplateExecutor) Execute(wr io.Writer, data interface{}) (err error) {
	var view *template.Template
	if view, err = executor.provider.Template(); err != nil {
		return err
	}
	return view.Execute(wr, data)
}
