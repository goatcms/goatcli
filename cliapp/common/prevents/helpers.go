package prevents

import (
	"fmt"

	"github.com/goatcms/goatcore/filesystem"
)

var (
	errNoGoatProject = fmt.Errorf("Directory doesn't contain a goat project")
)

// RequireGoatProject return error if filespace does'n contains a goat project
func RequireGoatProject(fs filesystem.Filespace) error {
	if fs.IsDir(".goat") {
		return nil
	}
	return errNoGoatProject
}
