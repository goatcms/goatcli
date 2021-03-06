package template

import (
	"bytes"
	"testing"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcore/app/goatapp"
)

func TestTFuncMathSum(t *testing.T) {
	var (
		err      error
		template = `{{sum "16" "32"}}`
	)
	t.Parallel()
	// prepare mockup application & data
	mapp, err := goatapp.NewMockupApp(goatapp.Params{})
	if err != nil {
		t.Error(err)
		return
	}
	fs := mapp.Filespaces().CWD()
	if err = fs.WriteFile(".goat/build/templates/app/main.tmpl", []byte(template), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = simpletf.RegisterFunctions(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		TemplateService gcliservices.TemplateService `dependency:"TemplateService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	var executor gcliservices.TemplatesExecutor
	if executor, err = deps.TemplateService.TemplatesExecutor(); err != nil {
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
