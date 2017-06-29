package properties

import (
	"fmt"
	"io"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/config"
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
	return services.PropertiesService(instance), nil
}

// ReadDefFromFS read properties definitions from filespace
func (p *Properties) ReadDefFromFS(fs filesystem.Filespace) (properties []*config.Property, err error) {
	var json []byte
	if !fs.IsFile(PropertiesDefPath) {
		return make([]*config.Property, 0), nil
	}
	if json, err = fs.ReadFile(PropertiesDefPath); err != nil {
		return nil, err
	}
	if properties, err = config.NewProperties(json); err != nil {
		return nil, err
	}
	return properties, nil
}

// ReadDataFromFS read properties data from filespace
func (p *Properties) ReadDataFromFS(fs filesystem.Filespace) (data map[string]string, err error) {
	var json []byte
	if !fs.IsFile(PropertiesDataPath) {
		return make(map[string]string, 0), nil
	}
	if json, err = fs.ReadFile(PropertiesDataPath); err != nil {
		return nil, err
	}
	if data, err = plainmap.JSONToPlainStringMap(json); err != nil {
		return nil, err
	}
	return data, nil
}

// FillData read lost properties data to curent data map
func (p *Properties) FillData(def []*config.Property, data map[string]string, defaultData map[string]string) (isChanged bool, err error) {
	var (
		ok           bool
		defaultValue string
		input        string
	)
	for _, property := range def {
		if _, ok = data[property.Key]; ok {
			continue
		}
		if defaultValue, ok = data[property.Key]; !ok {
			switch strings.ToLower(property.Type) {
			case "numeric":
				defaultValue = varutil.RandString(property.Max, varutil.NumericBytes)
			case "alpha":
				defaultValue = varutil.RandString(property.Max, varutil.AlphaBytes)
			case "alnum":
				defaultValue = varutil.RandString(property.Max, varutil.AlphaNumericBytes)
			case "strong":
				defaultValue = varutil.RandString(property.Max, varutil.StrongBytes)
			default:
				return isChanged, fmt.Errorf("wrong property type %s (for property %s)", property.Type, property.Key)
			}
		}
		for {
			p.deps.Output.Printf("Insert property %s [%s]: ", property.Key, defaultValue)
			if input, err = p.deps.Input.ReadLine(); err != nil && err != io.EOF {
				return isChanged, err
			}
			if input == "" {
				isChanged = true
				data[property.Key] = defaultValue
				break
			}
			if len(input) < property.Min {
				p.deps.Output.Printf("Value is too short. Minimum length of the property value is %d characters.\n", property.Min)
				continue
			}
			if len(input) > property.Max {
				p.deps.Output.Printf("Value is too long. Maximum length of the property value is %d characters.\n", property.Max)
				continue
			}
			isChanged = true
			data[property.Key] = input
			break
		}
	}
	return isChanged, nil
}

// WriteDataToFS write properties data to fs file
func (p *Properties) WriteDataToFS(fs filesystem.Filespace, data map[string]string) (err error) {
	var json string
	if json, err = plainmap.PlainStringMapToJSON(data); err != nil {
		return err
	}
	if err = fs.WriteFile(".goat/properties.json", []byte(json), 0766); err != nil {
		return err
	}
	return nil
}
