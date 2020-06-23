package scriptsc

import (
	"strings"
	"testing"
	"time"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipservices"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func TestPipRunPromptStory(t *testing.T) {
	t.Parallel()
	var (
		err         error
		mapp        *mockupapp.App
		bootstraper app.Bootstrap
		deps        struct {
			TasksUnit pipservices.TasksUnit `dependency:"PipTasksUnit"`
		}
	)
	if mapp, bootstraper, err = newApp(mockupapp.MockupOptions{
		Args: []string{`appname`, `scripts:run`, `scriptName`},
	}); err != nil {
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
	if err = bootstraper.Run(); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.AppScope().Wait(); err != nil {
		t.Error(err)
		return
	}
	if mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	result := mapp.OutputBuffer().String()
	if strings.Index(result, "\n>\n") != -1 {
		t.Errorf("output should not contains prompt character and it is: %s", result)
		return
	}
}
