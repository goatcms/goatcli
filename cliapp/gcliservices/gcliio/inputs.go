package gcliio

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common"

	"github.com/goatcms/goatcli/cliapp/common/gclivarutil"

	"github.com/goatcms/goatcore/varutil/goaterr"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/prevents"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
)

// Inputs run application scripts
type Inputs struct {
	deps struct {
		Interactive       string                         `argument:"?interactive" ,command:"?interactive"`
		PropertiesService gcliservices.PropertiesService `dependency:"PropertiesService"`
		SecretsService    gcliservices.SecretsService    `dependency:"SecretsService"`
		DataService       gcliservices.DataService       `dependency:"DataService"`
	}
}

// InputsFactory create new Inputs instance
func InputsFactory(dp app.DependencyProvider) (result interface{}, err error) {
	r := &Inputs{}
	if err = dp.InjectTo(&r.deps); err != nil {
		return nil, err
	}
	return gcliservices.GCLIInputs(r), nil
}

// Inputs return goat cli app inputs
func (inputs *Inputs) Inputs(ctx app.IOContext) (propertiesData, secretsData map[string]string, appData gcliservices.ApplicationData, err error) {
	var (
		propertiesDef     []*config.Property
		secretsDef        []*config.Property
		isChanged         bool
		interactive       bool
		fs                = ctx.IO().CWD()
		data              map[string]string
		elasticProperties common.ElasticData
	)
	interactive = strings.ToLower(inputs.deps.Interactive) != "false"
	if err = prevents.RequireGoatProject(fs); err != nil {
		return nil, nil, appData, err
	}
	// load properties
	if propertiesDef, err = inputs.deps.PropertiesService.ReadDefFromFS(fs); err != nil {
		return nil, nil, appData, err
	}
	if propertiesData, err = inputs.deps.PropertiesService.ReadDataFromFS(fs); err != nil {
		return nil, nil, appData, err
	}
	if isChanged, err = inputs.deps.PropertiesService.FillData(ctx, propertiesDef, propertiesData, map[string]string{}, interactive); err != nil {
		return nil, nil, appData, err
	}
	if isChanged {
		if err = inputs.deps.PropertiesService.WriteDataToFS(fs, propertiesData); err != nil {
			return nil, nil, appData, err
		}
	}
	// load data
	if data, err = inputs.deps.DataService.ReadDataFromFS(fs); err != nil {
		return nil, nil, appData, err
	}
	if appData, err = am.NewApplicationData(data); err != nil {
		return nil, nil, appData, goaterr.Wrapf(err, "plaindata contains illegal keys or is incorrect")
	}
	// load secrets
	if elasticProperties, err = gclivarutil.NewElasticData(propertiesData); err != nil {
		return nil, nil, appData, err
	}
	if secretsDef, err = inputs.deps.SecretsService.ReadDefFromFS(fs, elasticProperties, appData); err != nil {
		return nil, nil, appData, err
	}
	if secretsData, err = inputs.deps.SecretsService.ReadDataFromFS(fs); err != nil {
		return nil, nil, appData, err
	}
	if isChanged, err = inputs.deps.SecretsService.FillData(ctx, secretsDef, secretsData, map[string]string{}, interactive); err != nil {
		return nil, nil, appData, err
	}
	if isChanged {
		if err = inputs.deps.SecretsService.WriteDataToFS(fs, secretsData); err != nil {
			return nil, nil, appData, err
		}
	}
	return propertiesData, secretsData, appData, nil
}
