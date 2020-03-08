package vcs

import "github.com/goatcms/goatcli/cliapp/gcliservices"

// Data contains vcs data like generated files list, persisted files etc
type Data struct {
	generatedFiles gcliservices.VCSGeneratedFiles
	persistedFiles gcliservices.VCSPersistedFiles
}

// NewData create new data instances
func NewData(generatedFiles gcliservices.VCSGeneratedFiles, persistedFiles gcliservices.VCSPersistedFiles) (data *Data) {
	return &Data{
		generatedFiles: generatedFiles,
		persistedFiles: persistedFiles,
	}
}

// VCSGeneratedFiles return generated file list
func (data *Data) VCSGeneratedFiles() gcliservices.VCSGeneratedFiles {
	return data.generatedFiles
}

// VCSPersistedFiles return persisted file list
func (data *Data) VCSPersistedFiles() gcliservices.VCSPersistedFiles {
	return data.persistedFiles
}
