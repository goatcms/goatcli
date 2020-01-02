package vcs

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
)

func TestNewGeneratedFilesFromStream(t *testing.T) {
	var (
		expectedTime time.Time
		newTime      time.Time
		instance     *GeneratedFiles
		row          *gcliservices.GeneratedFile
		err          error
	)
	t.Parallel()
	if instance, err = NewGeneratedFilesFromStream(strings.NewReader(`
		2009-11-10T23:00:00Z;/some/generated_file.go
		2009-11-10T23:00:00Z;/some/other_generated_file.go
	`)); err != nil {
		t.Error(err)
		return
	}
	if len(instance.All()) != 2 {
		t.Errorf("Expect two generated files and take %d", len(instance.All()))
	}
	if !instance.ContainsPath("/some/other_generated_file.go") {
		t.Errorf("Should contains /some/other_generated_file.go path")
	}
	if !instance.ContainsPath("/some/generated_file.go") {
		t.Errorf("Shouldn't contains /some/generated_file.go path")
	}
	if instance.ContainsPath("/some/unknow_file.go") {
		t.Errorf("Shouldn't contains /some/unknow_file.go path")
	}
	if expectedTime, err = time.Parse(time.RFC3339, "2009-11-10T23:00:00Z"); err != nil {
		t.Error(err)
		return
	}
	row = instance.Get("/some/generated_file.go")
	if row.ModTime != expectedTime {
		t.Errorf("Expected time equals to 2009-11-10T23:00:00Z")
	}
	if newTime, err = time.Parse(time.RFC3339, "2012-11-10T23:00:00Z"); err != nil {
		t.Error(err)
		return
	}
	instance.Add(&gcliservices.GeneratedFile{
		Path:    "/some/generated_file.go",
		ModTime: newTime,
	})
	if instance.Get("/some/generated_file.go").ModTime != newTime {
		t.Errorf("Expected time after update equals to 2012-11-10T23:00:00Z")
	}
}

func TestGeneratedFilesGelAllOrder(t *testing.T) {
	var (
		instance *GeneratedFiles
		err      error
	)
	t.Parallel()
	if instance, err = NewGeneratedFilesFromStream(strings.NewReader(`
		2009-11-10T23:00:00Z;/some/b.go
		2009-11-10T23:00:00Z;/some/a.go
		2009-11-10T23:00:00Z;/some/c.go
	`)); err != nil {
		t.Error(err)
		return
	}
	result := instance.All()
	if len(result) != 3 {
		t.Errorf("Expect two generated files and take %d", len(instance.All()))
	}
	if result[0].Path != "/some/a.go" {
		t.Errorf("Expected /some/a.go as first")
	}
	if result[1].Path != "/some/b.go" {
		t.Errorf("Expected /some/b.go as second")
	}
	if result[2].Path != "/some/c.go" {
		t.Errorf("Expected /some/c.go as third")
	}
}

func TestGeneratedWrite(t *testing.T) {
	var (
		instance *GeneratedFiles
		err      error
	)
	t.Parallel()
	if instance, err = NewGeneratedFilesFromStream(strings.NewReader(`
		2009-11-10T23:00:00Z;/some/generated_file.go

		2009-11-10T23:00:00Z;/some/other_generated_file.go
	`)); err != nil {
		t.Error(err)
		return
	}
	buf := new(bytes.Buffer)
	instance.WriteAll(buf)
	expected := `2009-11-10T23:00:00Z;/some/generated_file.go
2009-11-10T23:00:00Z;/some/other_generated_file.go
`
	if buf.String() != expected {
		t.Errorf("expected %s string and take %s", expected, buf.String())
	}
}
