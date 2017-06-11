package properties

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/varutil/plainmap"
)

const (
	testPropJSON     = `[{"key":"my.config.key", "type":"alnum", "min":1, "max":22}]`
	testPropDataJSON = `{"my.config.key":"data_from_file"}`
)

func TestProperties(t *testing.T) {
	var (
		value string
		err   error
	)
	t.Parallel()
	output := new(bytes.Buffer)
	mapp, err := mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("my_insert_value\n")),
		Output: gio.NewOutput(output),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/properties.def.json", []byte(testPropJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	var deps struct {
		Properties services.Properties `dependency:"PropertiesService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if value, err = deps.Properties.Get("my.config.key"); err != nil {
		t.Error(err)
		return
	}
	if value != "my_insert_value" {
		t.Errorf("property value for my.config.key should be equals to 'my_insert_value' and it is %s", value)
		return
	}
}

func TestGenerateProperty(t *testing.T) {
	var (
		value string
		err   error
	)
	t.Parallel()
	output := new(bytes.Buffer)
	mapp, err := mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("  \t\n")),
		Output: gio.NewOutput(output),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/properties.def.json", []byte(testPropJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	var deps struct {
		Properties services.Properties `dependency:"PropertiesService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if value, err = deps.Properties.Get("my.config.key"); err != nil {
		t.Error(err)
		return
	}
	if len(value) != 22 {
		t.Errorf("should generate property value (lenght should be equals to max lenght 22 and it is %d)", len(value))
		return
	}
	var dataJSON []byte
	var saved map[string]string
	if dataJSON, err = mapp.RootFilespace().ReadFile(".goat/properties.json"); err != nil {
		t.Error(err)
		return
	}
	if saved, err = plainmap.JSONToPlainStringMap(dataJSON); err != nil {
		t.Error(err)
		return
	}
	if value != saved["my.config.key"] {
		t.Errorf("saved value incorrect (%s != %s)", value, saved["my.config.key"])
		return
	}
}

func TestPropertiesReadFromFile(t *testing.T) {
	var (
		value string
		err   error
	)
	t.Parallel()
	output := new(bytes.Buffer)
	mapp, err := mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("  \t\n")),
		Output: gio.NewOutput(output),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/properties.def.json", []byte(testPropJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/properties.json", []byte(testPropDataJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	var deps struct {
		Properties services.Properties `dependency:"PropertiesService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if value, err = deps.Properties.Get("my.config.key"); err != nil {
		t.Error(err)
		return
	}
	if value != "data_from_file" {
		t.Errorf("should generate property value (lenght should be equals to max lenght 22 and it is %d)", len(value))
		return
	}
}
