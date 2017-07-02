package mockups

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/filesystem"
)

// RepositoriesService is a simple mockup for  services.Repositories
type RepositoriesService struct {
	data map[string]filesystem.Filespace
}

// NewRepositoriesService create mockup for services.Repositories
func NewRepositoriesService(data map[string]filesystem.Filespace) services.RepositoriesService {
	return RepositoriesService{
		data: data,
	}
}

// Filespace return filespace for repository and revision
func (r RepositoriesService) Filespace(repository, rev string) (filesystem.Filespace, error) {
	key := repository + "!" + rev
	fs, ok := r.data[key]
	if !ok {
		return nil, fmt.Errorf("key not exist %s", key)
	}
	return fs, nil
}
