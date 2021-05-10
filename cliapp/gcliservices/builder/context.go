package builder

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/builder/executor"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/fscache"
	"github.com/goatcms/goatcore/varutil"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// Command is command to run
type Command struct {
	Path string
	CMD  string
}

// Context contains build process data
type Context struct {
	ctx                app.IOContext
	appModel           interface{}
	data               common.ElasticData
	properties         common.ElasticData
	secrets            common.ElasticData
	service            *Service
	vcsData            gcliservices.VCSData
	commitCommands     []Command
	remoteFS           filesystem.Filespace
	cache              *fscache.Cache
	generatorExecutors []*executor.GeneratorExecutor
}

// Build scaffolding a new app and clone dependencies
func (c *Context) Build(fs filesystem.Filespace) (err error) {
	var (
		deps         []*config.Dependency
		buildConfigs []*config.Build
	)
	// create cache
	c.remoteFS = fs
	if c.cache, err = fscache.NewMemCache(fs); err != nil {
		return err
	}
	// clone dependencies
	depService := c.service.deps.Dependencies
	if deps, err = depService.ReadDefFromFS(c.cache); err != nil {
		return err
	}
	if err = depService.CloneDependencies(c.cache, deps); err != nil {
		return err
	}
	// read build config
	if buildConfigs, err = c.service.ReadDefFromFS(c.cache); err != nil {
		return err
	}
	// read vcs data
	if c.vcsData, err = c.service.deps.VCSService.ReadDataFromFS(c.cache); err != nil {
		return err
	}
	// bind commit event
	c.ctx.Scope().On(gcliservices.BuildCommitevent, c.commit)
	// build
	if err = c.build(c.cache, "", buildConfigs); err != nil {
		return err
	}
	if err = c.ctx.Scope().Wait(); err != nil {
		return err
	}
	for _, executor := range c.generatorExecutors {
		if err = executor.ExecuteHook("vcs", c.vcsData); err != nil {
			return err
		}
	}
	return nil
}

// Build project files and directories from data
func (c *Context) build(fs filesystem.Filespace, subPath string, buildConfigs []*config.Build) (err error) {
	var (
		templatesExecutor gcliservices.TemplatesExecutor
		generatorExecutor *executor.GeneratorExecutor
	)
	// build modules
	if err = c.buildModules(fs, subPath); err != nil {
		return err
	}
	if templatesExecutor, err = c.service.deps.TemplateService.TemplatesExecutor(); err != nil {
		return err
	}
	if generatorExecutor, err = executor.NewGeneratorExecutor(c.ctx.Scope(), executor.SharedData{
		AM:   c.appModel,
		Data: c.data,
		Properties: executor.GlobalProperties{
			Project: c.properties,
			Secrets: c.secrets,
		},
		FS:      fs,
		VCSData: c.vcsData,
	}, c.service.limit, templatesExecutor); err != nil {
		return err
	}
	c.generatorExecutors = append(c.generatorExecutors, generatorExecutor)
	for _, config := range buildConfigs {
		if config.AfterBuild != "" {
			c.commitCommands = append(c.commitCommands, Command{
				Path: subPath,
				CMD:  config.AfterBuild,
			})
		}
		if err = generatorExecutor.ExecuteView(config.Layout, config.Template, config.Properties, TaskData{
			From: config.From,
			To:   config.To,
		}); err != nil {
			return err
		}
	}
	return nil
}

// commit build save data to filespace, run after build commands , and update generated files modtime
func (c *Context) commit(data interface{}) (err error) {
	// persist files
	if err = c.cache.Commit(); err != nil {
		return err
	}
	// run after build commands
	for _, command := range c.commitCommands {
		if err = c.runCommand(command); err != nil {
			return err
		}
	}
	// update and persist vcs data
	generatedFiles := c.vcsData.VCSGeneratedFiles()
	for _, file := range generatedFiles.New() {
		var info os.FileInfo
		if info, err = c.remoteFS.Lstat(file.Path); err != nil {
			return err
		}
		generatedFiles.Add(&gcliservices.GeneratedFile{
			Path:    file.Path,
			ModTime: info.ModTime(),
		})
	}
	vcsService := c.service.deps.VCSService
	return vcsService.WriteDataToFS(c.remoteFS, c.vcsData)
}

func (c *Context) runCommand(command Command) (err error) {
	var (
		out  bytes.Buffer
		args []string
	)
	if args, _, err = varutil.SplitArguments(command.CMD); err != nil {
		return err
	}
	cwd := c.service.deps.CWD + command.Path
	for i := range args {
		// replace it here because argument.cwd can contains space (for example in home directory name)
		args[i] = strings.Replace(args[i], "{{argument.cwd}}", cwd, -1)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err = cmd.Run(); err != nil {
		return goaterr.Errorf("external app fail %v: %v %v", args, err, string(out.Bytes()))
	}
	c.ctx.IO().Out().Printf("%s", out.Bytes())
	return nil
}

// buildModules build project modules
func (c *Context) buildModules(fs filesystem.Filespace, subPath string) (err error) {
	var (
		modules []*config.Module
	)
	if modules, err = c.service.deps.Modules.ReadDefFromFS(fs); err != nil {
		return err
	}
	for _, module := range modules {
		var (
			modulefs     filesystem.Filespace
			buildConfigs []*config.Build
		)
		if !fs.IsExist(module.SourceDir) {
			return goaterr.Errorf("builder.buildModules: Module '%s' is not exist", module.SourceDir)
		}
		if err = fs.MkdirAll(module.SourceDir, 0766); err != nil {
			return err
		}
		if modulefs, err = fs.Filespace(module.SourceDir); err != nil {
			return err
		}
		if buildConfigs, err = c.service.ReadDefFromFS(modulefs); err != nil {
			return err
		}
		moduleSubPath := subPath + module.SourceDir
		if err = c.build(modulefs, moduleSubPath, buildConfigs); err != nil {
			return err
		}
	}
	return nil
}
