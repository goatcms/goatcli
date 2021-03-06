package template

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcore/app/goatapp"
)

func TestTFuncError(t *testing.T) {
	var (
		err      error
		template = `{{error "error message"}}`
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
		Config          gcliservices.TemplateConfig  `dependency:"TemplateConfig"`
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
	if err = executor.Execute("default", "app", buf, nil); err == nil {
		t.Error("expected error")
		return
	}
	if !strings.Contains(err.Error(), "error message") {
		t.Errorf("expected error with message 'error message' and take %s", err.Error())
		return
	}
}
