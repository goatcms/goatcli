package cloner

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/filesystem"
)

func TestModulesFromFile(t *testing.T) {
	var (
		destFS              filesystem.Filespace
		repositoriesService services.Repositories
		replaces            []*config.Replace
		err                 error
		mapp                app.App
		deps                struct {
			Cloner services.Cloner `dependency:"ClonerService"`
		}
	)
	t.Parallel()
	// prepare data
	if destFS, err = buildDestFilespace(); err != nil {
		t.Error(err)
		return
	}
	if repositoriesService, err = buildRepositoriesService(); err != nil {
		t.Error(err)
		return
	}
	replaces = buildReplaces()
	//propertiesResult = buildPropertiesResult()
	// new app
	if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("  \t\n")),
		Output: gio.NewOutput(new(bytes.Buffer)),
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().SetDefault("RepositoriesService", repositoriesService); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if err = deps.Cloner.Clone("https://github.com/goatcms/mockup", "master", destFS, replaces); err != nil {
		t.Error(err)
		return
	}
	if destFS.IsExist(".git") {
		t.Errorf("Clone should omit .git directory")
		return
	}
	if !destFS.IsFile("main.go") {
		t.Errorf("Clone should clone main.go")
		return
	}
	if !destFS.IsFile("docs/main.md") {
		t.Errorf("Clone should clone docs/main.md")
		return
	}
}
