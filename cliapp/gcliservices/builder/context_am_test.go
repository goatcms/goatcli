package builder

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/data"
	"github.com/goatcms/goatcli/cliapp/gcliservices/dependencies"
	"github.com/goatcms/goatcli/cliapp/gcliservices/modules"
	"github.com/goatcms/goatcli/cliapp/gcliservices/repositories"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcli/cliapp/gcliservices/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

const (
	testCTXAMLayout = `{{- define "out/file.txt"}}
		{{- $ctx := .}}
		{{- $entities := $ctx.AM.Entities }}
		{{- index $ctx.PlainData "datakey" }}
		{{- index $ctx.Properties.Project "propkey" }}
		{{- index $ctx.Properties.Secrets "secretkey" }}

		Entities:
		{{- range $index, $entity := $entities }}
			{{$entity.Name.CamelCaseUF}}
		{{- end -}}
	{{- end}}`
	testCTXAMTemplate = `
	{{$ctx.RenderOnce "out/file.txt" "" "" "out/file.txt" $ctx.DotData}}
	`
	testCTXAMConfig = `[{
	  "from":"ignore",
	  "to":"ignore",
	  "template":"names",
	  "layout":"default"
	}]`
	testCTXAMData = `
	{
	  "model": {
	    "user": {
	      "name": "user",
	      "plural": "users",
	      "label":"firstname",
	      "fields": {
	        "0": {
	          "name": "firstname",
	          "system": "n",
	          "type": "string",
	          "unique": "n",
	          "required": "y"
	        },
	      }
	    },
			"session": {
				"name": "session",
				"plural": "sessions",
				"label":"data",
				"fields": {
					"0": {
						"name": "data",
						"system": "n",
						"type": "string",
						"unique": "n",
						"required": "y"
					},
				}
			}
	  }
	}
`
)

func TestCTXBuilderAfterBuild(t *testing.T) {
	var (
		mapp      app.App
		err       error
		fileBytes []byte
		ctxData   map[string]string
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
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		fs.WriteFile(".goat/build/layouts/default/main.tmpl", []byte(testCTXAMLayout), 0766),
		fs.WriteFile(".goat/build/templates/names/main.ctrl", []byte(testCTXAMTemplate), 0766),
		fs.WriteFile(".goat/data/model/user.json", []byte(testCTXAMData), 0766),
		fs.WriteFile(".goat/build.def.json", []byte(testCTXAMConfig), 0766))); err != nil {
		t.Error(err)
		return
	}
	dp := mapp.DependencyProvider()
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		RegisterDependencies(dp),
		modules.RegisterDependencies(dp),
		dependencies.RegisterDependencies(dp),
		repositories.RegisterDependencies(dp),
		template.RegisterDependencies(dp),
		vcs.RegisterDependencies(dp),
		data.RegisterDependencies(dp),
		simpletf.RegisterFunctions(dp))); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		BuilderService gcliservices.BuilderService `dependency:"BuilderService"`
		DataService    gcliservices.DataService    `dependency:"DataService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	ctx := mapp.IOContext()
	if ctxData, err = deps.DataService.ReadDataFromFS(fs); err != nil {
		t.Error(err)
		return
	}
	appData := am.NewApplicationData(ctxData)
	if err = deps.BuilderService.Build(ctx, fs, appData, map[string]string{}, map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if err = ctx.Scope().Wait(); err != nil {
		t.Error(err)
		return
	}
	if err = ctx.Scope().Trigger(app.CommitEvent, nil); err != nil {
		t.Error(err)
		return
	}
	if !fs.IsFile("out/file.txt") {
		t.Errorf("out/file.txt is not exist")
		return
	}
	if fileBytes, err = fs.ReadFile("out/file.txt"); err != nil {
		t.Error(err)
		return
	}
	fileString := string(fileBytes)
	if !strings.Contains(fileString, "User") {
		t.Errorf("Result should contains 'User' and it is '%s'", fileString)
	}
	if !strings.Contains(fileString, "Session") {
		t.Errorf("Result should contains 'User' and it is '%s'", fileString)
	}
}
