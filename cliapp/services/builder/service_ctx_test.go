package builder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/template"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

const (
	testCTXBuilderLayout   = `{{define "list"}} {{range .}}{{println "," .}}{{end}} {{end}}`
	testCTXBuilderTemplate = `
	{{if not (.Filesystem.IsFile "out/file.txt")}}
		{{.Out.File "out/file.txt"}}
			{{- index .Data "datakey" -}}
			{{- index .Properties.Project "propkey" -}}
			{{- index .Properties.Secrets "secretkey" -}}
		{{.Out.EOF}}
	{{end}}`
)

func TestCTXBuilder(t *testing.T) {
	var (
		err     error
		context []byte
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
	if err = mapp.RootFilespace().WriteFile(".goat/build/layouts/default/main.tmpl", []byte(testCTXBuilderLayout), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/build/templates/names/main.tmpl", []byte(testCTXBuilderTemplate), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = template.RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = template.InitDependencies(mapp); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		BuilderService services.BuilderService `dependency:"BuilderService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	buildConfig := []*config.Build{
		&config.Build{
			From:     "ignore",
			To:       "ignore",
			Template: "names",
			Layout:   "default",
		},
	}
	fs := mapp.RootFilespace()
	if err = deps.BuilderService.Build(fs, buildConfig, map[string]string{
		"datakey": "Ala",
	}, map[string]string{
		"propkey": " ma",
	}, map[string]string{
		"secretkey": " kota",
	}); err != nil {
		t.Error(err)
		return
	}
	if !fs.IsFile("out/file.txt") {
		t.Errorf("out/file.txt is not exist")
		return
	}
	if context, err = fs.ReadFile("out/file.txt"); err != nil {
		t.Error(err)
		return
	}
	if string(context) != "Ala ma kota" {
		t.Errorf("File content must be equals to 'Ala ma kota' and it is '%s'", context)
		return
	}
}
