package template

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestAssetsServiceLayout(t *testing.T) {
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
	if err = mapp.RootFilespace().WriteFile(".goat/build/helpers/main.tmpl", []byte(`
		Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}
	`), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/build/layouts/layoutName/main.tmpl", []byte(`
		{{define "list"}} {{join . ", "}}{{end}}
	`), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		AssetsProvider services.TemplateAssetsProvider `dependency:"TemplateAssetsProvider"`
		Config         services.TemplateConfig         `dependency:"TemplateConfig"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if err = deps.Config.AddFunc("join", strings.Join); err != nil {
		t.Error(err)
		return
	}
	var template *template.Template
	if template, err = deps.AssetsProvider.Layout("layoutName"); err != nil {
		t.Error(err)
		return
	}
	buf := new(bytes.Buffer)
	if err = template.Execute(buf, guardians); err != nil {
		t.Error(err)
		return
	}
	result := buf.String()
	if !strings.Contains(result, "Gamora,") {
		t.Errorf("after render overlay should contains 'Gamora,' and other characters. It is: %s", result)
		return
	}
}
