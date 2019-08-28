package builder

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/builder/executor"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
)

// Service build structure
type Service struct {
	deps struct {
		CWD             string                       `argument:"?cwd"`
		ExecutorLimit   string                       `argument:"?executor.limit"`
		TemplateService services.TemplateService     `dependency:"TemplateService"`
		Modules         services.ModulesService      `dependency:"ModulesService"`
		Dependencies    services.DependenciesService `dependency:"DependenciesService"`
		Repositories    services.RepositoriesService `dependency:"RepositoriesService"`
	}
	limit int64
}

// ServiceFactory create new repositories instance
func ServiceFactory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &Service{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	if instance.deps.CWD == "" {
		instance.deps.CWD = "./"
	}
	if instance.deps.ExecutorLimit == "" {
		instance.limit = DefaultExecutorLimit
	} else {
		if instance.limit, err = strconv.ParseInt(instance.deps.ExecutorLimit, 10, 64); err != nil {
			return nil, err
		}
	}
	return services.BuilderService(instance), nil
}

// Build scaffolding a new app and clone dependencies
func (s *Service) Build(ctxScope app.Scope, fs filesystem.Filespace, buildConfigs []*config.Build, data, properties, secrets map[string]string) (err error) {
	// clone dependencies
	if err = s.CloneDependencies(fs); err != nil {
		return err
	}
	return s.build(ctxScope, "", fs, buildConfigs, data, properties, secrets)
}

// Build project files and directories from data
func (s *Service) build(ctxScope app.Scope, subPath string, fs filesystem.Filespace, buildConfigs []*config.Build, data, projectProperties, secrets map[string]string) (err error) {
	var (
		templateExecutor  services.TemplateExecutor
		generatorExecutor *executor.GeneratorExecutor
	)
	// build modules
	if err = s.BuildModules(ctxScope, subPath, fs, data, projectProperties, secrets); err != nil {
		return err
	}
	if templateExecutor, err = s.deps.TemplateService.Build(fs); err != nil {
		return err
	}
	if generatorExecutor, err = executor.NewGeneratorExecutor(ctxScope, executor.SharedData{
		PlainData: data,
		Properties: executor.GlobalProperties{
			Project: projectProperties,
			Secrets: secrets,
		},
		FS: fs,
	}, s.limit, templateExecutor); err != nil {
		return err
	}
	for _, c := range buildConfigs {
		generatorExecutor.ExecuteTask(executor.Task{
			Template: executor.TemplateHandler{
				Layout: c.Layout,
				Path:   c.Template,
			},
			DotData: TaskData{
				From: c.From,
				To:   c.To,
			},
			BuildProperties: c.Properties,
			FSPath:          "",
		})
		if c.AfterBuild != "" {
			ctxScope.On(app.CommitEvent, func(d interface{}) error {
				return s.afterBuild(subPath, c.AfterBuild)
			})
		}
	}
	return nil
}

func (s *Service) afterBuild(subPath string, command string) (err error) {
	var (
		out  bytes.Buffer
		args []string
	)
	command = strings.Replace(command, "\\\"", "\"", -1)
	args = strings.Split(command, " ")
	cwd := s.deps.CWD + subPath
	for i := range args {
		// replace it here because argument.cwd can contains space (for example in home directory name)
		args[i] = strings.Replace(args[i], "{{argument.cwd}}", cwd, -1)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("external app fail %v: %v %v", args, err, string(out.Bytes()))
	}
	return nil
}

// CloneDependencies download project dependencies
func (s *Service) CloneDependencies(fs filesystem.Filespace) (err error) {
	var (
		deps []*config.Dependency
	)
	if deps, err = s.deps.Dependencies.ReadDefFromFS(fs); err != nil {
		return err
	}
	return s.deps.Dependencies.CloneDependencies(fs, deps)
}

// BuildModules build project modules
func (s *Service) BuildModules(ctxScope app.Scope, subPath string, fs filesystem.Filespace, data, properties, secrets map[string]string) (err error) {
	var (
		modules []*config.Module
	)
	if modules, err = s.deps.Modules.ReadDefFromFS(fs); err != nil {
		return err
	}
	for _, module := range modules {
		var (
			modulefs     filesystem.Filespace
			buildConfigs []*config.Build
		)
		if !fs.IsExist(module.SourceDir) {
			return fmt.Errorf("builder.BuildModules: Module '%s' is not exist", module.SourceDir)
		}
		if err = fs.MkdirAll(module.SourceDir, 0766); err != nil {
			return err
		}
		if modulefs, err = fs.Filespace(module.SourceDir); err != nil {
			return err
		}
		if buildConfigs, err = s.ReadDefFromFS(modulefs); err != nil {
			return err
		}
		moduleSubPath := subPath + module.SourceDir
		if err = s.build(ctxScope, moduleSubPath, modulefs, buildConfigs, data, properties, secrets); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) buildModule(ctxScope app.Scope, subPath string, module *config.Module, fs filesystem.Filespace, data, properties, secrets map[string]string) (err error) {
	var (
		modulefs     filesystem.Filespace
		buildConfigs []*config.Build
	)
	if !fs.IsExist(module.SourceDir) {
		return fmt.Errorf("builder.BuildModules: Module '%s' is not exist", module.SourceDir)
	}
	if err = fs.MkdirAll(module.SourceDir, 0766); err != nil {
		return err
	}
	if modulefs, err = fs.Filespace(module.SourceDir); err != nil {
		return err
	}
	if buildConfigs, err = s.ReadDefFromFS(modulefs); err != nil {
		return err
	}
	moduleSubPath := subPath + module.SourceDir
	return s.build(ctxScope, moduleSubPath, modulefs, buildConfigs, data, properties, secrets)
}

// ReadDefFromFS return data definition
func (s *Service) ReadDefFromFS(fs filesystem.Filespace) (builds []*config.Build, err error) {
	var json []byte
	if !fs.IsFile(BuildDefPath) {
		return make([]*config.Build, 0), nil
	}
	if json, err = fs.ReadFile(BuildDefPath); err != nil {
		return nil, err
	}
	if builds, err = config.NewBuilds(json); err != nil {
		return nil, err
	}
	return builds, nil
}
