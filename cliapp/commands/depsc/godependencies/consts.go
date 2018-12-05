package godependencies

import "regexp"

const (
	// MaxImportDepth is default value for max depths during recursive import of dependency
	MaxImportDepth = 404
)

// PathMappingRow is single row to mapping
type PathMappingRow struct {
	From string
	To   string
}

var (
	// AlwaysIgnored is set of ignored strings
	AlwaysIgnored = []*regexp.Regexp{}

	// GOPathMapping is list of default path to replace to non-standard repositories
	GOPathMapping = []PathMappingRow{
		PathMappingRow{
			From: "golang.org/x/",
			To:   "github.com/golang/",
		},
		PathMappingRow{
			From: "google.golang.org/",
			To:   "github.com/golang/",
		},
		PathMappingRow{
			From: "cloud.google.com/go/",
			To:   "github.com/GoogleCloudPlatform/google-cloud-go/",
		},
	}
)
