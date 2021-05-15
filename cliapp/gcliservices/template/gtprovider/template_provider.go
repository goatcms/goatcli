package gtprovider

import (
	"os"
	"sync"
	"text/template"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/fsloop"
)

// TemplateProvider provide templates api
type TemplateProvider struct {
	baseTemplate *template.Template
	fs           filesystem.Filespace
	path         string
	tmpl         *template.Template
	isCached     bool
	mu           sync.Mutex
}

// NewTemplateProvider create TemplateProvider instance
func NewTemplateProvider(baseTemplate *template.Template, fs filesystem.Filespace, path string, isCached bool) *TemplateProvider {
	return &TemplateProvider{
		fs:           fs,
		path:         path,
		baseTemplate: baseTemplate,
		isCached:     isCached,
	}
}

// Template load and return template
func (provider *TemplateProvider) Template() (tmpl *template.Template, err error) {
	if provider.tmpl != nil {
		return provider.tmpl, nil
	}
	return provider.load()
}

func (provider *TemplateProvider) load() (tmpl *template.Template, err error) {
	var subFS filesystem.Filespace
	provider.mu.Lock()
	defer provider.mu.Unlock()
	// check after lock
	if provider.tmpl != nil {
		return provider.tmpl, nil
	}
	// create a new view
	if tmpl, err = provider.baseTemplate.Clone(); err != nil {
		return nil, err
	}
	if subFS, err = provider.fs.Filespace(provider.path); err != nil {
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
		provider.tmpl = tmpl
	}
	return tmpl, nil
}
