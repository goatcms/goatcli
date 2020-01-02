package data

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/cio"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/plainmap"
)

var (
	numericReg = regexp.MustCompile("^[0-9]+$")
	alphaReg   = regexp.MustCompile("^[A-Za-z]+$")
	alnumReg   = regexp.MustCompile("^[A-Za-z0-9]+$")
)

// Data provider
type Data struct {
	deps struct {
		Input  app.Input  `dependency:"InputService"`
		Output app.Output `dependency:"OutputService"`
	}
}

// BuilderFactory create new repositories instance
func BuilderFactory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &Data{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return gcliservices.DataService(instance), nil
}

// HasDataFromFS return true if storage data by prefix
func (d *Data) HasDataFromFS(fs filesystem.Filespace, prefix string) bool {
	path := BaseDataPath + strings.Replace(prefix, ".", "/", -1) + ".json"
	return fs.IsExist(path)
}

// ReadDefFromFS return data definition
func (d *Data) ReadDefFromFS(fs filesystem.Filespace) (dataSets []*config.DataSet, err error) {
	var (
		json  []byte
		nodes []os.FileInfo
		sets  []*config.DataSet
	)
	if fs.IsFile(DataDefPath) {
		if json, err = fs.ReadFile(DataDefPath); err != nil {
			return nil, err
		}
		if dataSets, err = config.NewDataSets(json); err != nil {
			return nil, err
		}
	} else {
		dataSets = make([]*config.DataSet, 0)
	}
	// Read separated data.def files
	if !fs.IsDir(BaseDataDefPath) {
		return dataSets, nil
	}
	if nodes, err = fs.ReadDir(BaseDataDefPath); err != nil {
		return nil, err
	}
	for _, node := range nodes {
		var path = BaseDataDefPath + node.Name()
		if !fs.IsFile(path) || !strings.HasSuffix(path, DataDefSuffix) {
			continue
		}
		if json, err = fs.ReadFile(path); err != nil {
			return nil, err
		}
		if sets, err = config.NewDataSets(json); err != nil {
			return nil, err
		}
		dataSets = append(dataSets, sets...)
	}
	return dataSets, nil
}

// ReadDataFromFS return data
func (d *Data) ReadDataFromFS(fs filesystem.Filespace) (data map[string]string, err error) {
	data = make(map[string]string)
	if !fs.IsDir(BaseDataPath) {
		return map[string]string{}, nil
	}
	if err = readDataFromFS(data, fs, BaseDataPath); err != nil {
		return nil, err
	}
	return data, nil
}

// ConsoleReadData create new data from Filespace
func (d *Data) ConsoleReadData(def *config.DataSet) (data map[string]string, err error) {
	data = make(map[string]string)
	if _, err = cio.ReadDataSet("", d.deps.Input, d.deps.Output, def, data); err != nil {
		return nil, err
	}
	return data, nil
}

// WriteDataToFS write data to filespace
func (d *Data) WriteDataToFS(fs filesystem.Filespace, prefix string, data map[string]string) (err error) {
	var json string
	outmap := map[string]string{}
	path := BaseDataPath + strings.Replace(prefix, ".", "/", -1) + ".json"
	prefix += "."
	for key, value := range data {
		outmap[prefix+key] = value
	}
	if json, err = plainmap.PlainStringMapToFormattedJSON(outmap); err != nil {
		return err
	}
	if fs.IsExist(path) {
		return fmt.Errorf("DataService.WriteDataToFS: %s exists", path)
	}
	if err = fs.MkdirAll(filepath.Dir(path), 0766); err != nil {
		return err
	}
	return fs.WriteFile(path, []byte(json), 0766)
}
