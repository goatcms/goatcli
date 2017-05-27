package clonec

import (
	"fmt"

	"github.com/goatcms/goatcore/app"
	"gopkg.in/src-d/go-git.v4"
)

// Run run command in app.App context
func Run(a app.App) error {
	var deps struct {
		Command       string `argument:"$1"`
		RepositoryURL string `argument:"?$2"`
		DestPath      string `argument:"?$3"`
	}
	if err := a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if deps.RepositoryURL == "" {
		fmt.Println("Unknown url to clone")
		return nil
	}
	if deps.DestPath == "" {
		fmt.Println("Unknown destination path")
		return nil
	}
	_, err := git.PlainClone(deps.DestPath, false, &git.CloneOptions{
		URL: deps.RepositoryURL,
	})
	if err != nil {
		return err
	}
	fmt.Println("cloned")
	return nil
}
