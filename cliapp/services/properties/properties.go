package properties

import (
	"io"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/result"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil"
	"github.com/goatcms/goatcore/varutil/plainmap"
)

// Properties provide project properties data
type Properties struct {
	deps struct {
		FS     filesystem.Filespace `filespace:"root"`
		Input  app.Input            `dependency:"InputService"`
		Output app.Output           `dependency:"OutputService"`
	}
}

// Factory create new repositories instance
func Factory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &Properties{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return services.Properties(instance), nil
}

func (p *Properties) Get(fs filesystem.Filespace) (common.PropertiesResult, error) {
	var (
		data map[string]string
		err  error
	)
	defJSON, err := p.deps.FS.ReadFile(".goat/properties.def.json")
	if err != nil {
		return nil, err
	}
	properties, err := config.NewProperties(defJSON)
	if err != nil {
		return nil, err
	}
	if p.deps.FS.IsFile(".goat/properties.json") {
		var dataJSON []byte
		dataJSON, err = p.deps.FS.ReadFile(".goat/properties.json")
		if err != nil {
			return nil, err
		}
		data, err = plainmap.JSONToPlainStringMap(dataJSON)
		if err != nil {
			return nil, err
		}
	} else {
		data = make(map[string]string)
	}
	isChanged := false
	for _, property := range properties {
		if _, ok := data[property.Key]; !ok {
			var genvalue string
			var input string
			// load data for property
			if !isChanged {
				p.deps.Output.Printf("Insert lost properties:\n")
				isChanged = true
			}
			switch strings.ToLower(property.Type) {
			case "numeric":
				genvalue = varutil.RandString(property.Max, varutil.NumericBytes)
			case "alpha":
				genvalue = varutil.RandString(property.Max, varutil.AlphaBytes)
			case "alnum":
				genvalue = varutil.RandString(property.Max, varutil.AlphaNumericBytes)
			case "strong":
				genvalue = varutil.RandString(property.Max, varutil.StrongBytes)
			}
			for {
				p.deps.Output.Printf(">%s [%s]: ", property.Key, genvalue)
				if input, err = p.deps.Input.ReadLine(); err != nil && err != io.EOF {
					return nil, err
				}
				if input == "" {
					data[property.Key] = genvalue
					break
				}
				if len(input) < property.Min {
					p.deps.Output.Printf("Value is too short. Minimum length of the value is %d characters.\n", property.Min)
					continue
				}
				if len(input) > property.Max {
					p.deps.Output.Printf("Value is too long. Maximum length of the value is %d characters.\n", property.Max)
					continue
				}
				data[property.Key] = input
				break
			}
		}
	}
	if isChanged {
		json, err := plainmap.PlainStringMapToJSON(data)
		if err != nil {
			return nil, err
		}
		p.deps.FS.WriteFile(".goat/properties.json", []byte(json), 0766)
	}
	return result.NewPropertiesResult(data), nil
}
