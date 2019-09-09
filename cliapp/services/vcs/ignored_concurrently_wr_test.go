package vcs

import (
	"testing"

	"github.com/goatcms/goatcore/varutil"
	"github.com/goatcms/goatcore/workers"
)

func TestIgnoredConcurrentlyRead(t *testing.T) {
	var (
		ignored = NewIgnoredFiles(true)
	)
	t.Parallel()
	// create directories
	for ij := workers.AsyncTestReapeat; ij > 0; ij-- {
		path := varutil.RandString(10, varutil.AlphaNumericBytes)
		ignored.AddPath(path)
		for i := workers.MaxJob; i > 0; i-- {
			go (func() {
				ignored.AddPath(path)
			})()
			go (func() {
				if !ignored.ContainsPath(path) {
					t.Errorf("Should contains path")
					return
				}
			})()
		}
	}
}
