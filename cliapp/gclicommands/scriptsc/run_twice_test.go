package scriptsc

import (
	"strings"
	"testing"
	"time"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func TestPipRunTwiceStory(t *testing.T) {
	t.Parallel()
	var (
		err         error
		mapp        *mockupapp.App
		bootstraper app.Bootstrap
		counter     int
	)
	if mapp, bootstraper, err = newApp(mockupapp.MockupOptions{
		Input: strings.NewReader(`
			scripts:run scriptName
			scripts:run scriptName
		`),
		Args: []string{`appname`, `terminal`},
	}); err != nil {
		t.Error(err)
		return
	}
	if err = goaterr.ToErrors(goaterr.AppendError(nil, app.RegisterCommand(mapp, "testCommand", func(a app.App, ctx app.IOContext) (err error) {
		time.Sleep(10 * time.Millisecond)
		counter++
		return ctx.IO().Out().Printf("test_output%d", counter)
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
	result := mapp.OutputBuffer().String()
	if strings.Index(result, "test_output1") == -1 {
		t.Errorf("expected 'test_output1' in application output and take: %v", result)
		return
	}
	if strings.Index(result, "test_output1") == -1 {
		t.Errorf("expected 'test_output2' in application output and take: %v", result)
		return
	}
}
