package vcs

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestVCSReadDataFromFS(t *testing.T) {
	var err error
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
	if err = mapp.RootFilespace().WriteFile(IgnoredFilesPath, []byte(`
		/first/ignored.file
		/second/ignored.file
		`), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(GeneratedFilesPath, []byte(`
		2009-11-10T23:00:00Z;/some/generated_file.go
		2009-11-10T23:00:00Z;/some/other_generated_file.go
		`), 0766); err != nil {
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
	var vcsData services.VCSData
	if vcsData, err = deps.VCSService.ReadDataFromFS(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	if len(vcsData.VCSGeneratedFiles().All()) != 2 {
		t.Errorf("expected two generated file and take %d", len(vcsData.VCSGeneratedFiles().All()))
		return
	}
	if len(vcsData.VCSIgnoredFiles().All()) != 2 {
		t.Errorf("expected two ignored file and take %d", len(vcsData.VCSIgnoredFiles().All()))
		return
	}
}
