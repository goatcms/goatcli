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

func TestPipRunNestedStory(t *testing.T) {
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
		Args: []string{`appname`, `scripts:run`, `first`},
	}); err != nil {
		t.Error(err)
		return
	}
	if err = goaterr.ToError(goaterr.AppendError(nil, app.RegisterCommand(mapp, "echoone", func(a app.App, ctx app.IOContext) (err error) {
		time.Sleep(10 * time.Millisecond)
		return ctx.IO().Out().Printf("1")
	}, ""), app.RegisterCommand(mapp, "echotwo", func(a app.App, ctx app.IOContext) (err error) {
		time.Sleep(10 * time.Millisecond)
		return ctx.IO().Out().Printf("2")
	}, ""))); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.RootFilespace()
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
	if strings.Index(result, "[first:echo]... success") == -1 {
		t.Errorf("output should contains '[first:echo]... success' and it is: %s", result)
		return
	}
	if strings.Index(result, "[first:runsecond:second:echo]... success") == -1 {
		t.Errorf("output should contains '[first:runsecond:second:echo]... success' and it is: %s", result)
		return
	}
}
