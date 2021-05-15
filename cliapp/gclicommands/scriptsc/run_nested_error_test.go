package scriptsc

import (
	"testing"
	"time"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/goatapp"
	"github.com/goatcms/goatcore/app/terminal"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func TestPipRunNestedErrorStory(t *testing.T) {
	t.Parallel()
	var (
		err         error
		mapp        *goatapp.MockupApp
		bootstraper app.Bootstrap
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
				return goaterr.Errorf("error")
			},
			Name: "error",
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
