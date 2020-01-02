package vcs

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestVCSWriteIgnoredToFS(t *testing.T) {
	var (
		err     error
		vcsData services.VCSData
		data    []byte
		dataStr string
	)
	t.Parallel()
	// prepare mockup application & data
	mapp, err := mockupapp.NewApp(mockupapp.MockupOptions{
		Input: gio.NewInput(strings.NewReader("")),
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
	if vcsData.VCSIgnoredFiles().Modified() != false {
		t.Errorf("Ignored files should be unmodified after read")
		return
	}
	vcsData.VCSIgnoredFiles().AddPath("/first/ignored.file")
	vcsData.VCSIgnoredFiles().AddPath("/second/ignored.file")
	if vcsData.VCSIgnoredFiles().Modified() != true {
		t.Errorf("after add paths to ignored files Modified flag must be true")
		return
	}
	if err = deps.VCSService.WriteDataToFS(mapp.RootFilespace(), vcsData); err != nil {
		t.Error(err)
		return
	}
	if data, err = mapp.RootFilespace().ReadFile(IgnoredFilesPath); err != nil {
		t.Error(err)
		return
	}
	dataStr = string(data)
	expected := `/first/ignored.file
/second/ignored.file
`
	if dataStr != expected {
		t.Errorf("take '%s' and expect '%s'", dataStr, expected)
		return
	}
}
