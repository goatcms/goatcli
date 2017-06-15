package cloner

import (
	"io"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/fsloop"
)

// Cloner clone and process new project
type Cloner struct {
	deps struct {
		Repositories services.Repositories `dependency:"RepositoriesService"`
	}
}

// Factory create new repositories instance
func Factory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &Cloner{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return services.Cloner(instance), nil
}

// Clone clone repository
func (cloner *Cloner) Clone(repository, rev string, destfs filesystem.Filespace) (err error) {
	var sourcefs filesystem.Filespace
	if sourcefs, err = cloner.deps.Repositories.Filespace(repository, rev); err != nil {
		return err
	}
	loop := fsloop.NewLoop(&fsloop.LoopData{
		Filespace: sourcefs,
		DirFilter: func(fs filesystem.Filespace, subPath string) bool {
			return subPath != ".git"
		},
		OnDir: func(fs filesystem.Filespace, subPath string) error {
			return destfs.MkdirAll(subPath, 0777)
		},
		OnFile: func(fs filesystem.Filespace, subPath string) error {
			var (
				reader io.Reader
				writer io.Writer
				err    error
			)
			if reader, err = sourcefs.Reader(subPath); err != nil {
				return err
			}
			if writer, err = destfs.Writer(subPath); err != nil {
				return err
			}
			if _, err = io.Copy(writer, reader); err != nil {
				return err
			}
			return nil
		},
		Consumers:  1,
		Producents: 1,
	}, nil)
	loop.Run(".")
	return nil
}
