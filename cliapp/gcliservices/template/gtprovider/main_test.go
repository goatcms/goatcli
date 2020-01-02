package gtprovider

const (
	masterTemplate  = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
	overlayTemplate = `{{define "list"}} {{join . ", "}}{{end}} `
	templateFile1   = `{{define "list"}} {{join . ", "}}{{end}} `
	templateFile2   = `{{define "unusedDef1"}} {{join . ": "}}{{end}} `
	templateFile3   = `{{define "unusedDef2"}} {{join . "| "}}{{end}} `
	templateFile4   = `{{define "unusedDef3"}} {{join . "/ "}}{{end}} `
)
