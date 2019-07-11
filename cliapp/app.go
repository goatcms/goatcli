package cliapp

import (
	"os"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/goatapp"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/app/scope/argscope"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/dependency/provider"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
)

// CLIApp represent command line application
type CLIApp struct {
	name    string
	version string

	rootFilespace filesystem.Filespace

	engineScope     app.Scope
	argsScope       app.Scope
	filespaceScope  app.Scope
	configScope     app.Scope
	dependencyScope app.Scope
	appScope        app.Scope
	commandScope    app.Scope

	dp dependency.Provider
}

// NewCLIApp create new app instance
func NewCLIApp(name, version string) (app.App, error) {
	capp := &CLIApp{
		name:    name,
		version: version,
	}
	if err := capp.initEngineScope(); err != nil {
		return nil, err
	}
	if err := capp.initArgsScope(); err != nil {
		return nil, err
	}
	if err := capp.initFilespaceScope(); err != nil {
		return nil, err
	}
	if err := capp.initConfigScope(); err != nil {
		return nil, err
	}
	if err := capp.initDependencyScope(); err != nil {
		return nil, err
	}
	if err := capp.initAppScope(); err != nil {
		return nil, err
	}
	if err := capp.initCommandScope(); err != nil {
		return nil, err
	}
	capp.dp.SetDefault(app.EngineScope, capp.engineScope)
	capp.dp.SetDefault(app.ArgsScope, capp.argsScope)
	capp.dp.SetDefault(app.FilespaceScope, capp.filespaceScope)
	capp.dp.SetDefault(app.ConfigScope, capp.configScope)
	capp.dp.SetDefault(app.DependencyScope, capp.dependencyScope)
	capp.dp.SetDefault(app.AppScope, capp.appScope)
	capp.dp.SetDefault(app.CommandScope, capp.commandScope)
	capp.dp.SetDefault(app.InputService, gio.NewAppInput(os.Stdin))
	capp.dp.SetDefault(app.OutputService, gio.NewAppOutput(os.Stdout))
	capp.dp.AddInjectors([]dependency.Injector{
		capp.commandScope,
		capp.appScope,
		// capp.dependencyScope, <- it is wraper for dependency injection and musn't
		// contains recursive injection
		capp.configScope,
		capp.filespaceScope,
		capp.argsScope,
		capp.engineScope,
	})

	capp.dp.SetDefault(app.AppService, app.App(capp))
	return capp, nil
}

func (capp *CLIApp) initEngineScope() error {
	capp.engineScope = scope.NewScope(app.EngineTagName)
	capp.engineScope.Set(app.GoatVersion, app.GoatVersionValue)
	return nil
}

func (capp *CLIApp) initArgsScope() (err error) {
	capp.argsScope, err = argscope.NewScope(os.Args, app.ArgsTagName)
	return err
}

func (capp *CLIApp) initFilespaceScope() (err error) {
	var (
		currentfs filesystem.Filespace
		deps      struct {
			CWD string `argument:"?cwd"`
		}
	)
	if capp.rootFilespace, err = diskfs.NewFilespace(RootPath); err != nil {
		return err
	}
	if err = capp.argsScope.InjectTo(&deps); err != nil {
		return err
	}
	if deps.CWD == "" {
		deps.CWD = CurrentPath
	}
	if currentfs, err = diskfs.NewFilespace(deps.CWD); err != nil {
		return err
	}
	capp.filespaceScope = scope.NewScope(app.FilespaceTagName)
	capp.filespaceScope.Set(app.RootFilespace, capp.rootFilespace)
	capp.filespaceScope.Set(CurrentFilesystem, currentfs)
	return nil
}

func (capp *CLIApp) initConfigScope() error {
	capp.configScope = scope.NewScope(app.CommandTagName)
	return nil
}

func (capp *CLIApp) initCommandScope() error {
	capp.commandScope = scope.NewScope(app.CommandTagName)
	return nil
}

func (capp *CLIApp) initDependencyScope() error {
	capp.dp = provider.NewProvider(app.DependencyTagName)
	capp.dependencyScope = goatapp.NewDependencyScope(capp.dp)
	return nil
}

func (capp *CLIApp) initAppScope() error {
	capp.appScope = scope.NewScope(app.AppTagName)
	capp.appScope.Set(app.AppName, capp.name)
	capp.appScope.Set(app.AppVersion, capp.version)
	return nil
}

// Name return app name
func (capp *CLIApp) Name() string {
	return capp.name
}

// Version return app version
func (capp *CLIApp) Version() string {
	return capp.version
}

// EngineScope return engine scope
func (capp *CLIApp) EngineScope() app.Scope {
	return capp.engineScope
}

// ArgsScope return app scope
func (capp *CLIApp) ArgsScope() app.Scope {
	return capp.argsScope
}

// FilespaceScope return filespace scope
func (capp *CLIApp) FilespaceScope() app.Scope {
	return capp.filespaceScope
}

// ConfigScope return config scope
func (capp *CLIApp) ConfigScope() app.Scope {
	return capp.configScope
}

// DependencyScope return dependency scope
func (capp *CLIApp) DependencyScope() app.Scope {
	return capp.dependencyScope
}

// AppScope return app scope
func (capp *CLIApp) AppScope() app.Scope {
	return capp.appScope
}

// CommandScope return command scope
func (capp *CLIApp) CommandScope() app.Scope {
	return capp.commandScope
}

// DependencyProvider return dependency provider
func (capp *CLIApp) DependencyProvider() dependency.Provider {
	return capp.dp
}

// RootFilespace return main filespace for application (current directory by default)
func (capp *CLIApp) RootFilespace() filesystem.Filespace {
	return capp.rootFilespace
}
