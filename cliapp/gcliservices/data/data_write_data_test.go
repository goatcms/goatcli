package data

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/goatapp"
	"github.com/goatcms/goatcore/varutil/plainmap"
)

func TestWriteDataToFS(t *testing.T) {
	var (
		err        error
		mapp       app.App
		filedata   []byte
		resultdata map[string]string
	)
	t.Parallel()
	// prepare mockup application
	if mapp, err = goatapp.NewMockupApp(goatapp.Params{}); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		DataService gcliservices.DataService `dependency:"DataService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	datamap := map[string]string{
		"title":   "Tytuł strony",
		"content": "Treść strony",
	}
	if err = deps.DataService.WriteDataToFS(mapp.Filespaces().Root(), "pages.home", datamap); err != nil {
		t.Error(err)
		return
	}
	destPath := ".goat/data/pages/home.json"
	if !mapp.Filespaces().Root().IsExist(destPath) {
		t.Errorf("expected %s file", destPath)
		return
	}
	if filedata, err = mapp.Filespaces().Root().ReadFile(destPath); err != nil {
		t.Error(err)
		return
	}
	if resultdata, err = plainmap.JSONToPlainStringMap(filedata); err != nil {
		t.Error(err)
		return
	}
	if resultdata["pages.home.title"] != "Tytuł strony" {
		t.Errorf("expected pages.home.title equals to 'Tytuł strony' and it is '%s'", resultdata["pages.home.title"])
		return
	}
	if resultdata["pages.home.content"] != "Treść strony" {
		t.Errorf("expected pages.home.content equals to 'Treść strony' and it is '%s'", resultdata["pages.home.content"])
		return
	}
}
