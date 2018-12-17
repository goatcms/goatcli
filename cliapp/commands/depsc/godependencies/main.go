package godependencies

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

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
