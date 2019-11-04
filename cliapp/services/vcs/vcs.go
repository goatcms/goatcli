package vcs

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
)

// VCS provide project modules data
type VCS struct {
	deps struct {
		FS filesystem.Filespace `filespace:"current"`
	}
	config []*config.Module
}

// Factory create new repositories instance
func Factory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &VCS{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	return services.VCSService(instance), nil
}

// ReadDataFromFS read vcs data from filesystem
func (vcs *VCS) ReadDataFromFS(fs filesystem.Filespace) (vcsData services.VCSData, err error) {
	var (
		ignoredFiles    *IgnoredFiles
		ignoredReader   filesystem.Reader
		generatedFiles  *GeneratedFiles
		generatedReader filesystem.Reader
	)
	if vcs.deps.FS.IsExist(IgnoredFilesPath) {
		if ignoredReader, err = vcs.deps.FS.Reader(IgnoredFilesPath); err != nil {
			return nil, err
		}
		defer ignoredReader.Close()
		if ignoredFiles, err = NewIgnoredFilesFromStream(ignoredReader); err != nil {
			return nil, err
		}
	} else {
		ignoredFiles = NewIgnoredFiles(false)
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
	return NewData(generatedFiles, ignoredFiles), nil
}

// WriteDataToFS write vcs data to filespace
func (vcs *VCS) WriteDataToFS(fs filesystem.Filespace, data services.VCSData) (err error) {
	if data.VCSIgnoredFiles().Modified() {
		if err = vcs.WriteIgnoredToFS(fs, data); err != nil {
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
func (vcs *VCS) WriteGeneratedToFS(fs filesystem.Filespace, data services.VCSData) (err error) {
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

// WriteIgnoredToFS write ignored files list to filespace
func (vcs *VCS) WriteIgnoredToFS(fs filesystem.Filespace, data services.VCSData) (err error) {
	var (
		ignoredWriter filesystem.Writer
	)
	if err = vcs.deps.FS.MkdirAll(DataDirectoryPath, filesystem.DefaultUnixDirMode); err != nil {
		return err
	}
	if vcs.deps.FS.IsFile(IgnoredFilesPath) {
		if err = vcs.deps.FS.Remove(IgnoredFilesPath); err != nil {
			return err
		}
	}
	if len(data.VCSIgnoredFiles().All()) == 0 {
		return nil
	}
	if ignoredWriter, err = vcs.deps.FS.Writer(IgnoredFilesPath); err != nil {
		return err
	}
	defer ignoredWriter.Close()
	return data.VCSIgnoredFiles().WriteAll(ignoredWriter)
}
