package vcs

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestVCSWriteGeneratedToFS(t *testing.T) {
	var (
		err          error
		vcsData      services.VCSData
		data         []byte
		dataStr      string
		expectedTime time.Time
	)
	t.Parallel()
	// prepare mockup application & data
	output := new(bytes.Buffer)
	mapp, err := mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("")),
		Output: gio.NewOutput(output),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		VCSService services.VCSService `dependency:"VCSService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if vcsData, err = deps.VCSService.ReadDataFromFS(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	if expectedTime, err = time.Parse(time.RFC3339, "2009-11-10T23:00:00Z"); err != nil {
		t.Error(err)
		return
	}
	vcsData.VCSGeneratedFiles().Add(&services.GeneratedFile{
		Path:    "/first/generated.file",
		ModTime: expectedTime,
	})
	vcsData.VCSGeneratedFiles().Add(&services.GeneratedFile{
		Path:    "/second/generated.file",
		ModTime: expectedTime,
	})
	if err = deps.VCSService.WriteDataToFS(mapp.RootFilespace(), vcsData); err != nil {
		t.Error(err)
		return
	}
	if data, err = mapp.RootFilespace().ReadFile(GeneratedFilesPath); err != nil {
		t.Error(err)
		return
	}
	dataStr = string(data)
	expected := `2009-11-10T23:00:00Z;/first/generated.file
2009-11-10T23:00:00Z;/second/generated.file
`
	if dataStr != expected {
		t.Errorf("take '%s' and expect '%s'", dataStr, expected)
		return
	}
}
