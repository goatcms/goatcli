package godependencies

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
	"github.com/goatcms/goatcore/varutil"
	"golang.org/x/tools/go/vcs"
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
	MaxDep  int
	DevLogs bool
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
	if importer.project.relative, err = MatchGoSrcRelativePath(importer.goPath, importer.project.absolute); err != nil {
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
	var imports []string
	if imports, err = FSDepImports(importer.fs); err != nil {
		return err
	}
	return importer.processesImports(imports, 0)
}

func (importer *Importer) importDirectory(path string, dep int) (err error) {
	var (
		imports []string
		fs      filesystem.Filespace
	)
	if dep == importer.options.MaxDep {
		return fmt.Errorf("godependencies.Importer.importDirectory: To much import recursive depths. You have %v recursive import lvls", dep)
	}
	if fs, err = importer.fs.Filespace(path); err != nil {
		return err
	}
	if imports, err = FSDirImports(fs); err != nil {
		return err
	}
	return importer.processesImports(imports, dep)
}

func (importer *Importer) processesImports(imports []string, dep int) (err error) {
	imports = importer.reduceImports(imports)
	for _, importPath := range imports {
		var (
			row        *SetRow
			repoRoot   *vcs.RepoRoot
			dependency *config.Dependency
		)
		importer.logs.GOPath(importPath)
		if importer.set.Resolve(importPath) {
			continue
		}
		importer.set.SetResolve(importPath, true)
		if repoRoot, err = vcs.RepoRootForImportPath(importPath, importer.options.DevLogs); err != nil {
			return err
		}
		dependency = importer.toDependency(repoRoot)
		if row = importer.set.Row(dependency.Dest); row == nil {
			if err = importer.set.Add([]*config.Dependency{dependency}); err != nil {
				return err
			}
			if importer.logs.OnNewSource != nil {
				importer.logs.OnNewSource(dependency.Dest)
			}
			if row = importer.set.Row(dependency.Dest); row == nil {
				return fmt.Errorf("godependencies.Importer.importPath: add dependency fail")
			}
		}
		if row.Imported {
			continue
		}
		if err = importer.dependencies.CloneDependencies(importer.fs, []*config.Dependency{
			row.Dependency,
		}); err != nil {
			return err
		}
		row.Imported = true
		if err = importer.importDirectory(row.Dependency.Dest, dep+1); err != nil {
			return err
		}
	}
	return nil
}

func (importer *Importer) toDependency(repoRoot *vcs.RepoRoot) (dependency *config.Dependency) {
	dependency = &config.Dependency{
		Repo:   repoRoot.Repo,
		Branch: "master",
		Rev:    "",
		Dest:   "vendor/" + repoRoot.Root,
	}
	if !strings.Contains(dependency.Repo, "://") {
		dependency.Repo = "https://" + dependency.Repo
	}
	dependency.Repo = strings.ToLower(repoRoot.VCS.String()) + "+" + dependency.Repo
	return dependency
}

func (importer *Importer) reduceImports(imports []string) (result []string) {
	var index int
	result = []string{}
	for _, row := range imports {
		if row = varutil.FullGOPath(row); row == "" {
			continue
		}
		// skip self-import
		if strings.HasPrefix(row, importer.project.relative) {
			continue
		}
		// skip duplicates
		if varutil.IsArrContainStr(result, row) {
			continue
		}
		// skip ignored files
		if IsIgnoredPath(AlwaysIgnored, row) {
			continue
		}
		// ignore non domain url (support onlu links like domain.com/...)
		if index = strings.Index(row, "/"); index == -1 {
			continue
		}
		if strings.Index(row[:index], ".") == -1 {
			continue
		}
		result = append(result, row)
	}
	return result
}

// WriteDef save import result to definitions
func (importer *Importer) WriteDef() (err error) {
	deps := importer.set.Dependencies()
	return importer.dependencies.WriteDefToFS(importer.fs, deps)
}
