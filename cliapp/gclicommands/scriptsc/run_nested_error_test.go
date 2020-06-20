package scriptsc

import (
	"testing"
	"time"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func TestPipRunNestedErrorStory(t *testing.T) {
	t.Parallel()
	var (
		err         error
		mapp        *mockupapp.App
		bootstraper app.Bootstrap
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
	}, ""), app.RegisterCommand(mapp, "error", func(a app.App, ctx app.IOContext) (err error) {
		time.Sleep(10 * time.Millisecond)
		return goaterr.Errorf("error")
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
pip:run --name=kill --sandbox=self --body=<<EOF
error
EOF`), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	// test
	if err = bootstraper.Run(); err == nil {
		t.Errorf("expected error")
		return
	}
}
