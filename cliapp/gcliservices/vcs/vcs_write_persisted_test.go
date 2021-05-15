package vcs

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app/goatapp"
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
	mapp, err := goatapp.NewMockupApp(goatapp.Params{})
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
	fs := mapp.Filespaces().CWD()
	if vcsData, err = deps.VCSService.ReadDataFromFS(fs); err != nil {
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
	if err = deps.VCSService.WriteDataToFS(fs, vcsData); err != nil {
		t.Error(err)
		return
	}
	if data, err = fs.ReadFile(PersistedFilesPath); err != nil {
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
