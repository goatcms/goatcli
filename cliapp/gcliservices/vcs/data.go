package vcs

import "github.com/goatcms/goatcli/cliapp/gcliservices"

// Data contains vcs data like generated files list, ignored files etc
type Data struct {
	generatedFiles gcliservices.VCSGeneratedFiles
	ignoredFiles   gcliservices.VCSIgnoredFiles
}

// NewData create new data instances
func NewData(generatedFiles gcliservices.VCSGeneratedFiles, ignoredFiles gcliservices.VCSIgnoredFiles) (data *Data) {
	return &Data{
		generatedFiles: generatedFiles,
		ignoredFiles:   ignoredFiles,
	}
}

// VCSGeneratedFiles return generated file list
func (data *Data) VCSGeneratedFiles() gcliservices.VCSGeneratedFiles {
	return data.generatedFiles
}

// VCSIgnoredFiles return ignored file list
func (data *Data) VCSIgnoredFiles() gcliservices.VCSIgnoredFiles {
	return data.ignoredFiles
}
