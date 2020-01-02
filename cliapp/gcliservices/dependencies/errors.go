package dependencies

import "fmt"

var (
	errInitRequired = fmt.Errorf("DependenciesService: Init is required before use dependency")
)
