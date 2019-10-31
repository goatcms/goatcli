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

// Templates return template list for a view
func (executor *TemplateExecutor) Templates() (list []string, err error) {
	var view *template.Template
	if view, err = executor.provider.Template(); err != nil {
		return nil, err
	}
	list = []string{}
	for _, t := range view.Templates() {
		list = append(list, t.Name())
	}
	return list, nil
}

// IsEmpty run true if template is not empty
func (executor *TemplateExecutor) IsEmpty() (is bool, err error) {
	var tmpl *template.Template
	if tmpl, err = executor.provider.Template(); err != nil {
		return false, nil
	}
	return tmpl.Tree == nil || tmpl.Root == nil, nil
}

// Execute run main code of a view
func (executor *TemplateExecutor) Execute(wr io.Writer, data interface{}) (err error) {
	var view *template.Template
	if view, err = executor.provider.Template(); err != nil {
		return err
	}
	return view.Execute(wr, data)
}

// ExecuteTemplate run template of a view
func (executor *TemplateExecutor) ExecuteTemplate(templateName string, wr io.Writer, data interface{}) (err error) {
	var template *template.Template
	if template, err = executor.provider.Template(); err != nil {
		return err
	}
	return template.ExecuteTemplate(wr, templateName, data)
}
