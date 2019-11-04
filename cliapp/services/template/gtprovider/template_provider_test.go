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

func TestLoadTemplateWithDefaultLayout(t *testing.T) {
	t.Parallel()
	var (
		baseTemplate *template.Template
		funcs        = template.FuncMap{"join": strings.Join}
		guardians    = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
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
	if err := fs.WriteFile("templates/myview/main.tmpl", []byte(overlayTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		assetsProvider := NewAssetsProvider(fs, goattext.HelpersPath, goattext.LayoutPath, funcs, true)
		if baseTemplate, err = assetsProvider.Layout("default"); err != nil {
			t.Error(err)
			return
		}
		provider := NewTemplateProvider(baseTemplate, fs, "templates/myview", true)
		view, err := provider.Template()
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

func TestLoadTemplateWitchManyFiles(t *testing.T) {
	t.Parallel()
	var (
		baseTemplate *template.Template
		funcs        = template.FuncMap{"join": strings.Join}
		guardians    = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
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
	if err := fs.WriteFile("templates/myview/file1.tmpl", []byte(templateFile1), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("templates/myview/file2.tmpl", []byte(templateFile2), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("templates/myview/file3.tmpl", []byte(templateFile3), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("templates/myview/file4.tmpl", []byte(templateFile4), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		assetsProvider := NewAssetsProvider(fs, goattext.HelpersPath, goattext.LayoutPath, funcs, true)
		if baseTemplate, err = assetsProvider.Layout("default"); err != nil {
			t.Error(err)
			return
		}
		provider := NewTemplateProvider(baseTemplate, fs, "templates/myview", true)
		view, errs := provider.Template()
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

func TestTemplateNoEscape(t *testing.T) {
	t.Parallel()
	var (
		baseTemplate *template.Template
		funcs        = template.FuncMap{"join": strings.Join}
		guardians    = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
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
	if err := fs.WriteFile("templates/myview/main.tmpl", []byte("<-no escape"), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		assetsProvider := NewAssetsProvider(fs, goattext.HelpersPath, goattext.LayoutPath, funcs, true)
		if baseTemplate, err = assetsProvider.Layout("default"); err != nil {
			t.Error(err)
			return
		}
		provider := NewTemplateProvider(baseTemplate, fs, "templates/myview", true)
		view, errs := provider.Template()
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
