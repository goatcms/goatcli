package properties

import (
	"github.com/goatcms/goatcli/cliapp/common/cio"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
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
func (p *Properties) FillData(def []*config.Property, data map[string]string, defaultData map[string]string, interactive bool) (isChanged bool, err error) {
	return cio.ReadProperties("", p.deps.Input, p.deps.Output, def, data, defaultData, interactive)
}

// WriteDataToFS write properties data to fs file
func (p *Properties) WriteDataToFS(fs filesystem.Filespace, data map[string]string) (err error) {
	var json string
	if json, err = plainmap.PlainStringMapToFormattedJSON(data); err != nil {
		return err
	}
	return fs.WriteFile(PropertiesDataPath, []byte(json), 0766)
}
