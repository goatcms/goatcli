package data

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/goatapp"
)

const (
	testDataFile1JSON = `{"lang.pl":{"application":"moja aplikacja", "welcome":"Witajcie"}}`
	testDataFile2JSON = `{"lang.en":{"application":"my application", "welcome":"Welcome"}}`
	testDataFile3JSON = `{"page.home":{"title":"welcome", "content":"My homepage text"}}`
)

func TestDataFromFile(t *testing.T) {
	var (
		err  error
		mapp app.App
		data map[string]string
	)
	t.Parallel()
	// prepare mockup application
	if mapp, err = goatapp.NewMockupApp(goatapp.Params{}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(".goat/data/langs/pl.json", []byte(testDataFile1JSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(".goat/data/langs/en.json", []byte(testDataFile2JSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(".goat/data/pages/home.json", []byte(testDataFile3JSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		Data gcliservices.DataService `dependency:"DataService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if data, err = deps.Data.ReadDataFromFS(mapp.Filespaces().Root()); err != nil {
		t.Error(err)
		return
	}
	if data["lang.pl.application"] != "moja aplikacja" {
		t.Errorf("lang.pl.application expected 'moja aplikacja' and take %v", data["lang.pl.application"])
		return
	}
	if data["lang.en.application"] != "my application" {
		t.Errorf("lang.pl.application expected 'moja aplikacja' and take %v", data["lang.pl.application"])
		return
	}
}

func TestEmptyData(t *testing.T) {
	var (
		err  error
		mapp app.App
		data map[string]string
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
		Data gcliservices.DataService `dependency:"DataService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if data, err = deps.Data.ReadDataFromFS(mapp.Filespaces().Root()); err != nil {
		t.Error(err)
		return
	}
	if len(data) != 0 {
		t.Errorf("if .goat/data directory is empty or dosn't exist, ReadData")
		return
	}
}
