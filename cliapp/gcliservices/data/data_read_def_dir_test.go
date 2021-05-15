package data

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/goatapp"
)

const (
	testDataDefFirstJSON  = `[{"type":"first_type", "properties":[{"key":"key", "type":"alnum", "min":1, "max":22}]}]`
	testDataDefSecondJSON = `[{"type":"second_type", "properties":[{"key":"key", "type":"alnum", "min":1, "max":22}]}]`
)

func TestDataDefFromDirectory(t *testing.T) {
	var (
		err      error
		mapp     app.App
		dataSets []*config.DataSet
	)
	t.Parallel()
	// prepare mockup application
	if mapp, err = goatapp.NewMockupApp(goatapp.Params{}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(".goat/data.def/first.json", []byte(testDataDefFirstJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(".goat/data.def/second.json", []byte(testDataDefSecondJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(".goat/data.def/wrong.ex", []byte("WRONG_FILE"), 0766); err != nil {
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
	if dataSets, err = deps.Data.ReadDefFromFS(mapp.Filespaces().Root()); err != nil {
		t.Error(err)
		return
	}
	if len(dataSets) != 2 {
		t.Errorf("expected two modules and take %d", len(dataSets))
		return
	}
	// check read file
	var (
		hasFirstType  = false
		hasSecondType = false
	)
	// TODO: Check types
	for _, row := range dataSets {
		if row.TypeName == "first_type" {
			hasFirstType = true
		}
		if row.TypeName == "second_type" {
			hasSecondType = true
		}
	}
	if !hasFirstType {
		t.Errorf("File with first type description is not load")
		return
	}
	if !hasSecondType {
		t.Errorf("File with second type description is not load")
		return
	}
}
