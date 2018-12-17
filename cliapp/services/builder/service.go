package builder

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/builder/bcontext"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil"
)

// Service build structure
type Service struct {
	deps struct {
		CWD             string                       `argument:"?cwd"`
		TemplateService services.TemplateService     `dependency:"TemplateService"`
		Modules         services.ModulesService      `dependency:"ModulesService"`
		Dependencies    services.DependenciesService `dependency:"DependenciesService"`
		Repositories    services.RepositoriesService `dependency:"RepositoriesService"`
	}
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
	return services.BuilderService(instance), nil
}

// Build scaffolding a new app and clone dependencies
func (s *Service) Build(fs filesystem.Filespace, buildConfigs []*config.Build, data, properties, secrets map[string]string) (err error) {
	// clone dependencies
	if err = s.CloneDependencies(fs); err != nil {
		return err
	}
	return s.build("", fs, buildConfigs, data, properties, secrets)
}

// Build project files and directories from data
func (s *Service) build(subPath string, fs filesystem.Filespace, buildConfigs []*config.Build, data, properties, secrets map[string]string) (err error) {
	var (
		templateExecutor services.TemplateExecutor
		writer           *FSWriter
		hash             string
	)
	// build modules
	if err = s.BuildModules(subPath, fs, data, properties, secrets); err != nil {
		return err
	}
	// build main code
	hash = varutil.RandString(30, varutil.AlphaNumericBytes)
	if templateExecutor, err = s.deps.TemplateService.Build(fs); err != nil {
		return err
	}
	writer = NewFSWriter(fs, hash)
	for _, c := range buildConfigs {
		context := bcontext.NewBuildContext(&bcontext.Options{
			From: c.From,
			To:   c.To,
			FS:   fs,
			Data: data,
			Hash: hash,
			Properties: bcontext.PropertieOptions{
				Project: properties,
				Secrets: secrets,
				Build:   c.Properties,
			},
		})
		if err = templateExecutor.Execute(c.Layout, c.Template, writer, context); err != nil {
			return err
		}
		if c.AfterBuild != "" {
			var (
				command string
				out     bytes.Buffer
				args    []string
			)
			command = strings.Replace(c.AfterBuild, "\\\"", "\"", -1)
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
		}
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
func (s *Service) BuildModules(subPath string, fs filesystem.Filespace, data, properties, secrets map[string]string) (err error) {
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
		if err = s.build(moduleSubPath, modulefs, buildConfigs, data, properties, secrets); err != nil {
			return err
		}
	}
	return nil
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
