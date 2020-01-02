package vcs

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
)

// NewGeneratedFile create new GeneratedFile instance
func NewGeneratedFile() (instance *gcliservices.GeneratedFile) {
	return &gcliservices.GeneratedFile{}
}

// NewGeneratedFileFromText create new GeneratedFile instance from text
func NewGeneratedFileFromText(text string) (instance *gcliservices.GeneratedFile, err error) {
	var parts = strings.SplitN(text, ";", 2)
	instance = NewGeneratedFile()
	if instance.ModTime, err = time.Parse(time.RFC3339, parts[0]); err != nil {
		return instance, err
	}
	instance.Path = parts[1]
	instance.Path = strings.Trim(instance.Path, " \t")
	return instance, nil
}

// GeneratedFiles contains generated files
type GeneratedFiles struct {
	mu       sync.RWMutex
	rows     []*gcliservices.GeneratedFile
	indexes  map[string]*gcliservices.GeneratedFile
	news     map[string]*gcliservices.GeneratedFile
	modified bool
}

// NewGeneratedFiles create new GeneratedFiles instance
func NewGeneratedFiles(modified bool) (instance *GeneratedFiles) {
	return &GeneratedFiles{
		rows:     []*gcliservices.GeneratedFile{},
		indexes:  map[string]*gcliservices.GeneratedFile{},
		news:     map[string]*gcliservices.GeneratedFile{},
		modified: modified,
	}
}

// NewGeneratedFilesFromStream create new GeneratedFiles instance and load data from stream
func NewGeneratedFilesFromStream(reader io.Reader) (instance *GeneratedFiles, err error) {
	var (
		scanner = bufio.NewScanner(reader)
	)
	instance = NewGeneratedFiles(false)
	for scanner.Scan() {
		instance.addRow(scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return instance, nil
}

// All return all generated files list
func (generated *GeneratedFiles) All() []*gcliservices.GeneratedFile {
	sort.Slice(generated.rows, func(i, j int) bool {
		return generated.rows[i].Path < generated.rows[j].Path
	})
	return generated.rows
}

// New return all new rows
func (generated *GeneratedFiles) New() map[string]*gcliservices.GeneratedFile {
	return generated.news
}

// Get return generated file description or nil
func (generated *GeneratedFiles) Get(path string) *gcliservices.GeneratedFile {
	generated.mu.RLock()
	defer generated.mu.RUnlock()
	return generated.indexes[path]
}

// ContainsPath return true if contains path
func (generated *GeneratedFiles) ContainsPath(path string) bool {
	generated.mu.RLock()
	defer generated.mu.RUnlock()
	_, ok := generated.indexes[path]
	return ok
}

// AddRow add new row from text to generated
func (generated *GeneratedFiles) AddRow(text string) (err error) {
	var row *gcliservices.GeneratedFile
	generated.mu.Lock()
	defer generated.mu.Unlock()
	if emptyRegexp.MatchString(text) {
		return nil
	}
	if row, err = generated.addRow(text); err != nil {
		return err
	}
	generated.news[row.Path] = row
	return nil
}

func (generated *GeneratedFiles) addRow(text string) (row *gcliservices.GeneratedFile, err error) {
	if emptyRegexp.MatchString(text) {
		return nil, nil
	}
	text = strings.Trim(text, " \t")
	if row, err = NewGeneratedFileFromText(text); err != nil {
		return nil, err
	}
	generated.add(row)
	generated.modified = true
	return row, nil
}

// Add create or update row in collection
func (generated *GeneratedFiles) Add(row *gcliservices.GeneratedFile) {
	generated.mu.Lock()
	defer generated.mu.Unlock()
	generated.add(row)
	generated.news[row.Path] = row
}

func (generated *GeneratedFiles) add(row *gcliservices.GeneratedFile) {
	if current, ok := generated.indexes[row.Path]; ok {
		current.ModTime = row.ModTime
		return
	}
	generated.rows = append(generated.rows, row)
	generated.indexes[row.Path] = row
	generated.modified = true
}

// WriteAll write all rows to write stream
func (generated *GeneratedFiles) WriteAll(writer io.Writer) (err error) {
	generated.mu.RLock()
	defer generated.mu.RUnlock()
	for _, row := range generated.rows {
		if _, err = writer.Write([]byte(fmt.Sprintf("%s;%s\n", row.ModTime.Format(time.RFC3339), row.Path))); err != nil {
			return err
		}
	}
	return nil
}

// Modified return true if object was modiefid
func (generated *GeneratedFiles) Modified() bool {
	return generated.modified
}
