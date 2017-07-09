package tfunc

import (
	"strings"

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
	deps.TemplateService.AddFunc("join", strings.Join)
	deps.TemplateService.AddFunc("hasPrefix", strings.HasPrefix)
	deps.TemplateService.AddFunc("hasSuffix", strings.HasSuffix)
	deps.TemplateService.AddFunc("regexp", Regexp)
	deps.TemplateService.AddFunc("strainMap", StrainMap)
	deps.TemplateService.AddFunc("keys", Keys)
	return nil
}
