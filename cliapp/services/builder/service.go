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
		CWD             string                   `argument:"?cwd"`
		TemplateService services.TemplateService `dependency:"TemplateService"`
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

// Build project files and directories from data
func (s *Service) Build(fs filesystem.Filespace, buildConfigs []*config.Build, data, properties map[string]string) (err error) {
	var (
		templateExecutor services.TemplateExecutor
		writer           *FSWriter
		hash             string
	)
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
			for i, _ := range args {
				// replace it here because argument.cwd can contains space (for example in home directory name)
				args[i] = strings.Replace(args[i], "{{argument.cwd}}", s.deps.CWD, -1)
			}
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = &out
			cmd.Stderr = &out
			if err = cmd.Run(); err != nil {
				fmt.Printf("external app fail %v: %v %v\n", args, err, string(out.Bytes()))
				// ignore errors from external apps (only show it)
			}
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
