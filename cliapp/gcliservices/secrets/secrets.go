package secrets

import (
	"os"
	"sort"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/cio"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/secrets/executor"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/plainmap"
)

// Secrets provide project secrets data
type Secrets struct {
	deps struct {
		FS              filesystem.Filespace         `filespace:"current"`
		TemplateService gcliservices.TemplateService `dependency:"TemplateService"`
	}
}

// Factory create new repositories instance
func Factory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &Secrets{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return gcliservices.SecretsService(instance), nil
}

// ReadDefFromFS read secrets definitions from filespace
func (p *Secrets) ReadDefFromFS(fs filesystem.Filespace, properties map[string]string, appData gcliservices.ApplicationData) (secrets []*config.Property, err error) {
	var (
		json  []byte
		nodes []os.FileInfo
		props []*config.Property
	)
	if fs.IsFile(SecretsDefPath) {
		if json, err = fs.ReadFile(SecretsDefPath); err != nil {
			return nil, err
		}
		if secrets, err = config.NewProperties(json); err != nil {
			return nil, err
		}
	} else {
		secrets = make([]*config.Property, 0)
	}
	// Read separated data.def files
	if !fs.IsDir(BaseSecretsDefPath) {
		return secrets, nil
	}
	if nodes, err = fs.ReadDir(BaseSecretsDefPath); err != nil {
		return nil, err
	}
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].Name() < nodes[j].Name()
	})
	for _, node := range nodes {
		var path = BaseSecretsDefPath + node.Name()
		if !fs.IsFile(path) || !strings.HasSuffix(path, SecretsDefSuffix) {
			continue
		}
		if json, err = fs.ReadFile(path); err != nil {
			return nil, err
		}
		if props, err = config.NewProperties(json); err != nil {
			return nil, err
		}
		secrets = append(secrets, props...)
	}
	// Generate secret template
	if props, err = p.readDefFromFTemplate(fs, properties, appData); err != nil {
		return nil, err
	}
	secrets = append(secrets, props...)
	return secrets, nil
}

// ReadDefFromFS read secrets definitions from filespace
func (p *Secrets) readDefFromFTemplate(fs filesystem.Filespace, properties map[string]string, appData gcliservices.ApplicationData) (secrets []*config.Property, err error) {
	var (
		templateExecutor gcliservices.TemplateExecutor
		secretsExecutor  *executor.SecretsExecutor
		executorScope    = scope.NewScope(scope.Params{})
	)
	if templateExecutor, err = p.deps.TemplateService.TemplateExecutor(".goat/secrets.def"); err != nil {
		return nil, err
	}
	if secretsExecutor, err = executor.NewSecretsExecutor(executorScope, executor.SharedData{
		AppData: appData,
		Properties: executor.GlobalProperties{
			Project: properties,
		},
		DotData: nil,
	}, 10, templateExecutor); err != nil {
		return nil, err
	}
	if err = secretsExecutor.Execute(); err != nil {
		return nil, err
	}
	if err = executorScope.Wait(); err != nil {
		return nil, err
	}
	return secretsExecutor.Secrets()
}

// ReadDataFromFS read secrets data from filespace
func (p *Secrets) ReadDataFromFS(fs filesystem.Filespace) (data map[string]string, err error) {
	var json []byte
	if !fs.IsFile(SecretsDataPath) {
		return make(map[string]string, 0), nil
	}
	if json, err = fs.ReadFile(SecretsDataPath); err != nil {
		return nil, err
	}
	if data, err = plainmap.JSONToPlainStringMap(json); err != nil {
		return nil, err
	}
	return data, nil
}

// FillData read lost secrets data to curent data map
func (p *Secrets) FillData(ctx app.IOContext, def []*config.Property, data map[string]string, defaultData map[string]string, interactive bool) (isChanged bool, err error) {
	io := ctx.IO()
	return cio.ReadProperties("", io.In(), io.Out(), def, data, defaultData, interactive)
}

// WriteDataToFS write secrets data to fs file
func (p *Secrets) WriteDataToFS(fs filesystem.Filespace, data map[string]string) (err error) {
	var json string
	if json, err = plainmap.PlainStringMapToFormattedJSON(data); err != nil {
		return err
	}
	return fs.WriteFile(SecretsDataPath, []byte(json), 0766)
}
