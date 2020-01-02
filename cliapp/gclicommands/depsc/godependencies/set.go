package godependencies

import (
	"fmt"
	"sync"

	"github.com/goatcms/goatcli/cliapp/common/config"
)

// Set is a class accumulate golang imports
type Set struct {
	resolved map[string]bool
	m        map[string]*SetRow
	mu       sync.RWMutex
}

// NewSet create new Set instance
func NewSet() *Set {
	return &Set{
		resolved: map[string]bool{},
		m:        map[string]*SetRow{},
	}
}

// Dependencies return all dependencies config data
func (set *Set) Dependencies() (result []*config.Dependency) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	result = []*config.Dependency{}
	for _, row := range set.m {
		result = append(result, row.Dependency)
	}
	return result
}

// Row return single dependency row by path
func (set *Set) Row(destPath string) (result *SetRow) {
	return set.m[destPath]
}

// Add append dependencies to accumulator
func (set *Set) Add(dependencies []*config.Dependency) (err error) {
	for _, dep := range dependencies {
		if _, ok := set.m[dep.Dest]; ok {
			return fmt.Errorf("godependencies.Add: duplicated dependencies for %s destination path", dep.Dest)
		}
		set.m[dep.Dest] = &SetRow{
			Dependency: dep,
			Imported:   false,
		}
	}
	return nil
}

// Resolve return status for path. True for resolved path or false for unresolved
func (set *Set) Resolve(gopath string) bool {
	return set.resolved[gopath]
}

// SetResolve set status for path
func (set *Set) SetResolve(gopath string, status bool) {
	set.resolved[gopath] = status
}

// AddSource add repositories master branches by URL or golang source path
/*func (set *Set) AddSource(urlOrGOPath string) (row *SetRow, err error) {
	var (
		relativePath string
	)
	if relativePath = varutil.GOPath(urlOrGOPath); relativePath == "" {
		return nil, fmt.Errorf("incorrect go path %s", urlOrGOPath)
	}
	row = &SetRow{
		Dependency: &config.Dependency{
			Repo:   "https://" + varutil.GOPath(MapPath(GOPathMapping, relativePath)) + ".git",
			Branch: "master",
			Rev:    "",
			Dest:   "vendor/" + relativePath,
		},
		Imported: false,
	}
	if _, ok := set.m[row.Dependency.Dest]; ok {
		return nil, fmt.Errorf("dependency path '%v' is already defined", urlOrGOPath)
	}
	set.m[row.Dependency.Dest] = row
	return row, nil
}
*/
