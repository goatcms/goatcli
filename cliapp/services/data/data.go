package compiler

import (
	"regexp"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
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
	return services.Data(instance), nil
}

// ReadDefFromFS return data definition
func (d *Data) ReadDefFromFS(fs filesystem.Filespace) (dataSets []*config.DataSet, err error) {
	var json []byte
	if !fs.IsFile(DataDefPath) {
		return make([]*config.DataSet, 0), nil
	}
	if json, err = fs.ReadFile(DataDefPath); err != nil {
		return nil, err
	}
	if dataSets, err = config.NewDataSets(json); err != nil {
		return nil, err
	}
	return dataSets, nil
}

// ReadDataFromFS return data
func (d *Data) ReadDataFromFS(fs filesystem.Filespace) (data map[string]string, err error) {
	data = make(map[string]string)
	if err = readDataFromFS(data, fs, ".goat/data/"); err != nil {
		return nil, err
	}
	return data, nil
}

// ConsoleReadData create new data from Filespace
func (d *Data) ConsoleReadData(def *config.DataSet) (data map[string]string, err error) {
	var line string
	data = make(map[string]string)
	for _, property := range def.Properties {
		for {
			d.deps.Output.Printf("%s: ", property.Prompt)
			if line, err = d.deps.Input.ReadLine(); err != nil {
				return nil, err
			}
			if len(line) > property.Max {
				d.deps.Output.Printf("Max value length is %d (value length is %d)\n", property.Max, len(line))
				continue
			}
			if len(line) < property.Min {
				d.deps.Output.Printf("Min value length is %d (value length is %d)\n", property.Min, len(line))
				continue
			}
			if property.Type == "numeric" && !numericReg.MatchString(line) {
				d.deps.Output.Printf("Require numeric value\n")
				continue
			}
			if property.Type == "alpha" && !alphaReg.MatchString(line) {
				d.deps.Output.Printf("Require alpha-numeric value\n")
				continue
			}
			if property.Type == "alnum" && !alnumReg.MatchString(line) {
				d.deps.Output.Printf("Require alpha-numeric value\n")
				continue
			}
			data[property.Key] = line
			break
		}
	}
	return data, nil
}
