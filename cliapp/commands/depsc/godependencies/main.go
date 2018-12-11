package godependencies

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// GetExternalImportList return all external project imports by path
// go list -f '{{ join .Imports "\n" }}' ./...
/*func GetExternalImportList(path string) (result []string, err error) {
	var (
		fs      filesystem.Filespace
		imports []string
		index   int
	)
	if fs, err = diskfs.NewFilespace(path); err != nil {
		return nil, err
	}
	if imports, err = FSImports(fs); err != nil {
		return nil, err
	}
	result = []string{}
	for _, row := range imports {
		if row = varutil.FullGOPath(row); row == "" {
			continue
		}
		if IsIgnoredPath(AlwaysIgnored, row) {
			continue
		}
		if varutil.IsArrContainStr(result, row) {
			continue
		}
		if index = strings.Index(row, "/"); index == -1 {
			continue
		}
		if strings.Index(row[:index], ".") == -1 {
			continue
		}
		result = append(result, row)
	}
	return result, nil
}*/

// MatchGoSrcRelativePath match path relative to GOPATH
func MatchGoSrcRelativePath(GOPath, cwd string) (result string, err error) {
	var (
		cwdabs        string
		packagePrefix string
	)
	GOSourcePath := GOPath + "/src"
	if cwdabs, err = filepath.Abs(cwd); err != nil {
		return "", err
	}
	if !strings.HasPrefix(cwdabs, GOSourcePath) {
		return "", fmt.Errorf("Current working directory must be in golang source directory ($GOPATH/src) ")
	}
	packagePrefix = cwdabs[len(GOSourcePath):]
	if strings.HasPrefix(packagePrefix, "/") {
		packagePrefix = packagePrefix[1:]
	}
	return packagePrefix, nil
}

// IsIgnoredPath return true if path is ignored
func IsIgnoredPath(coll []*regexp.Regexp, path string) bool {
	for _, row := range coll {
		if row.MatchString(path) {
			return true
		}
	}
	return false
}
