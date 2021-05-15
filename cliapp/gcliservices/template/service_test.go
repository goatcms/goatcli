package template

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app/goatapp"
)

const (
	testLayoutContent   = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
	testTemplateContent = `{{define "list"}} {{join . ", "}}{{end}} `
)

func TestServiceTemplateExecute(t *testing.T) {
	var (
		err       error
		guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)
	t.Parallel()
	// prepare mockup application & data
	mapp, err := goatapp.NewMockupApp(goatapp.Params{})
	if err != nil {
		t.Error(err)
		return
	}
	fs := mapp.Filespaces().CWD()
	if err = fs.WriteFile(".goat/build/helpers/main.tmpl", []byte(testLayoutContent), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile(".goat/properties.def/main.tmpl", []byte(testTemplateContent), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
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
	if err = deps.Config.AddFunc("join", strings.Join); err != nil {
		t.Error(err)
		return
	}
	var executor gcliservices.TemplateExecutor
	if executor, err = deps.TemplateService.TemplateExecutor(".goat/properties.def"); err != nil {
		t.Error(err)
		return
	}
	buf := new(bytes.Buffer)
	if err = executor.Execute(buf, guardians); err != nil {
		t.Error(err)
		return
	}
	result := buf.String()
	if !strings.Contains(result, "Gamora,") {
		t.Errorf("after render overlay should contains 'Gamora,' and other characters. It is: %s", result)
		return
	}
}

func TestServiceViewsExecute(t *testing.T) {
	var (
		err       error
		guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)
	t.Parallel()
	// prepare mockup application & data
	mapp, err := goatapp.NewMockupApp(goatapp.Params{})
	if err != nil {
		t.Error(err)
		return
	}
	fs := mapp.Filespaces().CWD()
	if err = fs.WriteFile(".goat/build/layouts/default/main.tmpl", []byte(testLayoutContent), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile(".goat/build/templates/app/main.tmpl", []byte(testTemplateContent), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
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
	if err = deps.Config.AddFunc("join", strings.Join); err != nil {
		t.Error(err)
		return
	}
	var executor gcliservices.TemplatesExecutor
	if executor, err = deps.TemplateService.TemplatesExecutor(); err != nil {
		t.Error(err)
		return
	}
	buf := new(bytes.Buffer)
	if err = executor.Execute("default", "app", buf, guardians); err != nil {
		t.Error(err)
		return
	}
	result := buf.String()
	if !strings.Contains(result, "Gamora,") {
		t.Errorf("after render overlay should contains 'Gamora,' and other characters. It is: %s", result)
		return
	}
}
