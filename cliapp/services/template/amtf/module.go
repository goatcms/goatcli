package amtf

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
)

// RegisterFunctions add default helper functions for templates
func RegisterFunctions(di dependency.Injector) (err error) {
	var deps struct {
		Config services.TemplateConfig `dependency:"TemplateConfig"`
	}
	if err = di.InjectTo(&deps); err != nil {
		return err
	}
	deps.Config.AddFunc("amLinkFieldUF", LinkFieldUF)
	deps.Config.AddFunc("amLinkFieldLF", LinkFieldLF)
	deps.Config.AddFunc("amLinkRelationUF", LinkRelationUF)
	deps.Config.AddFunc("amLinkRelationLF", LinkRelationLF)
	deps.Config.AddFunc("amStructClassName", StructClassName)
	return nil
}
