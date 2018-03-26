package mockups

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/repositories"
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
func (r RepositoriesService) Filespace(repoURL string, version repositories.Version) (filesystem.Filespace, error) {
	key := repoURL + "." + version.Branch + "." + version.Revision
	fs, ok := r.data[key]
	if !ok {
		return nil, fmt.Errorf("key not exist %s", key)
	}
	return fs, nil
}
