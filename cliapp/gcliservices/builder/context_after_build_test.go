package builder

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/gclivarutil"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/dependencies"
	"github.com/goatcms/goatcli/cliapp/gcliservices/modules"
	"github.com/goatcms/goatcli/cliapp/gcliservices/repositories"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcli/cliapp/gcliservices/vcs"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

const (
	testCTXBuilderAfterBuildTemplate = `
	{{$ctx := .}}
	`
	testCTXBuilderAfterBuildConfig = `[{
	  "from":"ignore",
	  "to":"ignore",
	  "template":"names",
	  "layout":"default",
	  "afterBuild": "echo \"TestOK\""
	}]`
)

func TestCTXBuilder(t *testing.T) {
	var (
		mapp             *mockupapp.App
		appData          gcliservices.ApplicationData
		emptyElasticData common.ElasticData
		err              error
	)
	t.Parallel()
	// prepare mockup application & data
	if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
		Input: gio.NewInput(strings.NewReader("")),
	}); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.RootFilespace()
	if err = goaterr.ToError(goaterr.AppendError(nil,
		fs.WriteFile(".goat/build/templates/names/main.tmpl", []byte(testCTXBuilderAfterBuildTemplate), 0766),
		fs.WriteFile(".goat/build.def.json", []byte(testCTXBuilderAfterBuildConfig), 0766))); err != nil {
		t.Error(err)
		return
	}
	dp := mapp.DependencyProvider()
	if err = goaterr.ToError(goaterr.AppendError(nil,
		RegisterDependencies(dp),
		modules.RegisterDependencies(dp),
		dependencies.RegisterDependencies(dp),
		repositories.RegisterDependencies(dp),
		template.RegisterDependencies(dp),
		vcs.RegisterDependencies(dp),
		simpletf.RegisterFunctions(dp))); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		BuilderService gcliservices.BuilderService `dependency:"BuilderService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	ctx := mapp.IOContext()
	if emptyElasticData, err = gclivarutil.NewElasticData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if appData, err = am.NewApplicationData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if err = deps.BuilderService.Build(ctx, fs, appData, emptyElasticData, emptyElasticData); err != nil {
		t.Error(err)
		return
	}
	if err = ctx.Scope().Wait(); err != nil {
		t.Error(err)
		return
	}
	if err = ctx.Scope().Trigger(gcliservices.BuildCommitevent, nil); err != nil {
		t.Error(err)
		return
	}
	outBuffer := mapp.OutputBuffer()
	if !strings.Contains(outBuffer.String(), "TestOK") {
		t.Errorf("expected TestOK in afterBuild command output and take '%s'", outBuffer.String())
		return
	}
}
