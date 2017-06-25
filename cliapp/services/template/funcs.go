package template

import (
	"github.com/goatcms/goatcms/cmsapp/services"
	"github.com/goatcms/goatcore/dependency"
)

// DefaultFuncs is a container for template helper functions
type DefaultFuncs struct {
	deps struct {
		Template services.Template `dependency:"TemplateService"`
	}
}

// NewDefaultFuncs create new instance of DefaultFuncs
func NewDefaultFuncs(di dependency.Injector) (*DefaultFuncs, error) {
	df := &DefaultFuncs{}
	if err := di.InjectTo(df.deps); err != nil {
		return nil, err
	}
	return df, nil
}

// RegisterFunctions add default helper functions for templates
func RegisterFunctions(di dependency.Injector) (err error) {
	//var functions *DefaultFuncs
	if _, err = NewDefaultFuncs(di); err != nil {
		return err
	}
	//df.Template.AddFunc(services.CutTextTF, df.CutText)
	//df.Template.AddFunc(services.MessagesTF, df.Messages)
	return nil
}
