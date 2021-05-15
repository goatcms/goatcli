package buildc

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/goatapp"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func TestPipRunWaitStory(t *testing.T) {
	t.Parallel()
	var (
		bootstraper app.Bootstrap
		err         error
		mapp        *goatapp.MockupApp
		resultBytes []byte
	)
	if mapp, bootstraper, err = newApp(goatapp.Params{
		Arguments: []string{`appname`, `build`},
	}); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.Filespaces().CWD()
	if err = goaterr.ToError(goaterr.AppendError(nil,
		fs.WriteFile(".goat/build/layouts/default/main.tmpl", []byte(`{{- define "out/file.txt"}}
			{{- $ctx := .}}
			expected content
		{{- end}}`), 0766), fs.WriteFile(".goat/build/templates/names/main.ctrl", []byte(`
			{{$ctx.RenderOnce "out/file.txt" "" "" "out/file.txt" $ctx.DotData}}
		`), 0766), fs.WriteFile(".goat/build.def.json", []byte(`[{
		  "template":"names",
		  "layout":"default"
		}]`), 0766))); err != nil {
		t.Error(err)
		return
	}
	// test
	if err = bootstraper.Run(); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Scopes().App().Wait(); err != nil {
		t.Error(err)
		return
	}
	if resultBytes, err = fs.ReadFile("out/file.txt"); err != nil {
		t.Error(err)
		return
	}
	result := string(resultBytes)
	if !strings.Contains(result, "expected content") {
		t.Errorf("expected 'test_output' in application output and take: %v", result)
		return
	}
}
