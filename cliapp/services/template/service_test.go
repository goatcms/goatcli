package template

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

const (
	testLayoutContent   = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
	testTemplateContent = `{{define "list"}} {{join . ", "}}{{end}} `
)

func TestServiceExecute(t *testing.T) {
	var (
		err       error
		guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
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
	if err = mapp.RootFilespace().WriteFile(".goat/build/layouts/default/main.tmpl", []byte(testLayoutContent), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/build/templates/app/main.tmpl", []byte(testTemplateContent), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
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
	if err = deps.TemplateService.AddFunc("join", strings.Join); err != nil {
		t.Error(err)
		return
	}
	var executor services.TemplateExecutor
	if executor, err = deps.TemplateService.Build(mapp.RootFilespace()); err != nil {
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
