package scriptsc

import (
	"os"
	"strings"
	"testing"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func TestPipRunLogsStory(t *testing.T) {
	t.Parallel()
	var (
		err            error
		mapp           *mockupapp.App
		bootstraper    app.Bootstrap
		dirInfo        []os.FileInfo
		ioLogsPath     string
		logsContent    []byte
		summaryPath    string
		summaryContent []byte
	)
	if mapp, bootstraper, err = newApp(mockupapp.MockupOptions{
		Args: []string{`appname`, `scripts:run`, `scriptName`},
	}); err != nil {
		t.Error(err)
		return
	}
	if err = goaterr.ToError(goaterr.AppendError(nil, app.RegisterCommand(mapp, "testCommand", func(a app.App, ctx app.IOContext) (err error) {
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
	cwd := mapp.IOContext().IO().CWD()
	if dirInfo, err = cwd.ReadDir(".goat/tmp/logs/scripts"); err != nil {
		t.Error(err)
		return
	}
	if len(dirInfo) != 2 {
		t.Errorf("expected two log files")
		return
	}
	for _, node := range dirInfo {
		name := node.Name()
		if strings.Contains(name, ".io.log") {
			ioLogsPath = ".goat/tmp/logs/scripts/" + name
		} else if strings.Contains(name, ".summary.log") {
			summaryPath = ".goat/tmp/logs/scripts/" + name
		}
	}
	if ioLogsPath == "" {
		t.Errorf("expected single logs files")
		return
	}
	if summaryPath == "" {
		t.Errorf("expected single summary files")
		return
	}
	// test logs content
	if logsContent, err = cwd.ReadFile(ioLogsPath); err != nil {
		t.Error(err)
		return
	}
	if !strings.Contains(string(logsContent), "test_output") {
		t.Errorf("expected 'test_output' in logs and take '%s'", logsContent)
	}
	// test summary content
	if summaryContent, err = cwd.ReadFile(summaryPath); err != nil {
		t.Error(err)
		return
	}
	if !strings.Contains(string(summaryContent), "test_output") {
		t.Errorf("expected 'test_output' in summary and take '%s'", summaryContent)
	}
}
