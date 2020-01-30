package buildc

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func TestPipRunWaitStory(t *testing.T) {
	t.Parallel()
	var (
		err         error
		mapp        *mockupapp.App
		bootstraper app.Bootstrap
		resultBytes []byte
	)
	if mapp, bootstraper, err = newApp(mockupapp.MockupOptions{
		Args: []string{`appname`, `build`},
	}); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.RootFilespace()
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
	if err = mapp.AppScope().Wait(); err != nil {
		t.Error(err)
		return
	}
	if resultBytes, err = fs.ReadFile("out/file.txt"); err != nil {
		t.Error(err)
		return
	}
	result := string(resultBytes)
	if strings.Index(result, "expected content") == -1 {
		t.Errorf("expected 'test_output' in application output and take: %v", result)
		return
	}
}
