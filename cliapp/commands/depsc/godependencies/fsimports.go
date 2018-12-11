package godependencies

import (
	"os"
	"strings"
	"sync"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/fsloop"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// FSDepImports return list of all golang imports in filespace.
// It find all *.go files to find imports.
func FSDepImports(sourcefs filesystem.Filespace) (imports []string, err error) {
	var mu sync.Mutex
	imports = []string{}
	loop := fsloop.NewLoop(&fsloop.LoopData{
		Filespace: sourcefs,
		DirFilter: func(fs filesystem.Filespace, subPath string) bool {
			return subPath != "./.git"
		},
		OnFile: func(fs filesystem.Filespace, subPath string) (err error) {
			var (
				newImports []string
				code       []byte
			)
			if !strings.HasSuffix(subPath, ".go") {
				return nil
			}
			if code, err = fs.ReadFile(subPath); err != nil {
				return err
			}
			if newImports, err = FindImports(string(code)); err != nil {
				return err
			}
			if len(newImports) == 0 {
				return nil
			}
			mu.Lock()
			imports = append(imports, newImports...)
			mu.Unlock()
			return nil
		},
		Consumers:  1,
		Producents: 1,
	}, nil)
	loop.Run("")
	loop.Wait()
	if len(loop.Errors()) != 0 {
		return nil, goaterr.NewErrors(loop.Errors())
	}
	return imports, nil
}

// FSDirImports return list of all golang imports in root directory.
// It find all *.go files to find imports.
func FSDirImports(sourcefs filesystem.Filespace) (imports []string, err error) {
	var (
		nodes      []os.FileInfo
		code       []byte
		newImports []string
	)
	imports = []string{}
	if nodes, err = sourcefs.ReadDir("./"); err != nil {
		return nil, err
	}
	for _, node := range nodes {
		if node.IsDir() || !strings.HasSuffix(node.Name(), ".go") {
			continue
		}
		if code, err = sourcefs.ReadFile(node.Name()); err != nil {
			return nil, err
		}
		if newImports, err = FindImports(string(code)); err != nil {
			return nil, err
		}
		imports = append(imports, newImports...)
	}
	return imports, nil
}
