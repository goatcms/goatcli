package vcs

import "github.com/goatcms/goatcli/cliapp/services"

// Data contains vcs data like generated files list, ignored files etc
type Data struct {
	generatedFiles *GeneratedFiles
	ignoredFiles   *IgnoredFiles
}

// NewData create new data instances
func NewData(generatedFiles *GeneratedFiles, ignoredFiles *IgnoredFiles) (data *Data) {
	return &Data{
		generatedFiles: generatedFiles,
		ignoredFiles:   ignoredFiles,
	}
}

// GeneratedFiles return generated file list
func (data *Data) GeneratedFiles() *GeneratedFiles {
	return data.generatedFiles
}

// IgnoredFiles return ignored file list
func (data *Data) IgnoredFiles() *IgnoredFiles {
	return data.ignoredFiles
}

// VCSGeneratedFiles return generated file list
func (data *Data) VCSGeneratedFiles() services.VCSGeneratedFiles {
	return data.generatedFiles
}

// VCSIgnoredFiles return ignored file list
func (data *Data) VCSIgnoredFiles() services.VCSIgnoredFiles {
	return data.ignoredFiles
}
