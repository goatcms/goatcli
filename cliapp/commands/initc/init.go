package initc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// RunInit run init command
func RunInit(a app.App, ctxScope app.Scope) (err error) {
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
	if err = deps.CurrentFS.MkdirAll(".goat", filesystem.DefaultUnixDirMode); err != nil {
		return err
	}
	if err = deps.CurrentFS.WriteFile(".goat/build.def.json", []byte(InitBildDefJSON), filesystem.DefaultUnixFileMode); err != nil {
		return err
	}
	if err = deps.CurrentFS.WriteFile(".goat/build/helpers/main.tmpl", []byte(InitMainBuildHelper), filesystem.DefaultUnixFileMode); err != nil {
		return err
	}
	if err = deps.CurrentFS.WriteFile(".goat/build/templates/main/main.tmpl", []byte(InitMainTemplateMainTMPL), filesystem.DefaultUnixFileMode); err != nil {
		return err
	}
	if err = deps.CurrentFS.WriteFile(".goat/build/templates/main/readme.tmpl", []byte(InitMainTemplateReadmeTMPL), filesystem.DefaultUnixFileMode); err != nil {
		return err
	}
	if err = deps.CurrentFS.WriteFile(".goat/build/templates/main/gitignore.tmpl", []byte(InitMainTemplateGitignoreTMPL), filesystem.DefaultUnixFileMode); err != nil {
		return err
	}
	if err = deps.CurrentFS.WriteFile(".goat/data.def.json", []byte(InitDataDefJSON), filesystem.DefaultUnixFileMode); err != nil {
		return err
	}
	if err = deps.CurrentFS.WriteFile(".goat/dependencies.def.json", []byte(InitDependenciesDefJSON), filesystem.DefaultUnixFileMode); err != nil {
		return err
	}
	if err = deps.CurrentFS.WriteFile(".goat/properties.def.json", []byte(InitPropertiesDefJSON), filesystem.DefaultUnixFileMode); err != nil {
		return err
	}
	if err = deps.CurrentFS.WriteFile(".goat/replaces.def.json", []byte(InitReplacesDefJSON), filesystem.DefaultUnixFileMode); err != nil {
		return err
	}
	if err = deps.CurrentFS.WriteFile(".goat/secrets.def.json", []byte(InitSecretsDefJSON), filesystem.DefaultUnixFileMode); err != nil {
		return err
	}
	deps.Output.Printf("inited")
	return nil
}
