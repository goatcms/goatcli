package gclimock

import (
	"github.com/goatcms/goatcore/filesystem"

	"github.com/goatcms/goatcore/app/mockupapp"
)

// NewApp create new mockup application
// it is base for other goat cli mockups
func NewApp(opt mockupapp.MockupOptions) (mapp *mockupapp.App, err error) {
	var gefs filesystem.Filespace
	if mapp, err = mockupapp.NewApp(opt); err != nil {
		return nil, err
	}
	if mapp.RootFilespace().MkdirAll(".goatcli/efs", filesystem.DefaultUnixDirMode); err != nil {
		return nil, err
	}
	if gefs, err = mapp.RootFilespace().Filespace(".goatcli/efs"); err != nil {
		return nil, err
	}
	if err = mapp.FilespaceScope().Set("gefs", gefs); err != nil {
		return nil, err
	}
	return mapp, nil
}
