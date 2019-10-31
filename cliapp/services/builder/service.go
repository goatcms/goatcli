package builder

import (
	"strconv"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
)

// Service build structure
type Service struct {
	deps struct {
		CWD             string                       `argument:"?cwd"`
		ExecutorLimit   string                       `argument:"?executor.limit"`
		TemplateService services.TemplateService     `dependency:"TemplateService"`
		Modules         services.ModulesService      `dependency:"ModulesService"`
		Dependencies    services.DependenciesService `dependency:"DependenciesService"`
		Repositories    services.RepositoriesService `dependency:"RepositoriesService"`
		VCSService      services.VCSService          `dependency:"VCSService"`
		Output          app.Output                   `dependency:"OutputService"`
	}
	limit int64
}

// ServiceFactory create new repositories instance
func ServiceFactory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &Service{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	if instance.deps.CWD == "" {
		instance.deps.CWD = "./"
	}
	if instance.deps.ExecutorLimit == "" {
		instance.limit = DefaultExecutorLimit
	} else {
		if instance.limit, err = strconv.ParseInt(instance.deps.ExecutorLimit, 10, 64); err != nil {
			return nil, err
		}
	}
	return services.BuilderService(instance), nil
}

// NewContext create new Context instance
func (s *Service) NewContext(scope app.Scope, appData services.ApplicationData, properties, secrets map[string]string) services.BuildContext {
	return &Context{
		scope:      scope,
		data:       appData.Plain,
		properties: properties,
		secrets:    secrets,
		service:    s,
		appModel:   appData.AM,
	}
}

// ReadDefFromFS return data definition
func (s *Service) ReadDefFromFS(fs filesystem.Filespace) (builds []*config.Build, err error) {
	var json []byte
	if !fs.IsFile(BuildDefPath) {
		return make([]*config.Build, 0), nil
	}
	if json, err = fs.ReadFile(BuildDefPath); err != nil {
		return nil, err
	}
	if builds, err = config.NewBuilds(json); err != nil {
		return nil, err
	}
	return builds, nil
}
