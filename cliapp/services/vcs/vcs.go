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
		FS filesystem.Filespace `filespace:"root"`
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
		ignoredFiles = NewIgnoredFiles()
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
		generatedFiles = NewGeneratedFiles()
	}
	return NewData(generatedFiles, ignoredFiles), nil
}

// WriteDataToFS write vcs data to filespace
func (vcs *VCS) WriteDataToFS(fs filesystem.Filespace, data services.VCSData) (err error) {
	var (
		ignoredWriter   filesystem.Writer
		generatedWriter filesystem.Writer
	)
	if generatedWriter, err = vcs.deps.FS.Writer(GeneratedFilesPath); err != nil {
		return err
	}
	defer generatedWriter.Close()
	if err = data.VCSGeneratedFiles().WriteAll(generatedWriter); err != nil {
		return err
	}
	if ignoredWriter, err = vcs.deps.FS.Writer(IgnoredFilesPath); err != nil {
		return err
	}
	defer ignoredWriter.Close()
	return data.VCSIgnoredFiles().WriteAll(ignoredWriter)
}
