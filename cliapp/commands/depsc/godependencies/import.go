package godependencies

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
	"github.com/goatcms/goatcore/varutil"
)

// ImporterPath is path contains relative and absolute paths to the same directory
type ImporterPath struct {
	relative string
	absolute string
}

// ImporterLogs is struct contains logs callbacks
type ImporterLogs struct {
	GOPath      func(path string)
	OnNewSource func(path string)
}

// ImporterOptions is set of importer options
type ImporterOptions struct {
	MaxDep int
}

// Importer is class contains logic of import go packages
type Importer struct {
	logs         ImporterLogs
	options      ImporterOptions
	project      ImporterPath
	vendor       ImporterPath
	goPath       string
	set          *Set
	dependencies services.DependenciesService
	fs           filesystem.Filespace
}

// NewImporter create new Importer instance
func NewImporter(projectPath string, logs ImporterLogs, options ImporterOptions, dependencies services.DependenciesService) (importer *Importer, err error) {
	var (
		list []*config.Dependency
	)
	importer = &Importer{
		project: ImporterPath{
			absolute: projectPath,
		},
		vendor: ImporterPath{
			absolute: projectPath + "/vendor",
		},
		goPath:       os.Getenv("GOPATH"),
		set:          NewSet(),
		logs:         logs,
		options:      options,
		dependencies: dependencies,
	}
	if importer.project.absolute, err = filepath.Abs(importer.project.absolute); err != nil {
		return nil, err
	}
	if importer.vendor.absolute, err = filepath.Abs(importer.vendor.absolute); err != nil {
		return nil, err
	}
	if importer.project.relative, err = MatchGoSrcRelativePath(importer.goPath, projectPath); err != nil {
		return nil, err
	}
	if importer.vendor.relative, err = MatchGoSrcRelativePath(importer.goPath, importer.vendor.absolute); err != nil {
		return nil, err
	}
	if importer.fs, err = diskfs.NewFilespace(projectPath); err != nil {
		return nil, err
	}
	if list, err = dependencies.ReadDefFromFS(importer.fs); err != nil {
		return nil, err
	}
	if err = importer.set.Add(list); err != nil {
		return nil, err
	}
	return importer, nil
}

// Import load all dependencies recursively
func (importer *Importer) Import() (err error) {
	return importer.importPath(importer.fs, importer.project.absolute, 0)
}

func (importer *Importer) importPath(fs filesystem.Filespace, path string, dep int) (err error) {
	var (
		list []string
	)
	if dep == importer.options.MaxDep {
		return fmt.Errorf("To much import recursive depths. You have %v recursive import lvls", dep)
	}
	if list, err = GetImportList(path, "./...", importer.vendor.relative); err != nil {
		return err
	}
	for _, depGOPath := range list {
		var (
			row        *SetRow
			vendorPath string
		)
		depGOPath = varutil.GOPath(depGOPath)
		vendorPath = "vendor/" + depGOPath
		importer.logs.GOPath(depGOPath)
		if row = importer.set.Row(vendorPath); row == nil {
			if IsIgnoredPath(AlwaysIgnored, depGOPath) || depGOPath == importer.project.relative {
				continue
			}
			if importer.logs.OnNewSource != nil {
				importer.logs.OnNewSource(depGOPath)
			}
			if row, err = importer.set.AddSource(depGOPath); err != nil {
				return err
			}
		}
		if row.Imported {
			continue
		}
		if err = importer.dependencies.CloneDependencies(fs, []*config.Dependency{
			row.Dependency,
		}); err != nil {
			return err
		}
		row.Imported = true
		depAbsolutePath := importer.project.absolute + "/" + row.Dependency.Dest
		if err = importer.importPath(fs, depAbsolutePath, dep+1); err != nil {
			return err
		}
	}
	return nil
}

// WriteDef save import result to definitions
func (importer *Importer) WriteDef() (err error) {
	deps := importer.set.Dependencies()
	return importer.dependencies.WriteDefToFS(importer.fs, deps)
}
