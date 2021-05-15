package vcs

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// VCS provide project modules data
type VCS struct {
	deps struct {
		FS filesystem.Filespace `filespace:"current"`
	}
}

// Factory create new repositories instance
func Factory(dp app.DependencyProvider) (interface{}, error) {
	var err error
	instance := &VCS{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return gcliservices.VCSService(instance), nil
}

// ReadDataFromFS read vcs data from filesystem
func (vcs *VCS) ReadDataFromFS(fs filesystem.Filespace) (vcsData gcliservices.VCSData, err error) {
	var (
		persistedFiles  *PersistedFiles
		persistedReader filesystem.Reader
		generatedFiles  *GeneratedFiles
		generatedReader filesystem.Reader
	)
	if vcs.deps.FS.IsExist(PersistedFilesPath) {
		if persistedReader, err = vcs.deps.FS.Reader(PersistedFilesPath); err != nil {
			return nil, err
		}
		defer persistedReader.Close()
		if persistedFiles, err = NewPersistedFilesFromStream(persistedReader); err != nil {
			return nil, err
		}
	} else {
		persistedFiles = NewPersistedFiles(false)
	}
	if vcs.deps.FS.IsExist(GeneratedFilesPath) {
		if generatedReader, err = vcs.deps.FS.Reader(GeneratedFilesPath); err != nil {
			return nil, err
		}
		defer generatedReader.Close()
		if generatedFiles, err = NewGeneratedFilesFromStream(generatedReader); err != nil {
			return nil, err
		}
	} else {
		generatedFiles = NewGeneratedFiles(false)
	}
	return NewData(generatedFiles, persistedFiles), nil
}

// WriteDataToFS write vcs data to filespace
func (vcs *VCS) WriteDataToFS(fs filesystem.Filespace, data gcliservices.VCSData) (err error) {
	if data.VCSPersistedFiles().Modified() {
		if err = vcs.WritePersistedToFS(fs, data); err != nil {
			return err
		}
	}
	if data.VCSGeneratedFiles().Modified() {
		if err = vcs.WriteGeneratedToFS(fs, data); err != nil {
			return err
		}
	}
	return nil
}

// WriteGeneratedToFS write generated files list to filespace
func (vcs *VCS) WriteGeneratedToFS(fs filesystem.Filespace, data gcliservices.VCSData) (err error) {
	var (
		generatedWriter filesystem.Writer
	)
	if err = vcs.deps.FS.MkdirAll(DataDirectoryPath, filesystem.DefaultUnixDirMode); err != nil {
		return err
	}
	if vcs.deps.FS.IsFile(GeneratedFilesPath) {
		if err = vcs.deps.FS.Remove(GeneratedFilesPath); err != nil {
			return err
		}
	}
	if len(data.VCSGeneratedFiles().All()) == 0 {
		return nil
	}
	if generatedWriter, err = vcs.deps.FS.Writer(GeneratedFilesPath); err != nil {
		return err
	}
	defer generatedWriter.Close()
	return data.VCSGeneratedFiles().WriteAll(generatedWriter)
}

// WritePersistedToFS write persisted files list to filespace
func (vcs *VCS) WritePersistedToFS(fs filesystem.Filespace, data gcliservices.VCSData) (err error) {
	var (
		persistedWriter filesystem.Writer
	)
	if err = vcs.deps.FS.MkdirAll(DataDirectoryPath, filesystem.DefaultUnixDirMode); err != nil {
		return err
	}
	if vcs.deps.FS.IsFile(PersistedFilesPath) {
		if err = vcs.deps.FS.Remove(PersistedFilesPath); err != nil {
			return err
		}
	}
	if len(data.VCSPersistedFiles().All()) == 0 {
		return nil
	}
	if persistedWriter, err = vcs.deps.FS.Writer(PersistedFilesPath); err != nil {
		return err
	}
	defer persistedWriter.Close()
	return data.VCSPersistedFiles().WriteAll(persistedWriter)
}
