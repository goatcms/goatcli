package tfunc

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/varutil"
)

// RegisterFunctions add default helper functions for templates
func RegisterFunctions(di dependency.Injector) (err error) {
	var deps struct {
		TemplateService services.TemplateService `dependency:"TemplateService"`
	}
	if err = di.InjectTo(&deps); err != nil {
		return err
	}
	deps.TemplateService.AddFunc("replace", Replace)
	deps.TemplateService.AddFunc("injectValues", InjectValues)
	deps.TemplateService.AddFunc("join", strings.Join)
	deps.TemplateService.AddFunc("split", strings.Split)
	deps.TemplateService.AddFunc("hasPrefix", strings.HasPrefix)
	deps.TemplateService.AddFunc("hasSuffix", strings.HasSuffix)
	deps.TemplateService.AddFunc("hasSuffix", strings.ToLower)
	deps.TemplateService.AddFunc("lower", strings.ToLower)
	deps.TemplateService.AddFunc("lowerFirst", ToLowerFirst)
	deps.TemplateService.AddFunc("upper", strings.ToUpper)
	deps.TemplateService.AddFunc("upperFirst", ToUpperFirst)
	deps.TemplateService.AddFunc("underscore", ToUnderscore)
	deps.TemplateService.AddFunc("camelcase", ToCamelCase)
	deps.TemplateService.AddFunc("camelcaself", ToCamelCaseLF)
	deps.TemplateService.AddFunc("camelcaseuf", ToCamelCaseUF)
	deps.TemplateService.AddFunc("title", strings.ToTitle)
	deps.TemplateService.AddFunc("regexp", Regexp)
	deps.TemplateService.AddFunc("strainMap", StrainMap)
	deps.TemplateService.AddFunc("random", varutil.RandString)
	deps.TemplateService.AddFunc("keys", Keys)
	deps.TemplateService.AddFunc("ctx", NewContext)
	deps.TemplateService.AddFunc("error", ToError)
	deps.TemplateService.AddFunc("union", Union)
	deps.TemplateService.AddFunc("unique", Unique)
	deps.TemplateService.AddFunc("except", Except)
	deps.TemplateService.AddFunc("sort", Sort)
	deps.TemplateService.AddFunc("intersect", Intersect)
	deps.TemplateService.AddFunc("randomValue", RandomValue)
	deps.TemplateService.AddFunc("valuesFor", ValuesFor)
	deps.TemplateService.AddFunc("findRow", FindRow)
	deps.TemplateService.AddFunc("findRows", FindRows)
	deps.TemplateService.AddFunc("noescape", Noescape)
	deps.TemplateService.AddFunc("dict", Dict)
	return nil
}
