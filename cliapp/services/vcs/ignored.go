package vcs

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

// IgnoredFiles contains ignored files
type IgnoredFiles struct {
	mu      sync.RWMutex
	rows    []string
	indexes map[string]string
}

// NewIgnoredFiles create new IgnoredFiles instance
func NewIgnoredFiles() (instance *IgnoredFiles) {
	return &IgnoredFiles{
		rows:    []string{},
		indexes: map[string]string{},
	}
}

// NewIgnoredFilesFromStream create new IgnoredFiles instance and load data from stream
func NewIgnoredFilesFromStream(reader io.Reader) (instance *IgnoredFiles, err error) {
	var (
		scanner = bufio.NewScanner(reader)
	)
	instance = NewIgnoredFiles()
	for scanner.Scan() {
		instance.AddPath(scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return instance, nil
}

// All return all ignored files list
func (ignored *IgnoredFiles) All() []string {
	return ignored.rows
}

// ContainsPath return true if contains path
func (ignored *IgnoredFiles) ContainsPath(path string) bool {
	ignored.mu.RLock()
	defer ignored.mu.RUnlock()
	_, ok := ignored.indexes[path]
	return ok
}

// AddPath add new path to ignored
func (ignored *IgnoredFiles) AddPath(path string) {
	ignored.mu.Lock()
	defer ignored.mu.Unlock()
	if emptyRegexp.MatchString(path) {
		return
	}
	path = strings.Trim(path, " \t")
	if _, ok := ignored.indexes[path]; ok {
		return
	}
	ignored.rows = append(ignored.rows, path)
	ignored.indexes[path] = path
}

// WriteAll write all rows to write stream
func (ignored *IgnoredFiles) WriteAll(writer io.Writer) (err error) {
	ignored.mu.Lock()
	defer ignored.mu.Unlock()
	for _, row := range ignored.rows {
		if _, err = writer.Write([]byte(row + "\n")); err != nil {
			return err
		}
	}
	return nil
}
