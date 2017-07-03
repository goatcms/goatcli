package data

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestConsoleRead(t *testing.T) {
	var (
		err     error
		mapp    app.App
		data    map[string]string
		dataSet *config.DataSet
	)
	t.Parallel()
	// prepare mockup application
	if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("sdasd\n111\nSebastian\n\nPozoga\na")),
		Output: gio.NewOutput(new(bytes.Buffer)),
	}); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	dataSet = &config.DataSet{
		Name: "UserFixture",
		Properties: []*config.Property{
			&config.Property{
				Key:    "id",
				Min:    1,
				Max:    20,
				Type:   "numeric",
				Prompt: "insert user id",
			},
			&config.Property{
				Key:    "firstname",
				Min:    1,
				Max:    20,
				Type:   "alnum",
				Prompt: "insert user fisrtname",
			},
			&config.Property{
				Key:    "lastname",
				Min:    1,
				Max:    20,
				Type:   "alnum",
				Prompt: "insert user lastname",
			},
		},
	}
	// test
	var deps struct {
		Data services.DataService `dependency:"DataService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if data, err = deps.Data.ConsoleReadData(dataSet); err != nil {
		t.Error(err)
		return
	}
	if len(data) != 3 {
		t.Errorf("expected 3 read elements and take %d", len(data))
		return
	}
}
