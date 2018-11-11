package template

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/template/tfunc"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestTFuncMathSum(t *testing.T) {
	var (
		err      error
		template = `{{sum "16" "32"}}`
	)
	t.Parallel()
	// prepare mockup application & data
	output := new(bytes.Buffer)
	mapp, err := mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("")),
		Output: gio.NewOutput(output),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/build/templates/app/main.tmpl", []byte(template), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = tfunc.RegisterFunctions(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		TemplateService services.TemplateService `dependency:"TemplateService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	var executor services.TemplateExecutor
	if executor, err = deps.TemplateService.Build(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	buf := new(bytes.Buffer)
	if err = executor.Execute("default", "app", buf, nil); err != nil {
		t.Error(err)
		return
	}
	result := buf.String()
	if result != "48" {
		t.Errorf("sum function should sum string values: expected 48 and take %s", result)
		return
	}
}