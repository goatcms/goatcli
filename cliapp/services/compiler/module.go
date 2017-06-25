package compiler

import "github.com/goatcms/goatcore/dependency"

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp dependency.Provider) error {
	if err := dp.AddDefaultFactory("CompilerService", CompilerFactory); err != nil {
		return err
	}
	if err := dp.AddDefaultFactory("FileCompilerService", FileCompilerFactory); err != nil {
		return err
	}
	if err := dp.AddDefaultFactory("TemplateCompilerService", TemplateCompilerFactory); err != nil {
		return err
	}
	return nil
}
