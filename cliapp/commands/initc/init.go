package initc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
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
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		deps.CurrentFS.MkdirAll(".goat", filesystem.DefaultUnixDirMode),
		deps.CurrentFS.WriteFile(".goat/build.def.json", []byte(InitBildDefJSON), filesystem.DefaultUnixFileMode),
		deps.CurrentFS.WriteFile(".goat/build/helpers/main.tmpl", []byte(InitMainBuildHelper), filesystem.DefaultUnixFileMode),
		deps.CurrentFS.WriteFile(".goat/build/templates/main/main.tmpl", []byte(InitMainTemplateMainTMPL), filesystem.DefaultUnixFileMode),
		deps.CurrentFS.WriteFile(".goat/build/templates/main/readme.tmpl", []byte(InitMainTemplateReadmeTMPL), filesystem.DefaultUnixFileMode),
		deps.CurrentFS.WriteFile(".goat/build/templates/main/gitignore.tmpl", []byte(InitMainTemplateGitignoreTMPL), filesystem.DefaultUnixFileMode),
		deps.CurrentFS.WriteFile(".goat/data.def.json", []byte(InitDataDefJSON), filesystem.DefaultUnixFileMode),
		deps.CurrentFS.WriteFile(".goat/dependencies.def.json", []byte(InitDependenciesDefJSON), filesystem.DefaultUnixFileMode),
		deps.CurrentFS.WriteFile(".goat/properties.def.json", []byte(InitPropertiesDefJSON), filesystem.DefaultUnixFileMode),
		deps.CurrentFS.WriteFile(".goat/replaces.def.json", []byte(InitReplacesDefJSON), filesystem.DefaultUnixFileMode),
		deps.CurrentFS.WriteFile(".goat/secrets.def.json", []byte(InitSecretsDefJSON), filesystem.DefaultUnixFileMode),
	)); err != nil {
		return err
	}
	return ctx.IO().Out().Printf("inited")
}
