package amtf

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
)

// RegisterFunctions add default helper functions for templates
func RegisterFunctions(di dependency.Injector) (err error) {
	var deps struct {
		TemplateService services.TemplateService `dependency:"TemplateService"`
	}
	if err = di.InjectTo(&deps); err != nil {
		return err
	}
	deps.TemplateService.AddFunc("amLinkFieldUF", LinkFieldUF)
	deps.TemplateService.AddFunc("amLinkFieldLF", LinkFieldLF)
	deps.TemplateService.AddFunc("amStructClassName", StructClassName)
	return nil
}
