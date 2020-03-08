package initc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// RunInit run init command
func RunInit(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			CurrentFS filesystem.Filespace `filespace:"current"`
		}
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if deps.CurrentFS.IsDir(".goat") {
		return ctx.IO().Err().Printf("already inited")
	}
	for path, fileContent := range goatStructureSnapshot {
		if err = deps.CurrentFS.WriteFile(path, []byte(fileContent), filesystem.DefaultUnixFileMode); err != nil {
			return err
		}
	}
	return ctx.IO().Out().Printf("inited")
}
