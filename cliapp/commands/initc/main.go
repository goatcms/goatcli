package initc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// Run run init command
func Run(a app.App) (err error) {
	var (
		deps struct {
			CurrentFS filesystem.Filespace `filespace:"current"`
			Output    app.Output           `dependency:"OutputService"`
		}
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if deps.CurrentFS.IsDir(".goat") {
		deps.Output.Printf("already inited")
		return nil
	}
	if err = deps.CurrentFS.MkdirAll(".goat", 0766); err != nil {
		return err
	}
	deps.Output.Printf("inited")
	return nil
}
