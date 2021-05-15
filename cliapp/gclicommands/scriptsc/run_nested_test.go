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

func TestPipRunNestedStory(t *testing.T) {
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
		Arguments: []string{`appname`, `scripts:run`, `first`},
	}); err != nil {
		t.Error(err)
		return
	}
	mapp.Terminal().SetCommand(
		terminal.NewCommand(terminal.CommandParams{
			Callback: func(a app.App, ctx app.IOContext) (err error) {
				time.Sleep(10 * time.Millisecond)
				return ctx.IO().Out().Printf("1")
			},
			Name: "echoone",
		}),
		terminal.NewCommand(terminal.CommandParams{
			Callback: func(a app.App, ctx app.IOContext) (err error) {
				time.Sleep(10 * time.Millisecond)
				return ctx.IO().Out().Printf("2")
			},
			Name: "echotwo",
		}),
	)
	fs := mapp.Filespaces().CWD()
	if err = fs.WriteFile(".goat/scripts/first/main.tmpl", []byte(`
{{- $ctx := . }}
pip:run --name=echo --sandbox=self --body=<<EOF
echoone
EOF
pip:run --name=runsecond --wait=echo --sandbox=self --body=<<EOF
scripts:run second
EOF
`), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile(".goat/scripts/second/main.tmpl", []byte(`
{{- $ctx := . }}
pip:run --name=echo --sandbox=self --body=<<EOF
echotwo
EOF`), filesystem.DefaultUnixFileMode); err != nil {
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
	if !strings.Contains(result, "[first:echo]... success") {
		t.Errorf("output should contains '[first:echo]... success' and it is: %s", result)
		return
	}
	if !strings.Contains(result, "[first:runsecond:second:echo]... success") {
		t.Errorf("output should contains '[first:runsecond:second:echo]... success' and it is: %s", result)
		return
	}
}
