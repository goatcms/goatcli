package vcs

import "github.com/goatcms/goatcli/cliapp/services"

// Data contains vcs data like generated files list, ignored files etc
type Data struct {
	generatedFiles services.VCSGeneratedFiles
	ignoredFiles   services.VCSIgnoredFiles
}

// NewData create new data instances
func NewData(generatedFiles services.VCSGeneratedFiles, ignoredFiles services.VCSIgnoredFiles) (data *Data) {
	return &Data{
		generatedFiles: generatedFiles,
		ignoredFiles:   ignoredFiles,
	}
}

// VCSGeneratedFiles return generated file list
func (data *Data) VCSGeneratedFiles() services.VCSGeneratedFiles {
	return data.generatedFiles
}

// VCSIgnoredFiles return ignored file list
func (data *Data) VCSIgnoredFiles() services.VCSIgnoredFiles {
	return data.ignoredFiles
}
