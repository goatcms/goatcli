package godependencies

import (
	"fmt"
	"testing"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/memfs"
	"github.com/goatcms/goatcore/varutil"
	"github.com/goatcms/goatcore/workers"
)

const (
	firstGoFile = `
	package test

	import "github.com/goatcms/goatcore"

	import (
		"github.com/goatcms/goatcms"
		"github.com/goatcms/goatcli"
	)`
	secondGoFile = `
		package test

		import "google.golang.org/grpc"

		import (
			"google.golang.org/appengine"
			"google.golang.org/genproto"
		)`
)

func TestFSDepImports(t *testing.T) {
	var (
		fs     filesystem.Filespace
		result []string
		err    error
	)
	if fs, err = memfs.NewFilespace(); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile("path/to/example/file.go", []byte(firstGoFile), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile("path/to/test/file.go", []byte(secondGoFile), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < workers.AsyncTestReapeat; i++ {
		if result, err = FSDepImports(fs); err != nil {
			t.Error(err)
			return
		}
		//
		if len(result) != 6 {
			fmt.Errorf("Expected 6 results")
		}
		if !varutil.IsArrContainStr(result, "github.com/goatcms/goatcore") {
			fmt.Errorf("Expected github.com/goatcms/goatcore")
			return
		}
		if !varutil.IsArrContainStr(result, "github.com/goatcms/goatcms") {
			fmt.Errorf("Expected github.com/goatcms/goatcms")
			return
		}
		if !varutil.IsArrContainStr(result, "github.com/goatcms/goatcli") {
			fmt.Errorf("Expected github.com/goatcms/goatcli")
			return
		}
		if !varutil.IsArrContainStr(result, "google.golang.org/grpc") {
			fmt.Errorf("Expected google.golang.org/grpc")
			return
		}
		if !varutil.IsArrContainStr(result, "google.golang.org/appengine") {
			fmt.Errorf("Expected google.golang.org/appengine")
			return
		}
		if !varutil.IsArrContainStr(result, "google.golang.org/genproto") {
			fmt.Errorf("Expected google.golang.org/genproto")
			return
		}
	}
}

func TestFSDirImports(t *testing.T) {
	var (
		fs     filesystem.Filespace
		result []string
		err    error
	)
	if fs, err = memfs.NewFilespace(); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile("path/to/example/file.go", []byte(firstGoFile), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile("file.go", []byte(secondGoFile), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	if result, err = FSDirImports(fs); err != nil {
		t.Error(err)
		return
	}
	if len(result) != 3 {
		fmt.Errorf("Expected 3 results")
	}
	if !varutil.IsArrContainStr(result, "google.golang.org/grpc") {
		fmt.Errorf("Expected google.golang.org/grpc")
		return
	}
	if !varutil.IsArrContainStr(result, "google.golang.org/appengine") {
		fmt.Errorf("Expected google.golang.org/appengine")
		return
	}
	if !varutil.IsArrContainStr(result, "google.golang.org/genproto") {
		fmt.Errorf("Expected google.golang.org/genproto")
		return
	}
}
