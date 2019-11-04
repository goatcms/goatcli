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

func TestLoadViewWithDefaultLayout(t *testing.T) {
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
	if err := fs.WriteFile("layouts/default/main.tmpl", []byte(masterTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("templates/mytemplate/main.tmpl", []byte(overlayTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		assetsProvider := NewAssetsProvider(fs, goattext.HelpersPath, goattext.LayoutPath, funcs, true)
		provider := NewTemplatesProvider(assetsProvider, fs, "templates/{name}", true)
		view, err := provider.Template(goattext.DefaultLayout, "mytemplate")
		if err != nil {
			t.Errorf("Errors: %v", err)
			return
		}
		buf := new(bytes.Buffer)
		if err := view.Execute(buf, guardians); err != nil {
			t.Error(err)
			return
		}
		result := buf.String()
		if !strings.Contains(result, "Gamora,") {
			t.Errorf("layout template should be overwrited. result: %v", result)
			return
		}
	}
}

func TestLoadManyFiles(t *testing.T) {
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
	if err := fs.WriteFile("layouts/default/main.tmpl", []byte(masterTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("templates/mytemplate/file1.tmpl", []byte(templateFile1), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("templates/mytemplate/file2.tmpl", []byte(templateFile2), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("templates/mytemplate/file3.tmpl", []byte(templateFile3), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("templates/mytemplate/file4.tmpl", []byte(templateFile4), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		assetsProvider := NewAssetsProvider(fs, goattext.HelpersPath, goattext.LayoutPath, funcs, true)
		provider := NewTemplatesProvider(assetsProvider, fs, "templates/{name}", true)
		view, errs := provider.Template(goattext.DefaultLayout, "mytemplate")
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

func TestHelperLoad(t *testing.T) {
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
	if err := fs.MkdirAll("helpers/", 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("helpers/main.tmpl", []byte(masterTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("helpers/overlay.tmpl", []byte(overlayTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		assetsProvider := NewAssetsProvider(fs, goattext.HelpersPath, goattext.LayoutPath, funcs, true)
		provider := NewTemplatesProvider(assetsProvider, fs, goattext.HelpersPath, true)
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

func TestNoEscape(t *testing.T) {
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
	if err := fs.MkdirAll("helpers/", 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("helpers/default/main.tmpl", []byte(masterTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("templates/mytemplate/main.tmpl", []byte("<-no escape"), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		assetsProvider := NewAssetsProvider(fs, goattext.HelpersPath, goattext.LayoutPath, funcs, true)
		provider := NewTemplatesProvider(assetsProvider, fs, "templates/{name}", true)
		view, errs := provider.Template(goattext.DefaultLayout, "mytemplate")
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
		if result != "<-no escape" {
			t.Errorf("Return value should be unescape %v", result)
			return
		}
	}
}
