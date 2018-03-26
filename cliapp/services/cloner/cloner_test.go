package cloner

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/result"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/modules"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/repositories"
)

func TestModulesFromFile(t *testing.T) {
	var (
		destFS              filesystem.Filespace
		repositoriesService services.RepositoriesService
		propertiesResult    *result.PropertiesResult
		err                 error
		mapp                app.App
		content             []byte
		deps                struct {
			Cloner services.ClonerService `dependency:"ClonerService"`
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
	propertiesResult = buildPropertiesResult()
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
	if err = modules.RegisterDependencies(mapp.DependencyProvider()); err != nil {
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
	if err = deps.Cloner.Clone("https://github.com/goatcms/mockup", repositories.Version{
		Branch: "master",
	}, destFS, propertiesResult); err != nil {
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
	if !destFS.IsFile("module/module.go") {
		t.Errorf("Clone should clone module/module.go")
		return
	}
	if content, err = destFS.ReadFile("docs/main.md"); err != nil {
		t.Error(err)
		return
	}
	if string(content) != "Description your my_project" {
		t.Errorf("should inject project name to {{project_name}} (result is: %s)", content)
		return
	}
}
