package godependencies

import (
	"testing"
)

const (
	testInlineImportsCode = `
package godependencies
//import "github.com/goatcms/error"
` + "var asd = `\nimport (\n\"github.com/goatcms/error3\"\n`" + `
/*
import "github.com/goatcms/error2"
*/

		import "github.com/goatcms/goatcms"
import "github.com/goatcms/goatcli"
		import "github.com/goatcms/goatcore"`

	testMultilineImportsCode = `
		package godependencies
				import (
					"github.com/goatcms/goatcms"
					"github.com/goatcms/goatcli"
				)
				//import ( "github.com/goatcms/error" )
				/*
				import (
					"github.com/goatcms/error2"
				)
				*/
				import (
					"github.com/goatcms/goatcore"
				)`

	testMixImportsCode = `
				package godependencies
						import "github.com/goatcms/goatcms"
						import (
							"github.com/goatcms/goatcore"
						)`
)

func TestFindInlineImports(t *testing.T) {
	t.Parallel()
	var (
		result []string
		err    error
	)
	if result, err = FindImports(testInlineImportsCode); err != nil {
		t.Error(err)
		return
	}
	if len(result) != 3 {
		t.Errorf("Expected 3 results and take %v", len(result))
		return
	}
	if result[0] != "github.com/goatcms/goatcms" {
		t.Errorf("Expected github.com/goatcms/goatcms and take %v", result[0])
	}
	if result[1] != "github.com/goatcms/goatcli" {
		t.Errorf("Expected github.com/goatcms/goatcli and take %v", result[1])
	}
	if result[2] != "github.com/goatcms/goatcore" {
		t.Errorf("Expected github.com/goatcms/goatcore and take %v", result[2])
	}
}

func TestMultilineImportsCode(t *testing.T) {
	t.Parallel()
	var (
		result []string
		err    error
	)
	if result, err = FindImports(testMultilineImportsCode); err != nil {
		t.Error(err)
		return
	}
	if len(result) != 3 {
		t.Errorf("Expected 3 results and take %v", len(result))
		return
	}
	if result[0] != "github.com/goatcms/goatcms" {
		t.Errorf("Expected github.com/goatcms/goatcms and take %v", result[0])
	}
	if result[1] != "github.com/goatcms/goatcli" {
		t.Errorf("Expected github.com/goatcms/goatcli and take %v", result[1])
	}
	if result[2] != "github.com/goatcms/goatcore" {
		t.Errorf("Expected github.com/goatcms/goatcore and take %v", result[2])
	}
}

func TestFindImports(t *testing.T) {
	t.Parallel()
	var (
		expected string
		result   []string
		err      error
	)
	if result, err = FindImports(testMixImportsCode); err != nil {
		t.Error(err)
		return
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 results and take %v", result)
		return
	}
	expected = "github.com/goatcms/goatcms"
	if result[0] != expected && result[1] != expected {
		t.Errorf("Expected github.com/goatcms/goatcms")
	}
	expected = "github.com/goatcms/goatcore"
	if result[0] != expected && result[1] != expected {
		t.Errorf("Expected github.com/goatcms/goatcore ")
	}
}

func TestRemoveComments(t *testing.T) {
	var (
		result string
	)
	t.Parallel()
	if result = removeComments(`abc/*a\nsd\nsad*/def`); result != "abcdef" {
		t.Errorf("Expected abcdef and take %v", result)
	}
	if result = removeComments(`abc// dome comment text\n`); result != "abc\n" {
		t.Errorf("Expected abc and take %v", result)
	}
}

func TestReduceMultilineStrings(t *testing.T) {
	var (
		result string
	)
	t.Parallel()
	if result = reduceMultilineStrings("`test` `asd\nasd\n`"); result != "`` ``" {
		t.Errorf("Expected two empty multiline strings and take %v", result)
	}
}
