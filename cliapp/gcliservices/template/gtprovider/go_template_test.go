package gtprovider

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/goatcms/goatcore/workers"
)

const (
	master  = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
	overlay = `{{define "list"}} {{join . ", "}}{{end}} `
)

func TestGOTemplateParse(t *testing.T) {
	t.Parallel()
	var (
		funcs     = template.FuncMap{"join": strings.Join}
		guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
	)
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		masterTmpl, err := template.New("master").Funcs(funcs).Parse(master)
		if err != nil {
			t.Error(err)
			return
		}
		overlayTmpl, err := template.Must(masterTmpl.Clone()).Parse(overlay)
		if err != nil {
			t.Error(err)
			return
		}
		bufMaster := new(bytes.Buffer)
		if err := masterTmpl.Execute(bufMaster, guardians); err != nil {
			t.Error(err)
			return
		}
		result := bufMaster.String()
		if !strings.Contains(result, "- Gamora") {
			t.Errorf("after render master should contains '- Gamora' and other characters. It is: %s", result)
			return
		}
		bufOverlay := new(bytes.Buffer)
		if err := overlayTmpl.Execute(bufOverlay, guardians); err != nil {
			t.Error(err)
			return
		}
		result = bufOverlay.String()
		if !strings.Contains(result, "Gamora,") {
			t.Errorf("after render overlay should contains 'Gamora,' and other characters. It is: %s", result)
			return
		}
	}
}
