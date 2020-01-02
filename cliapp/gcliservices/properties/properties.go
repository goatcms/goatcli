package properties

import (
	"os"
	"sort"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/cio"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/plainmap"
)

// Properties provide project properties data
type Properties struct {
	deps struct {
		FS filesystem.Filespace `filespace:"current"`
	}
}

// Factory create new repositories instance
func Factory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &Properties{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return gcliservices.PropertiesService(instance), nil
}

// ReadDefFromFS read properties definitions from filespace
func (p *Properties) ReadDefFromFS(fs filesystem.Filespace) (properties []*config.Property, err error) {
	var (
		json  []byte
		nodes []os.FileInfo
		props []*config.Property
	)
	if fs.IsFile(PropertiesDefPath) {
		if json, err = fs.ReadFile(PropertiesDefPath); err != nil {
			return nil, err
		}
		if properties, err = config.NewProperties(json); err != nil {
			return nil, err
		}
	} else {
		properties = make([]*config.Property, 0)
	}
	// Read separated data.def files
	if !fs.IsDir(BasePropertiesDefPath) {
		return properties, nil
	}
	if nodes, err = fs.ReadDir(BasePropertiesDefPath); err != nil {
		return nil, err
	}
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].Name() < nodes[j].Name()
	})
	for _, node := range nodes {
		var path = BasePropertiesDefPath + node.Name()
		if !fs.IsFile(path) || !strings.HasSuffix(path, PropertiesDefSuffix) {
			continue
		}
		if json, err = fs.ReadFile(path); err != nil {
			return nil, err
		}
		if props, err = config.NewProperties(json); err != nil {
			return nil, err
		}
		properties = append(properties, props...)
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
func (p *Properties) FillData(ctx app.IOContext, def []*config.Property, data map[string]string, defaultData map[string]string, interactive bool) (isChanged bool, err error) {
	io := ctx.IO()
	return cio.ReadProperties("", io.In(), io.Out(), def, data, defaultData, interactive)
}

// WriteDataToFS write properties data to fs file
func (p *Properties) WriteDataToFS(fs filesystem.Filespace, data map[string]string) (err error) {
	var json string
	if json, err = plainmap.PlainStringMapToFormattedJSON(data); err != nil {
		return err
	}
	return fs.WriteFile(PropertiesDataPath, []byte(json), 0766)
}
