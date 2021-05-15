package shtf

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
)

// RegisterFunctions add default helper functions for templates
func RegisterFunctions(di app.DependencyProvider) (err error) {
	var deps struct {
		Config gcliservices.TemplateConfig `dependency:"TemplateConfig"`
	}
	if err = di.InjectTo(&deps); err != nil {
		return err
	}
	deps.Config.AddFunc("sh_secret", SecretEnv)
	return nil
}
