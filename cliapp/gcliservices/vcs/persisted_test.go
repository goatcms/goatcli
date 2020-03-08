package vcs

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewIgnoreFilesFromStream(t *testing.T) {
	var (
		instance *PersistedFiles
		err      error
	)
	t.Parallel()
	if instance, err = NewPersistedFilesFromStream(strings.NewReader(`
		/some/generated_file_to_persisted.go
	`)); err != nil {
		t.Error(err)
		return
	}
	if len(instance.All()) != 1 {
		t.Errorf("expected one element and take %v", len(instance.All()))
	}
	if !instance.ContainsPath("/some/generated_file_to_persisted.go") {
		t.Errorf("Should contains /some/generated_file_to_persisted.go path")
	}
	if instance.ContainsPath("/some/unknow/file.go") {
		t.Errorf("Shouldn't contains /some/unknow/file.go path")
	}
}

func TestPersistedWrite(t *testing.T) {
	var (
		instance *PersistedFiles
		err      error
	)
	t.Parallel()
	if instance, err = NewPersistedFilesFromStream(strings.NewReader(`
		/some/persisted_file.go

		/some/other_persisted_file.go
	`)); err != nil {
		t.Error(err)
		return
	}
	buf := new(bytes.Buffer)
	instance.WriteAll(buf)
	expected := `/some/persisted_file.go
/some/other_persisted_file.go
`
	if buf.String() != expected {
		t.Errorf("expected %s string and take %s", expected, buf.String())
	}
}

func TestPersistedRemove(t *testing.T) {
	var (
		instance *PersistedFiles
		err      error
	)
	t.Parallel()
	if instance, err = NewPersistedFilesFromStream(strings.NewReader(`
		/some/generated_file_to_persisted.go
		/some/generated_file_to_remove.go
	`)); err != nil {
		t.Error(err)
		return
	}
	if instance.Modified() != false {
		t.Errorf("Generated files should be unmodified after read")
		return
	}
	instance.RemovePath("/some/generated_file_to_remove.go")
	if instance.Modified() != true {
		t.Errorf("Generated files should be modified after remove path")
		return
	}
	if len(instance.All()) != 1 {
		t.Errorf("Should contains only /some/generated_file_to_persisted.go Expected one eleemnt and take %d", len(instance.All()))
	}
}
