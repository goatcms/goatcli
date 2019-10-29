package gtprovider

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/goatcms/goatcore/filesystem/filespace/memfs"
	"github.com/goatcms/goatcore/goattext"
	"github.com/goatcms/goatcore/workers"
)

const (
	masterOnceTemplate  = `Names: {{template "/file1.txt.once" .}}`
	overlayOnceTemplate = `{{join . ", "}}`
)

func TestLoadOnceFile(t *testing.T) {
	t.Parallel()
	var (
		funcs     = template.FuncMap{"join": strings.Join}
		guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)
	fs, err := memfs.NewFilespace()
	if err != nil {
		t.Error(err)
	}
	// create test data
	if err := fs.MkdirAll("layouts/", 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.MkdirAll("templates/", 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("layouts/default/main.tmpl", []byte(masterOnceTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("templates/myview/file1.txt.once", []byte(overlayOnceTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		assetsProvider := NewAssetsProvider(fs, "helpers/", "layouts/{name}/", funcs, true)
		provider := NewTemplatesProvider(assetsProvider, fs, "templates/{name}/", true)
		view, errs := provider.Template(goattext.DefaultLayout, "myview")
		if errs != nil {
			t.Errorf("Errors: %v", errs)
			return
		}
		buf := new(bytes.Buffer)
		if err := view.Execute(buf, guardians); err != nil {
			t.Error(err)
			return
		}
		result := buf.String()
		if !strings.Contains(result, "Gamora,") {
			t.Errorf("layout template should be overwrited. Result is: %v", result)
			return
		}
	}
}
