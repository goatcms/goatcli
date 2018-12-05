package godependencies

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goatcms/goatcore/varutil"
)

// GetImportList return all project imports by path
// go list -f '{{ join .Imports "\n" }}' ./...
func GetImportList(cwd string, path, prefix string) (result []string, err error) {
	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
		coll   []string
		index  int
		gp     string
	)
	cmd := exec.Command("go", "list", "-f", "{{ join .Imports \"\\n\" }}", path)
	cmd.Dir = cwd
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err = cmd.Run(); err != nil {
		return nil, err
	}
	coll = strings.Split(outBuf.String(), "\n")
	result = []string{}
	for _, row := range coll {
		if strings.HasPrefix(row, prefix) {
			row = row[len(prefix):]
		}
		if gp = varutil.GOPath(row); gp == "" {
			continue
		}
		if varutil.IsArrContainStr(result, gp) {
			continue
		}
		index = strings.Index(row, "/")
		if strings.Index(row[:index], ".") == -1 {
			continue
		}
		result = append(result, row)
	}
	return result, nil
}

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

// MapPath replace path
func MapPath(mapping []PathMappingRow, path string) string {
	for _, row := range mapping {
		if strings.HasPrefix(path, row.From) {
			return row.To + path[len(row.From):]
		}
	}
	return path
}
