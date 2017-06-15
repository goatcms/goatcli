package repositories

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/disk"
)

// Repositories provide access to repositories
type Repositories struct {
	deps struct {
	}
}

// Factory create new repositories instance
func Factory(dp dependency.Provider) (interface{}, error) {
	r := &Repositories{}
	if err := dp.InjectTo(&r.deps); err != nil {
		return nil, err
	}
	return services.Repositories(r), nil
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
	return nil, nil
}

func (repos *Repositories) clone(url, rev, destPath string) error {
	var (
		err  error
		repo *git.Repository
		tree *git.Worktree
	)
	if repo, err = git.PlainClone(destPath, false, &git.CloneOptions{
		URL: url,
	}); err != nil {
		return err
	}
	if tree, err = repo.Worktree(); err != nil {
		return err
	}
	if err = tree.Checkout(&git.CheckoutOptions{
		Force:  true,
		Branch: plumbing.ReferenceName(rev),
	}); err != nil {
		return err
	}
	return nil
}

func (repos *Repositories) update(destPath string) error {
	var (
		err  error
		repo *git.Repository
	)
	if repo, err = git.PlainOpen(destPath); err != nil {
		return err
	}
	if err = repo.Pull(&git.PullOptions{}); err != nil {
		return err
	}
	return nil
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
