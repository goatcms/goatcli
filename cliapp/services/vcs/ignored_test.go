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
