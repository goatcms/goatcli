package cloner

import (
	"github.com/goatcms/goatcli/cliapp/common/mockups"
	"github.com/goatcms/goatcli/cliapp/common/result"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/memfs"
)

func buildSrcFilespace() (fs filesystem.Filespace, err error) {
	if fs, err = memfs.NewFilespace(); err != nil {
		return nil, err
	}
	if err = fs.WriteFile(".git/noCopyGitDir.md", []byte(""), 0766); err != nil {
		return nil, err
	}
	if err = fs.WriteFile(".goat/replace.json", []byte(`[{"from":"\{\{project_name\}\}", "to":"{{property_project_name}}", "pattern":"[A-Za-z0-9_/]*.(md|txt)"}]`), 0766); err != nil {
		return nil, err
	}
	if err = fs.WriteFile("main.go", []byte("package main\n/*Main package*/\n"), 0766); err != nil {
		return nil, err
	}
	if err = fs.WriteFile("docs/main.md", []byte("Description your {{project_name}}"), 0766); err != nil {
		return nil, err
	}
	if err = fs.WriteFile("docs/main.txt", []byte("No proccess this file {{project_name}}"), 0766); err != nil {
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

func buildRepositoriesService() (services.Repositories, error) {
	var (
		fs  filesystem.Filespace
		err error
	)
	if fs, err = buildSrcFilespace(); err != nil {
		return nil, err
	}
	return mockups.NewRepositoriesService(map[string]filesystem.Filespace{
		"https://github.com/goatcms/mockup!master": fs,
	}), nil
}
