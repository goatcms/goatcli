package cloner

import (
	"github.com/goatcms/goatcli/cliapp/common/mockups"
	"github.com/goatcms/goatcli/cliapp/common/result"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/memfs"
)

const (
	testReplacesJSON = `[{"from":"\{\{project_name\}\}", "to":"{{property_project_name}}", "pattern":"[A-Za-z0-9_/]*.(md|txt)"}]`
	testModulesJSON  = `[{"srcClone":"https://github.com/goatcms/mockupmodule", "srcBranch":"master", "srcDir":"module"}]`
)

func buildSrcFilespace() (fs filesystem.Filespace, err error) {
	if fs, err = memfs.NewFilespace(); err != nil {
		return nil, err
	}
	if err = fs.WriteFile(".git/noCopyGitDir.md", []byte(""), 0766); err != nil {
		return nil, err
	}
	if err = fs.WriteFile(".goat/replaces.def.json", []byte(testReplacesJSON), 0766); err != nil {
		return nil, err
	}
	if err = fs.WriteFile(".goat/modules.def.json", []byte(testModulesJSON), 0766); err != nil {
		return nil, err
	}
	if err = fs.WriteFile("main.go", []byte("package main\n/*Main package*/\n"), 0766); err != nil {
		return nil, err
	}
	if err = fs.WriteFile("docs/main.md", []byte("Description your {{project_name}}"), 0766); err != nil {
		return nil, err
	}
	if err = fs.WriteFile("docs/main.txt", []byte("No process this file {{project_name}}"), 0766); err != nil {
		return nil, err
	}
	return fs, nil
}

func buildModuleSrcFilespace() (fs filesystem.Filespace, err error) {
	if fs, err = memfs.NewFilespace(); err != nil {
		return nil, err
	}
	if err = fs.WriteFile(".git/noCopyGitDir.md", []byte(""), 0766); err != nil {
		return nil, err
	}
	if err = fs.WriteFile("module.go", []byte("package main\n/*Main package*/\n"), 0766); err != nil {
		return nil, err
	}
	return fs, nil
}

func buildDestFilespace() (fs filesystem.Filespace, err error) {
	return memfs.NewFilespace()
}

func buildPropertiesResult() *result.PropertiesResult {
	return result.NewPropertiesResult(map[string]string{
		"property_project_name": "my_project",
	})
}

func buildRepositoriesService() (services.RepositoriesService, error) {
	var (
		moduleFS filesystem.Filespace
		rootFS   filesystem.Filespace
		err      error
	)
	if rootFS, err = buildSrcFilespace(); err != nil {
		return nil, err
	}
	if moduleFS, err = buildModuleSrcFilespace(); err != nil {
		return nil, err
	}
	return mockups.NewRepositoriesService(map[string]filesystem.Filespace{
		"https://github.com/goatcms/mockup.master.":       rootFS,
		"https://github.com/goatcms/mockupmodule.master.": moduleFS,
	}), nil
}
