package gcliio

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestProjectManagerStory(t *testing.T) {
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
	if _, err = deps.Keystorage.Password("key", "insert value"); err == nil {
		t.Errorf("expected error for unknow key and non interactive mode")
		return
	}
}
