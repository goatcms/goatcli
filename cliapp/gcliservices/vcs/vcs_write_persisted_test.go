package vcs

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/gclimock"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestVCSWritePersistedToFS(t *testing.T) {
	var (
		err     error
		vcsData gcliservices.VCSData
		data    []byte
		dataStr string
	)
	t.Parallel()
	// prepare mockup application & data
	mapp, err := gclimock.NewApp(mockupapp.MockupOptions{
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
		VCSService gcliservices.VCSService `dependency:"VCSService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if vcsData, err = deps.VCSService.ReadDataFromFS(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	if vcsData.VCSPersistedFiles().Modified() != false {
		t.Errorf("Persisted files should be unmodified after read")
		return
	}
	vcsData.VCSPersistedFiles().AddPath("/first/persisted.file")
	vcsData.VCSPersistedFiles().AddPath("/second/persisted.file")
	if vcsData.VCSPersistedFiles().Modified() != true {
		t.Errorf("after add paths to persisted files Modified flag must be true")
		return
	}
	if err = deps.VCSService.WriteDataToFS(mapp.RootFilespace(), vcsData); err != nil {
		t.Error(err)
		return
	}
	if data, err = mapp.RootFilespace().ReadFile(PersistedFilesPath); err != nil {
		t.Error(err)
		return
	}
	dataStr = string(data)
	expected := `/first/persisted.file
/second/persisted.file
`
	if dataStr != expected {
		t.Errorf("take '%s' and expect '%s'", dataStr, expected)
		return
	}
}
