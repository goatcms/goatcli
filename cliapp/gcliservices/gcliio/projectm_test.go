package gcliio

import (
	"testing"

	"github.com/goatcms/goatcore/filesystem"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestLoadCorrectIfDefIsEmptyForNonInteractiveModeStory(t *testing.T) {
	t.Parallel()
	var (
		err  error
		mapp app.App
		deps struct {
			ProjectManager gcliservices.GCLIProjectManager `dependency:"GCLIProjectManager"`
		}
	)
	if mapp, err = newMockupApp(mockupapp.MockupOptions{
		Args: []string{"interactive=false"},
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if _, err = deps.ProjectManager.Project(mapp.IOContext()); err != nil {
		t.Error(err)
		return
	}
}

func TestProjectManagerRequireSecretsForNonInteractiveModeStory(t *testing.T) {
	t.Parallel()
	var (
		err  error
		mapp app.App
		deps struct {
			ProjectManager gcliservices.GCLIProjectManager `dependency:"GCLIProjectManager"`
		}
	)
	if mapp, err = newMockupApp(mockupapp.MockupOptions{
		Args: []string{"interactive=false"},
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.HomeFilespace().WriteFile(".goat/secrets.def/00.json", []byte(`[
		{
		  "prompt": "Insert test value",
		  "key": "test",
		  "type": "line",
		  "min": 0,
		  "max": 1024
		}
	  ]`), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if _, err = deps.ProjectManager.Project(mapp.IOContext()); err == nil {
		t.Errorf("expected error for unknow key when non interactive mode")
		return
	}
}
