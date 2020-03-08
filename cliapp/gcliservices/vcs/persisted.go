package vcs

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

// PersistedFiles contains persisted files
type PersistedFiles struct {
	mu       sync.RWMutex
	rows     []string
	indexes  map[string]string
	modified bool
}

// NewPersistedFiles create new PersistedFiles instance
func NewPersistedFiles(modified bool) (instance *PersistedFiles) {
	return &PersistedFiles{
		rows:     []string{},
		indexes:  map[string]string{},
		modified: modified,
	}
}

// NewPersistedFilesFromStream create new PersistedFiles instance and load data from stream
func NewPersistedFilesFromStream(reader io.Reader) (instance *PersistedFiles, err error) {
	var (
		scanner = bufio.NewScanner(reader)
	)
	instance = NewPersistedFiles(false)
	for scanner.Scan() {
		instance.AddPath(scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	instance.modified = false
	return instance, nil
}

// All return all persisted files list
func (persisted *PersistedFiles) All() []string {
	return persisted.rows
}

// ContainsPath return true if contains path
func (persisted *PersistedFiles) ContainsPath(path string) bool {
	persisted.mu.RLock()
	defer persisted.mu.RUnlock()
	_, ok := persisted.indexes[path]
	return ok
}

// AddPath add new path to persisted
func (persisted *PersistedFiles) AddPath(path string) {
	persisted.mu.Lock()
	defer persisted.mu.Unlock()
	if emptyRegexp.MatchString(path) {
		return
	}
	path = strings.Trim(path, " \t")
	if path == "" {
		return
	}
	if _, ok := persisted.indexes[path]; ok {
		return
	}
	persisted.rows = append(persisted.rows, path)
	persisted.indexes[path] = path
	persisted.modified = true
}

// RemovePath add new path to persisted
func (persisted *PersistedFiles) RemovePath(path string) {
	persisted.mu.Lock()
	defer persisted.mu.Unlock()
	if emptyRegexp.MatchString(path) {
		return
	}
	path = strings.Trim(path, " \t")
	if _, ok := persisted.indexes[path]; !ok {
		return
	}
	for i, row := range persisted.rows {
		if row == path {
			persisted.rows = append(persisted.rows[:i], persisted.rows[i+1:]...)
			break
		}
	}
	delete(persisted.indexes, path)
	persisted.modified = true
}

// WriteAll write all rows to write stream
func (persisted *PersistedFiles) WriteAll(writer io.Writer) (err error) {
	persisted.mu.Lock()
	defer persisted.mu.Unlock()
	for _, row := range persisted.rows {
		if _, err = writer.Write([]byte(row + "\n")); err != nil {
			return err
		}
	}
	return nil
}

// Modified return true if object was modiefid
func (persisted *PersistedFiles) Modified() bool {
	return persisted.modified
}
