package simpletf

import (
	"math/rand"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/varutil"
)

// RegisterFunctions add default helper functions for templates
func RegisterFunctions(di dependency.Injector) (err error) {
	var deps struct {
		Config gcliservices.TemplateConfig `dependency:"TemplateConfig"`
	}
	if err = di.InjectTo(&deps); err != nil {
		return err
	}
	deps.Config.AddFunc("replace", Replace)
	deps.Config.AddFunc("injectValues", InjectValues)
	deps.Config.AddFunc("join", strings.Join)
	deps.Config.AddFunc("split", strings.Split)
	deps.Config.AddFunc("hasPrefix", strings.HasPrefix)
	deps.Config.AddFunc("hasSuffix", strings.HasSuffix)
	deps.Config.AddFunc("lower", strings.ToLower)
	deps.Config.AddFunc("title", strings.ToTitle)
	deps.Config.AddFunc("upper", strings.ToUpper)
	deps.Config.AddFunc("lowerFirst", naming.ToLowerFirst)
	deps.Config.AddFunc("upperFirst", naming.ToUpperFirst)
	deps.Config.AddFunc("underscore", naming.ToUnderscore)
	deps.Config.AddFunc("camelcase", naming.ToCamelCase)
	deps.Config.AddFunc("camelcaself", naming.ToCamelCaseLF)
	deps.Config.AddFunc("camelcaseuf", naming.ToCamelCaseUF)
	deps.Config.AddFunc("regexp", Regexp)
	deps.Config.AddFunc("strainMap", StrainMap)
	deps.Config.AddFunc("random", varutil.RandString)
	deps.Config.AddFunc("keys", Keys)
	deps.Config.AddFunc("subMap", SubMap)
	deps.Config.AddFunc("json", ToJSON)
	deps.Config.AddFunc("error", ToError)
	deps.Config.AddFunc("union", Union)
	deps.Config.AddFunc("unique", Unique)
	deps.Config.AddFunc("except", Except)
	deps.Config.AddFunc("sort", Sort)
	deps.Config.AddFunc("intersect", Intersect)
	deps.Config.AddFunc("randomValue", RandomValue)
	deps.Config.AddFunc("valuesFor", ValuesFor)
	deps.Config.AddFunc("findRow", FindRow)
	deps.Config.AddFunc("findRows", FindRows)
	deps.Config.AddFunc("noescape", Noescape)
	deps.Config.AddFunc("wrap", Wrap)
	deps.Config.AddFunc("safeHTMLAttr", SafeHTMLAttr)
	deps.Config.AddFunc("safeHTML", SafeHTML)
	deps.Config.AddFunc("safeURL", SafeURL)
	deps.Config.AddFunc("dict", Dict)
	deps.Config.AddFunc("sum", Sum)
	deps.Config.AddFunc("sub", Sub)
	deps.Config.AddFunc("log", Log)
	deps.Config.AddFunc("logBuildFail", LogBuildFail)
	deps.Config.AddFunc("logBuildSuccess", LogBuildSuccess)
	deps.Config.AddFunc("logWarning", LogWarning)
	deps.Config.AddFunc("contains", varutil.IsArrContainStr)
	deps.Config.AddFunc("repeatNTimes", RepeatNTimes)
	deps.Config.AddFunc("randIntn", rand.Intn)
	deps.Config.AddFunc("sequencer", NewSequencer)
	return nil
}
