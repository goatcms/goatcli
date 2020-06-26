package properties

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gclimock"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

const (
	testPropetiesDefFirstJSON  = `[{"key":"key1", "type":"alnum", "min":1, "max":22},{"key":"key2", "type":"alnum", "min":1, "max":22}]`
	testPropetiesDefSecondJSON = `[{"key":"key3", "type":"alnum", "min":1, "max":22},{"key":"key4", "type":"alnum", "min":1, "max":22}]`
)

func TestDataDefFromDirectory(t *testing.T) {
	var (
		err        error
		mapp       app.App
		properties []*config.Property
		deps       struct {
			Propeties gcliservices.PropertiesService `dependency:"PropertiesService"`
		}
	)
	t.Parallel()
	// prepare mockup application
	if mapp, err = gclimock.NewApp(mockupapp.MockupOptions{
		Input: gio.NewInput(strings.NewReader("")),
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/properties.def/01_first.json", []byte(testPropetiesDefFirstJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/properties.def/02_second.json", []byte(testPropetiesDefSecondJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	// test
	if properties, err = deps.Propeties.ReadDefFromFS(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	if len(properties) != 4 {
		t.Errorf("expected four properties definitions and take %d", len(properties))
		return
	}
	// check order
	if properties[0].Key != "key1" {
		t.Errorf("properties[0].Key should be equals to key1 and take %v", properties[0].Key)
		return
	}
	if properties[1].Key != "key2" {
		t.Errorf("properties[1].Key should be equals to key1 and take %v", properties[1].Key)
		return
	}
	if properties[2].Key != "key3" {
		t.Errorf("properties[2].Key should be equals to key1 and take %v", properties[2].Key)
		return
	}
	if properties[3].Key != "key4" {
		t.Errorf("properties[3].Keyshould be equals to key1 and take %v", properties[3].Key)
		return
	}
}
