package template

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/template/simpletf"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestTFuncError(t *testing.T) {
	var (
		err      error
		template = `{{error "error message"}}`
	)
	t.Parallel()
	// prepare mockup application & data
	mapp, err := mockupapp.NewApp(mockupapp.MockupOptions{
		Input: gio.NewInput(strings.NewReader("")),
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
	if err = simpletf.RegisterFunctions(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		TemplateService services.TemplateService `dependency:"TemplateService"`
		Config          services.TemplateConfig  `dependency:"TemplateConfig"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	var executor services.TemplatesExecutor
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
