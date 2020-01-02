package builder

const (
	// BuildDefPath is path to build definitions
	BuildDefPath = ".goat/build.def.json"
)

const (
	// DefaultExecutorLimit is file limit for single executor
	DefaultExecutorLimit = 1000000
)

// TaskData contains default data for single task
type TaskData struct {
	From, To string
}
