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
	testAssetsProviderMasterTemplate        = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
	testAssetsProviderOverlayAssetsTemplate = `{{define "list"}} {{join . ", "}}{{end}} `
	testAssetsProviderFile1                 = `{{define "list"}} {{join . ", "}}{{end}} `
	testAssetsProviderFile2                 = `{{define "unusedDef1"}} {{join . ": "}}{{end}} `
)

func TestLoadDefaultLayoutByAssetsProvider(t *testing.T) {
	t.Parallel()
	var (
		funcs     = template.FuncMap{"join": strings.Join}
		guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)
	fs, err := memfs.NewFilespace()
	if err != nil {
		t.Error(err)
		return
	}
	// create test data
	if err := fs.MkdirAll("layouts/", 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("layouts/default/main.tmpl", []byte(testAssetsProviderMasterTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		provider := NewAssetsProvider(fs, goattext.HelpersPath, goattext.LayoutPath, funcs, true)
		view, err := provider.Layout(goattext.DefaultLayout)
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
		if !strings.Contains(result, "- Gamora") {
			t.Errorf("after render should contains '- Gamora' and other characters. It is: %v", result)
			return
		}
	}
}

func TestLoadManyFilesByAssetsProvider(t *testing.T) {
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
	if err := fs.WriteFile("layouts/default/main.tmpl", []byte(testAssetsProviderMasterTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("layouts/default/file1.tmpl", []byte(testAssetsProviderFile1), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("layouts/default/file2.tmpl", []byte(testAssetsProviderFile2), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		provider := NewAssetsProvider(fs, goattext.HelpersPath, goattext.LayoutPath, funcs, true)
		view, errs := provider.Layout("default")
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

func TestHelperLoadByAssetsProvider(t *testing.T) {
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
	if err := fs.WriteFile("helpers/main.tmpl", []byte(testAssetsProviderMasterTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("helpers/someother.tmpl", []byte(testAssetsProviderOverlayAssetsTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		provider := NewAssetsProvider(fs, goattext.HelpersPath, goattext.LayoutPath, funcs, true)
		view, errs := provider.Base()
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

func TestNoEscapeByAssetsProvider(t *testing.T) {
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
	if err := fs.WriteFile("helpers/default/main.tmpl", []byte(testAssetsProviderMasterTemplate), 0777); err != nil {
		t.Error(err)
		return
	}
	if err := fs.WriteFile("layouts/default/main.tmpl", []byte("<-no escape"), 0777); err != nil {
		t.Error(err)
		return
	}
	// test loop
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		provider := NewAssetsProvider(fs, goattext.HelpersPath, goattext.LayoutPath, funcs, true)
		view, errs := provider.Layout("default")
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
