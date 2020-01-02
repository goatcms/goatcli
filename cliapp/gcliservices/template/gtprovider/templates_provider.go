package gtprovider

import (
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/fsloop"
	"github.com/goatcms/goatcore/goathtml"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// TemplatesProvider provide templates api
type TemplatesProvider struct {
	fs             filesystem.Filespace
	assetsProvider gcliservices.TemplateAssetsProvider
	templatesPath  string
	templatesMutex sync.Mutex
	templates      map[string]*template.Template
	//funcs          template.FuncMap
	isCached bool
}

// NewTemplatesProvider create TemplatesProvider instance
func NewTemplatesProvider(assetsProvider gcliservices.TemplateAssetsProvider, fs filesystem.Filespace, templatesPath string, isCached bool) *TemplatesProvider {
	return &TemplatesProvider{
		fs:             fs,
		assetsProvider: assetsProvider,
		templates:      map[string]*template.Template{},
		templatesPath:  templatesPath,
		isCached:       isCached,
	}
}

// Template return template by name. It contains selected layout definitions and helpers
func (provider *TemplatesProvider) Template(layoutName, templateName string) (tmpl *template.Template, err error) {
	var (
		ok  bool
		key string
	)
	if layoutName == "" {
		layoutName = goathtml.DefaultLayout
	}
	if templateName == "" {
		return nil, goaterr.Errorf("goathtml.TemplatesProvider: A template name is required")
	}
	key = layoutName + ":" + templateName
	if tmpl, ok = provider.templates[key]; ok {
		return tmpl, nil
	}
	return provider.template(layoutName, templateName, key)
}

func (provider *TemplatesProvider) template(layoutName, templateName, key string) (tmpl *template.Template, err error) {
	var (
		ok    bool
		path  string
		subFS filesystem.Filespace
	)
	provider.templatesMutex.Lock()
	defer provider.templatesMutex.Unlock()
	// check after lock
	if tmpl, ok = provider.templates[key]; ok {
		return tmpl, nil
	}
	// create a new template
	if tmpl, err = provider.assetsProvider.Layout(layoutName); err != nil {
		return nil, err
	}
	path = strings.Replace(provider.templatesPath, "{name}", templateName, 1)
	if !provider.fs.IsDir(path) {
		if provider.isCached {
			provider.templates[key] = tmpl
		}
		return tmpl, nil
	}
	if tmpl, err = tmpl.Clone(); err != nil {
		return nil, err
	}
	if subFS, err = provider.fs.Filespace(path); err != nil {
		return nil, err
	}
	templateLoader := NewTemplateLoader(tmpl)
	if err = fsloop.WalkFS(subFS, "", func(currentPath string, info os.FileInfo) (err error) {
		_, err = templateLoader.LoadByExtension(subFS, currentPath)
		return err
	}, nil); err != nil {
		return nil, err
	}
	if provider.isCached {
		provider.templates[key] = tmpl
	}
	return tmpl, nil
}
