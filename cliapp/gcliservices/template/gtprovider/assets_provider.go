package gtprovider

import (
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/fsloop"
	"github.com/goatcms/goatcore/goathtml"
)

// AssetsProvider provide templates api
type AssetsProvider struct {
	fs           filesystem.Filespace
	helpersPath  string
	layoutPath   string
	baseMutex    sync.Mutex
	baseTemplate *template.Template
	layoutMutex  sync.Mutex
	layouts      map[string]*template.Template
	funcs        template.FuncMap
	isCached     bool
}

// NewAssetsProvider create AssetsProvider instance
func NewAssetsProvider(fs filesystem.Filespace, helpersPath, layoutPath string, funcs template.FuncMap, isCached bool) *AssetsProvider {
	return &AssetsProvider{
		fs:          fs,
		layoutPath:  layoutPath,
		helpersPath: helpersPath,
		layouts:     map[string]*template.Template{},
		funcs:       funcs,
		isCached:    isCached,
	}
}

// Base return base template (with loaded helpers)
func (provider *AssetsProvider) Base() (*template.Template, error) {
	if provider.baseTemplate != nil {
		return provider.baseTemplate, nil
	}
	return provider.base()
}

func (provider *AssetsProvider) base() (baseTemplate *template.Template, err error) {
	var subFS filesystem.Filespace
	provider.baseMutex.Lock()
	defer provider.baseMutex.Unlock()
	if provider.baseTemplate != nil {
		return provider.baseTemplate, nil
	}
	baseTemplate = template.New("baseTemplate")
	baseTemplate.Funcs(provider.funcs)
	if !provider.fs.IsDir(provider.helpersPath) {
		if provider.isCached {
			provider.baseTemplate = baseTemplate
		}
		return baseTemplate, nil
	}
	templateLoader := NewTemplateLoader(baseTemplate)
	if subFS, err = provider.fs.Filespace(provider.helpersPath); err != nil {
		return nil, err
	}
	if err = fsloop.WalkFS(subFS, "", func(currentPath string, info os.FileInfo) (err error) {
		_, err = templateLoader.LoadByExtension(subFS, currentPath)
		return err
	}, nil); err != nil {
		return nil, err
	}
	if provider.isCached {
		provider.baseTemplate = baseTemplate
	}
	return baseTemplate, nil
}

// Layout return template for named layout (with loaded helpers and layout definitions)
func (provider *AssetsProvider) Layout(name string) (tmpl *template.Template, err error) {
	var (
		ok bool
	)
	if name == "" {
		name = goathtml.DefaultLayout
	}
	if tmpl, ok = provider.layouts[name]; ok {
		return tmpl, nil
	}
	return provider.layout(name)
}

func (provider *AssetsProvider) layout(name string) (layoutTemplate *template.Template, err error) {
	var (
		ok    bool
		subFS filesystem.Filespace
	)
	provider.layoutMutex.Lock()
	defer provider.layoutMutex.Unlock()
	if layoutTemplate, ok = provider.layouts[name]; ok {
		return layoutTemplate, nil
	}
	if layoutTemplate, err = provider.Base(); err != nil {
		return nil, err
	}
	path := strings.Replace(provider.layoutPath, "{name}", name, 1)
	if !provider.fs.IsDir(path) {
		if provider.isCached {
			provider.layouts[name] = layoutTemplate
		}
		return layoutTemplate, nil
	}
	if layoutTemplate, err = layoutTemplate.Clone(); err != nil {
		return nil, err
	}
	templateLoader := NewTemplateLoader(layoutTemplate)
	if subFS, err = provider.fs.Filespace(path); err != nil {
		return nil, err
	}
	if err = fsloop.WalkFS(subFS, "", func(currentPath string, info os.FileInfo) (err error) {
		_, err = templateLoader.LoadByExtension(subFS, currentPath)
		return err
	}, nil); err != nil {
		return nil, err
	}
	if provider.isCached {
		provider.layouts[name] = layoutTemplate
	}
	return layoutTemplate, nil
}
