package simpletf

import (
	"math/rand"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
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
	deps.TemplateService.AddFunc("lower", strings.ToLower)
	deps.TemplateService.AddFunc("title", strings.ToTitle)
	deps.TemplateService.AddFunc("upper", strings.ToUpper)
	deps.TemplateService.AddFunc("lowerFirst", naming.ToLowerFirst)
	deps.TemplateService.AddFunc("upperFirst", naming.ToUpperFirst)
	deps.TemplateService.AddFunc("underscore", naming.ToUnderscore)
	deps.TemplateService.AddFunc("camelcase", naming.ToCamelCase)
	deps.TemplateService.AddFunc("camelcaself", naming.ToCamelCaseLF)
	deps.TemplateService.AddFunc("camelcaseuf", naming.ToCamelCaseUF)
	deps.TemplateService.AddFunc("regexp", Regexp)
	deps.TemplateService.AddFunc("strainMap", StrainMap)
	deps.TemplateService.AddFunc("random", varutil.RandString)
	deps.TemplateService.AddFunc("keys", Keys)
	deps.TemplateService.AddFunc("subMap", SubMap)
	deps.TemplateService.AddFunc("json", ToJSON)
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
	deps.TemplateService.AddFunc("wrap", Wrap)
	deps.TemplateService.AddFunc("safeHTMLAttr", SafeHTMLAttr)
	deps.TemplateService.AddFunc("safeHTML", SafeHTML)
	deps.TemplateService.AddFunc("safeURL", SafeURL)
	deps.TemplateService.AddFunc("dict", Dict)
	deps.TemplateService.AddFunc("sum", Sum)
	deps.TemplateService.AddFunc("sub", Sub)
	deps.TemplateService.AddFunc("log", Log)
	deps.TemplateService.AddFunc("logBuildFail", LogBuildFail)
	deps.TemplateService.AddFunc("logBuildSuccess", LogBuildSuccess)
	deps.TemplateService.AddFunc("logWarning", LogWarning)
	deps.TemplateService.AddFunc("contains", varutil.IsArrContainStr)
	deps.TemplateService.AddFunc("repeatNTimes", RepeatNTimes)
	deps.TemplateService.AddFunc("randIntn", rand.Intn)
	return nil
}
