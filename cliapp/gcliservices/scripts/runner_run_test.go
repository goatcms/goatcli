package scripts

import (
	"strings"
	"testing"
	"time"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func TestPipRunWaitStory(t *testing.T) {
	t.Parallel()
	var (
		err  error
		mapp *mockupapp.App
		deps struct {
			ScriptsRunner gcliservices.ScriptsRunner `dependency:"ScriptsRunner"`
		}
	)
	if mapp, _, err = newApp(mockupapp.MockupOptions{}); err != nil {
		t.Error(err)
		return
	}
	if err = goaterr.ToError(goaterr.AppendError(nil, app.RegisterCommand(mapp, "testCommand", func(a app.App, ctx app.IOContext) (err error) {
		time.Sleep(10 * time.Millisecond)
		return ctx.IO().Out().Printf("test_output")
	}, ""))); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.RootFilespace()
	if err = fs.WriteFile(".goat/scripts/scriptName/main.tmpl", []byte(`testCommand`), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	// test
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	appData := am.NewApplicationData(map[string]string{})
	if err = deps.ScriptsRunner.RunByName(mapp.IOContext(), gcliservices.ScriptsRunnerParams{
		FS:         fs,
		ScriptName: "scriptName",
		Properties: map[string]string{},
		Secrets:    map[string]string{},
		Data:       appData,
	}); err != nil {
		t.Error(err)
		return
	}
	output := mapp.OutputBuffer().String()
	if strings.Index(output, "test_output") == -1 {
		t.Errorf("expected output contains 'test_output' and take: %s", output)
	}
}
