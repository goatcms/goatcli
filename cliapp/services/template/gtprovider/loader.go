package gtprovider

import (
	"strings"
	"sync"
	"text/template"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// TemplateLoader provide method to load templates from filesystem
type TemplateLoader struct {
	template   *template.Template
	muTemplate sync.Mutex
}

// NewTemplateLoader create TemplateLoader instance
func NewTemplateLoader(template *template.Template) *TemplateLoader {
	return &TemplateLoader{
		template: template,
	}
}

// LoadByExtension file by extension (if file has unknown extension it is skipped and return nil, false)
func (loader *TemplateLoader) LoadByExtension(fs filesystem.Filespace, subPath string) (loaded bool, err error) {
	if strings.HasSuffix(subPath, TemplateExtension) {
		return true, loader.LoadFullTemplate(fs, subPath)
	} else if strings.HasSuffix(subPath, OnceTemplateExtension) {
		return true, loader.LoadSingleFileTemplate(fs, subPath)
	} else if strings.HasSuffix(subPath, RenderTemplateExtension) {
		return true, loader.LoadSingleFileTemplate(fs, subPath)
	} else if strings.HasSuffix(subPath, CtrlTemplateExtension) {
		return true, loader.LoadSingleFileTemplate(fs, subPath)
	} else if strings.HasSuffix(subPath, DefTemplateExtension) {
		return true, loader.LoadSingleFileTemplate(fs, subPath)
	}
	return false, nil
}

// LoadFullTemplate get all templates code form files in subPath and add it to template
func (loader *TemplateLoader) LoadFullTemplate(fs filesystem.Filespace, subPath string) error {
	bytes, err := fs.ReadFile(subPath)
	if err != nil {
		return err
	}
	loader.muTemplate.Lock()
	defer loader.muTemplate.Unlock()
	if len(bytes) == 0 {
		return goaterr.Errorf("empty file")
	}
	if _, err := loader.template.Parse(string(bytes)); err != nil {
		return goaterr.Errorf("%v: %v", subPath, err)
	}
	return nil
}

// LoadSingleFileTemplate load file contins single template
func (loader *TemplateLoader) LoadSingleFileTemplate(fs filesystem.Filespace, subPath string) error {
	bytes, err := fs.ReadFile(subPath)
	if err != nil {
		return err
	}
	loader.muTemplate.Lock()
	defer loader.muTemplate.Unlock()
	fullTemplateCode := `{{define "` + subPath + `"}}{{ $ctx := . }}` + string(bytes) + `{{end}}`
	if _, err := loader.template.Parse(fullTemplateCode); err != nil {
		return goaterr.Errorf("%v: %v", subPath, err)
	}
	return nil
}

// Template return loaded template
func (loader *TemplateLoader) Template() *template.Template {
	return loader.template
}
