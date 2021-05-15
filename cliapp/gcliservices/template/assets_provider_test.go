package template

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app/goatapp"
)

func TestAssetsServiceLayout(t *testing.T) {
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
	if err = fs.WriteFile(".goat/build/helpers/main.tmpl", []byte(`
		Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}
	`), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile(".goat/build/layouts/layoutName/main.tmpl", []byte(`
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
		AssetsProvider gcliservices.TemplateAssetsProvider `dependency:"TemplateAssetsProvider"`
		Config         gcliservices.TemplateConfig         `dependency:"TemplateConfig"`
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
