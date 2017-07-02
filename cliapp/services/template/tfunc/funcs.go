package tfunc

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
	//deps.TemplateService.AddFunc("IsFile", IsFile)
	//deps.TemplateService.AddFunc("ISDir", IsDir)
	//deps.TemplateService.AddFunc("IsExist", IsExist)
	return nil
}
