package scriptsc

import (
	"strings"
	"testing"
	"time"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/goatapp"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipservices"
	"github.com/goatcms/goatcore/app/terminal"
	"github.com/goatcms/goatcore/filesystem"
)

func TestPipRunPromptStory(t *testing.T) {
	t.Parallel()
	var (
		err         error
		mapp        *goatapp.MockupApp
		bootstraper app.Bootstrap
		deps        struct {
			TasksUnit pipservices.TasksUnit `dependency:"PipTasksUnit"`
		}
	)
	if mapp, bootstraper, err = newApp(goatapp.Params{
		Arguments: []string{`appname`, `scripts:run`, `scriptName`},
	}); err != nil {
		t.Error(err)
		return
	}
	mapp.Terminal().SetCommand(
		terminal.NewCommand(terminal.CommandParams{
			Callback: func(a app.App, ctx app.IOContext) (err error) {
				time.Sleep(10 * time.Millisecond)
				return ctx.IO().Out().Printf("test_output")
			},
			Name: "testCommand",
		}),
	)
	fs := mapp.Filespaces().CWD()
	if err = fs.WriteFile(".goat/scripts/scriptName/main.tmpl", []byte(`testCommand`), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	// test
	if mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if err = bootstraper.Run(); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Scopes().App().Wait(); err != nil {
		t.Error(err)
		return
	}
	result := mapp.OutputBuffer().String()
	if strings.Contains(result, "\n>\n") {
		t.Errorf("output should not contains prompt character and it is: %s", result)
		return
	}
}
