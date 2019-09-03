package builder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/data"
	"github.com/goatcms/goatcli/cliapp/services/dependencies"
	"github.com/goatcms/goatcli/cliapp/services/modules"
	"github.com/goatcms/goatcli/cliapp/services/repositories"
	"github.com/goatcms/goatcli/cliapp/services/template"
	"github.com/goatcms/goatcli/cliapp/services/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/scope"
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
	{{$ctx := .}}
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
		Input:  gio.NewInput(strings.NewReader("")),
		Output: gio.NewOutput(new(bytes.Buffer)),
	}); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.RootFilespace()
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		fs.WriteFile(".goat/build/layouts/default/main.tmpl", []byte(testCTXAMLayout), 0766),
		fs.WriteFile(".goat/build/templates/names/main.tmpl", []byte(testCTXAMTemplate), 0766),
		fs.WriteFile(".goat/data/model/user.json", []byte(testCTXAMData), 0766),
		fs.WriteFile(".goat/build.def.json", []byte(testCTXAMConfig), 0766))); err != nil {
		t.Error(err)
		return
	}
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		RegisterDependencies(mapp.DependencyProvider()),
		modules.RegisterDependencies(mapp.DependencyProvider()),
		dependencies.RegisterDependencies(mapp.DependencyProvider()),
		repositories.RegisterDependencies(mapp.DependencyProvider()),
		template.RegisterDependencies(mapp.DependencyProvider()),
		vcs.RegisterDependencies(mapp.DependencyProvider()),
		data.RegisterDependencies(mapp.DependencyProvider()),
		template.InitDependencies(mapp))); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		BuilderService services.BuilderService `dependency:"BuilderService"`
		DataService    services.DataService    `dependency:"DataService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	ctxScope := scope.NewScope("test")
	if ctxData, err = deps.DataService.ReadDataFromFS(fs); err != nil {
		t.Error(err)
		return
	}
	appModel := am.NewApplicationModel(ctxData)
	buildContext := deps.BuilderService.NewContext(ctxScope, appModel, ctxData, map[string]string{}, map[string]string{})
	if err = buildContext.Build(fs); err != nil {
		t.Error(err)
		return
	}
	if err = ctxScope.Wait(); err != nil {
		t.Error(err)
		return
	}
	if err = ctxScope.Trigger(app.CommitEvent, nil); err != nil {
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