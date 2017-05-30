package services

import "github.com/goatcms/goatcore/filesystem"

const (
	// RepositoriesService provide git repository access
	RepositoriesService = "Repositories"
)

// Repositories provide git repository access
type Repositories interface {
	Filespace(repository, rev string) (filesystem.Filespace, error)
}
