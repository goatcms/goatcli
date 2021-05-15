package data

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/goatapp"
)

const (
	testDataDefJSON = `[{"type":"data_name", "properties":[{"key":"key", "type":"alnum", "min":1, "max":22}]}]`
)

func TestDataDefFromFile(t *testing.T) {
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
	if err = mapp.Filespaces().Root().WriteFile(DataDefPath, []byte(testDataDefJSON), 0766); err != nil {
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
	if len(dataSets) != 1 {
		t.Errorf("expected one module and take %d", len(dataSets))
		return
	}
}

func TestDataDefDefaultEmpty(t *testing.T) {
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
	if len(dataSets) != 0 {
		t.Errorf("should return empty array when def file is not exist (result %v)", dataSets)
		return
	}
}
