package datac

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/prevents"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// RunAdd run data:add command
func RunAdd(a app.App, ctxScope app.Scope) (err error) {
	var (
		deps struct {
			TypeName    string               `command:"?$1"`
			DataName    string               `command:"?$2"`
			CurrentFS   filesystem.Filespace `filespace:"current"`
			DataService services.DataService `dependency:"DataService"`
			Input       app.Input            `dependency:"InputService"`
			Output      app.Output           `dependency:"OutputService"`
		}
		defs    []*config.DataSet
		def     *config.DataSet
		datamap map[string]string
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctxScope.InjectTo(&deps); err != nil {
		return err
	}
	if deps.TypeName == "" {
		return fmt.Errorf("First argument type is required")
	}
	if deps.DataName == "" {
		return fmt.Errorf("Second argument name is required")
	}
	if err = prevents.RequireGoatProject(deps.CurrentFS); err != nil {
		return err
	}
	prefix := deps.TypeName + "." + deps.DataName
	if deps.DataService.HasDataFromFS(deps.CurrentFS, prefix) {
		return fmt.Errorf("Data exists")
	}
	if defs, err = deps.DataService.ReadDefFromFS(deps.CurrentFS); err != nil {
		return err
	}
	for _, d := range defs {
		if d.TypeName == deps.TypeName {
			def = d
			break
		}
	}
	if def == nil {
		return fmt.Errorf("Incorrect type '%s' (add your type to '.goat/data.def.json' file)", deps.TypeName)
	}
	if datamap, err = deps.DataService.ConsoleReadData(def); err != nil {
		return err
	}
	if err = deps.DataService.WriteDataToFS(deps.CurrentFS, prefix, datamap); err != nil {
		return err
	}
	deps.Output.Printf("finish.")
	return nil
}
