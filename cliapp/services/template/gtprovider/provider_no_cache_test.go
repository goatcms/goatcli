package gtprovider

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/memfs"
	"github.com/goatcms/goatcore/goattext"
	"github.com/goatcms/goatcore/workers"
)

func TestMainNoCacheChanges(t *testing.T) {
	t.Parallel()
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		var (
			fs   filesystem.Filespace
			err  error
			view *template.Template
		)
		if fs, err = memfs.NewFilespace(); err != nil {
			t.Error(err)
			return
		}
		// create test data
		if err = fs.WriteFile("views/test/main.tmpl", []byte("42"), 0777); err != nil {
			t.Error(err)
			return
		}
		provider := NewProvider(fs, goattext.HelpersPath, goattext.LayoutPath, goattext.ViewPath, nil, false)
		if view, err = provider.View(goattext.DefaultLayout, "test"); err != nil {
			t.Errorf("Errors: %v", err)
			return
		}
		buf := new(bytes.Buffer)
		if err = view.Execute(buf, nil); err != nil {
			t.Error(err)
			return
		}
		result := buf.String()
		if result != "42" {
			t.Errorf("Before changes Result should be equals to 42 and it is %s", result)
			return
		}
		if err = fs.WriteFile("views/test/main.tmpl", []byte("2018"), 0777); err != nil {
			t.Error(err)
			return
		}
		if view, err = provider.View(goattext.DefaultLayout, "test"); err != nil {
			t.Errorf("Errors: %v", err)
			return
		}
		buf = new(bytes.Buffer)
		if err = view.Execute(buf, nil); err != nil {
			t.Error(err)
			return
		}
		result = buf.String()
		if result != "2018" {
			t.Errorf("After changes result should be equals to 2018 and it is %s", result)
			return
		}
	}
}
