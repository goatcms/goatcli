package vcs

import (
	"testing"
	"time"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/varutil"
	"github.com/goatcms/goatcore/workers"
)

func TestGeneratedConcurrentlyRW(t *testing.T) {
	var (
		generated = NewGeneratedFiles(true)
		err       error
	)
	t.Parallel()
	// create directories
	for ij := workers.AsyncTestReapeat; ij > 0; ij-- {
		path := varutil.RandString(10, varutil.AlphaNumericBytes)
		generated.Add(&services.GeneratedFile{
			Path:    path,
			ModTime: time.Now(),
		})
		for i := workers.MaxJob; i > 0; i-- {
			go (func() {
				generated.Add(&services.GeneratedFile{
					Path:    path,
					ModTime: time.Now(),
				})
			})()
			go (func() {
				if err = generated.AddRow("2009-11-10T23:00:00Z;" + path); err != nil {
					t.Error(err)
					return
				}
			})()
			go (func() {
				take := generated.Get(path).Path
				if take != path {
					t.Errorf("Expected path equals to %s and take %s", path, take)
					return
				}
			})()
		}
	}
}
