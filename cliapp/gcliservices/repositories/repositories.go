package repositories

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/disk"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
	"github.com/goatcms/goatcore/repositories"
)

// Repositories provide access to repositories
type Repositories struct {
	deps struct {
		Connector repositories.Connector `dependency:"RepositoriesConnector"`
	}
}

// Factory create new repositories instance
func Factory(dp app.DependencyProvider) (interface{}, error) {
	r := &Repositories{}
	if err := dp.InjectTo(&r.deps); err != nil {
		return nil, err
	}
	return gcliservices.RepositoriesService(r), nil
}

// Filespace return filespace for repository
func (repos *Repositories) Filespace(repoURL string, version repositories.Version) (filesystem.Filespace, error) {
	var (
		destPath string
	)
	basepath, err := repos.srcpath()
	if err != nil {
		return nil, err
	}
	if version.Branch != "" && version.Revision != "" {
		destPath = basepath + reduceRepoURL(repoURL) + "." + version.Branch + "." + version.Revision
	} else if version.Revision != "" {
		destPath = basepath + reduceRepoURL(repoURL) + ".master." + version.Revision
	} else if version.Branch != "" {
		destPath = basepath + reduceRepoURL(repoURL) + "." + version.Branch
	} else {
		destPath = basepath + reduceRepoURL(repoURL)
	}
	if disk.IsDir(destPath) {
		if err = repos.update(destPath); err != nil {
			return nil, err
		}
	} else {
		os.RemoveAll(destPath)
		if err = repos.clone(repoURL, version, destPath); err != nil {
			return nil, err
		}
	}
	return diskfs.NewFilespace(destPath)
}

func (repos *Repositories) clone(url string, version repositories.Version, destPath string) (err error) {
	if _, err = repos.deps.Connector.Clone(url, version, destPath); err != nil {
		return err
	}
	return nil
}

func (repos *Repositories) update(destPath string) (err error) {
	var (
		repo repositories.Repository
	)
	if repo, err = repos.deps.Connector.Open(destPath); err != nil {
		return err
	}
	return repo.Pull()
}

func (repos *Repositories) srcpath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	path := home + "/.goatcli/src/"
	if err = os.MkdirAll(path, 0766); err != nil {
		return "", err
	}
	return path, nil
}
