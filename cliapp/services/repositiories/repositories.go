package repositories

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
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
func Factory(dp dependency.Provider) (interface{}, error) {
	r := &Repositories{}
	if err := dp.InjectTo(&r.deps); err != nil {
		return nil, err
	}
	return services.RepositoriesService(r), nil
}

// Filespace return filespace for repository
func (repos *Repositories) Filespace(repoURL, rev string) (filesystem.Filespace, error) {
	basepath, err := repos.srcpath()
	if err != nil {
		return nil, err
	}
	if rev == "" {
		rev = "master"
	}
	destPath := basepath + reduceRepoURL(repoURL) + "." + rev
	if disk.IsDir(destPath) {
		if err = repos.update(destPath); err != nil {
			return nil, err
		}
	} else {
		os.RemoveAll(destPath)
		if err = repos.clone(repoURL, rev, destPath); err != nil {
			return nil, err
		}
	}
	return diskfs.NewFilespace(destPath)
}

func (repos *Repositories) clone(url, version, destPath string) (err error) {
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
