package builder

import (
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
	return services.BuilderService(instance), nil
}

// Build project files and directories from data
func (s *Service) Build(fs filesystem.Filespace, buildConfigs []*config.Build, data map[string]string) (err error) {
	var (
		templateExecutor services.TemplateExecutor
		writer           *FSWriter
		hash             string
	)
	hash = varutil.RandString(20, varutil.AlphaNumericBytes)
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
				Build: c.Properties,
			},
		})
		if err = templateExecutor.Execute(c.Layout, c.View, writer, context); err != nil {
			return err
		}
	}
	return nil
}
