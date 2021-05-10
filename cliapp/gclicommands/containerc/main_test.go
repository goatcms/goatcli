package containerc

import (
	"os"
	"testing"

	"github.com/goatcms/goatcli/cliapp/gclicommands"
	"github.com/goatcms/goatcli/cliapp/gcliservices/builder"
	"github.com/goatcms/goatcli/cliapp/gcliservices/cloner"
	"github.com/goatcms/goatcli/cliapp/gcliservices/data"
	"github.com/goatcms/goatcli/cliapp/gcliservices/dependencies"
	"github.com/goatcms/goatcli/cliapp/gcliservices/gclicore"
	"github.com/goatcms/goatcli/cliapp/gcliservices/gcliio"
	"github.com/goatcms/goatcli/cliapp/gcliservices/modules"
	"github.com/goatcms/goatcli/cliapp/gcliservices/properties"
	"github.com/goatcms/goatcli/cliapp/gcliservices/repositories"
	"github.com/goatcms/goatcli/cliapp/gcliservices/secrets"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcli/cliapp/gcliservices/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/bootstrap"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/modules/commonm"
	"github.com/goatcms/goatcore/app/modules/ocm"
	"github.com/goatcms/goatcore/app/modules/terminalm"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func newApp(options mockupapp.MockupOptions) (mapp *mockupapp.App, bootstraper app.Bootstrap, err error) {
	if mapp, err = mockupapp.NewApp(options); err != nil {
		return nil, nil, err
	}
	bootstraper = bootstrap.NewBootstrap(mapp)
	dp := mapp.DependencyProvider()
	if err = goaterr.ToError(goaterr.AppendError(nil,
		bootstraper.Register(terminalm.NewModule()),
		bootstraper.Register(commonm.NewModule()),
		bootstraper.Register(ocm.NewModule()),
		builder.RegisterDependencies(dp),
		gclicore.RegisterDependencies(dp),
		gcliio.RegisterDependencies(dp),
		data.RegisterDependencies(dp),
		properties.RegisterDependencies(dp),
		secrets.RegisterDependencies(dp),
		template.RegisterDependencies(dp),
		modules.RegisterDependencies(dp),
		dependencies.RegisterDependencies(dp),
		repositories.RegisterDependencies(dp),
		cloner.RegisterDependencies(dp),
		vcs.RegisterDependencies(dp),
		app.RegisterCommand(mapp, "container:image", RunContainerImagePip, gclicommands.ContainerImage),
	)); err != nil {
		return nil, nil, err
	}
	if err = bootstraper.Init(); err != nil {
		return nil, nil, err
	}
	return mapp, bootstraper, nil
}

func newBaseFS(t *testing.T, testName string) (fs filesystem.Filespace, err error) {
	var (
		testBasePath = os.Getenv("GOATCORE_TEST_TMPDIR")
		path         = testBasePath + "./teststmp/cliapp/gclicommands/containerc/" + testName
	)
	if err = os.MkdirAll(path, filesystem.DefaultUnixDirMode); err != nil {
		return
	}
	t.Cleanup(func() {
		os.RemoveAll(path)
	})
	return diskfs.NewFilespace(path)
}
