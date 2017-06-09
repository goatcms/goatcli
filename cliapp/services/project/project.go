package project

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
	"github.com/goatcms/goatcore/filesystem/fsloop"
)

// Project represent project data
type Project struct {
	fs   filesystem.Filespace
	deps struct{}
}

// Factory create new repositories instance
func Factory(dp dependency.Provider) (interface{}, error) {
	var err error
	p := &Project{}
	if err = dp.InjectTo(&p.deps); err != nil {
		return nil, err
	}
	if p.fs, err = diskfs.NewFilespace("./"); err != nil {
		return nil, err
	}
	return services.Project(p), nil
}

// Filespace return project filespace
func (p *Project) Filespace() (filesystem.Filespace, error) {
	return p.fs, nil
}

// Init build repository data
func (p *Project) Init(repofs filesystem.Filespace) error {
	projectfs, err := p.Filespace()
	if err != nil {
		return err
	}
	loop := fsloop.NewLoop(&fsloop.LoopData{
		Filespace: repofs,
		DirFilter: func(fs filesystem.Filespace, subPath string) bool {
			return subPath != ".git"
		},
		OnDir: func(fs filesystem.Filespace, subPath string) error {
			return projectfs.MkdirAll(subPath, 0777)
		},
		Consumers:  1,
		Producents: 1,
	}, nil)
	loop.Run(".")
	return nil
}
