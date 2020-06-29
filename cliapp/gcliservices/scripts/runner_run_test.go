package scripts

import (
	"strings"
	"testing"
	"time"

	"github.com/goatcms/goatcli/cliapp/common/gclivarutil"

	"github.com/goatcms/goatcli/cliapp/common"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipservices"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipservices/namespaces"
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
		taskManager      pipservices.TasksManager
		appData          gcliservices.ApplicationData
		emptyElasticData common.ElasticData
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
	if appData, err = am.NewApplicationData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if emptyElasticData, err = gclivarutil.NewElasticData(map[string]string{}); err != nil {
		return
	}
	if taskManager, err = deps.ScriptsRunner.Run(gcliservices.ScriptsContext{
		Scope: mapp.IOContext().Scope(),
		CWD:   mapp.IOContext().IO().CWD(),
		Namespaces: namespaces.NewNamespaces(pipservices.NamasepacesParams{
			Task: "",
			Lock: "",
		}),
	}, fs, "scriptName", emptyElasticData, emptyElasticData, appData); err != nil {
		t.Error(err)
		return
	}
	if err = taskManager.Wait(); err != nil {
		t.Error(err)
		return
	}
	// Expect empty context output
	output := mapp.OutputBuffer().String()
	if output != "" {
		t.Errorf("Expected empty output")
	}
	// Output broadcast should contains tasks output
	oString := taskManager.OBroadcast().String()
	if strings.Index(oString, "test_output") == -1 {
		t.Errorf("Output broadcast should contains tasks output. And it is '%s'", oString)
	}
	// StatusBroadcast should contains task summary without task output
	sString := taskManager.StatusBroadcast().String()
	if strings.Index(sString, "[scriptName]") == -1 {
		t.Errorf("StatusBroadcast should contains '[scriptName]' task status. And it is '%s'", sString)
	}
	if strings.Index(sString, "test_output") != -1 {
		t.Errorf("StatusBroadcast should not contains tasks output. And it is '%s'", sString)
	}
}
