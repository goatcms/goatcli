package vcs

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewIgnoreFilesFromStream(t *testing.T) {
	var (
		instance *IgnoredFiles
		err      error
	)
	t.Parallel()
	if instance, err = NewIgnoredFilesFromStream(strings.NewReader(`
		/some/generated_file_to_ignored.go
	`)); err != nil {
		t.Error(err)
		return
	}
	if len(instance.All()) != 1 {
		t.Errorf("expected one element and take %v", len(instance.All()))
	}
	if !instance.ContainsPath("/some/generated_file_to_ignored.go") {
		t.Errorf("Should contains /some/generated_file_to_ignored.go path")
	}
	if instance.ContainsPath("/some/unknow/file.go") {
		t.Errorf("Shouldn't contains /some/unknow/file.go path")
	}
}

func TestIgnoredWrite(t *testing.T) {
	var (
		instance *IgnoredFiles
		err      error
	)
	t.Parallel()
	if instance, err = NewIgnoredFilesFromStream(strings.NewReader(`
		/some/ignored_file.go

		/some/other_ignored_file.go
	`)); err != nil {
		t.Error(err)
		return
	}
	buf := new(bytes.Buffer)
	instance.WriteAll(buf)
	expected := `/some/ignored_file.go
/some/other_ignored_file.go
`
	if buf.String() != expected {
		t.Errorf("expected %s string and take %s", expected, buf.String())
	}
}

func TestIgnoredRemove(t *testing.T) {
	var (
		instance *IgnoredFiles
		err      error
	)
	t.Parallel()
	if instance, err = NewIgnoredFilesFromStream(strings.NewReader(`
		/some/generated_file_to_ignored.go
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
		t.Errorf("Should contains only /some/generated_file_to_ignored.go Expected one eleemnt and take %d", len(instance.All()))
	}
}
