package commands

const (
	// Clone is a clone command help
	Clone = "[git_url, destination_path, [--rev=rev]] clone project (and modules)"
	// DataAdd is a data:add command help
	DataAdd = "[type_name, data_name] add new data set to project scaffolding data"
	// Build is a build command help
	Build = "build goat project in current directory"
	// Clean is a clean command help
	Clean = "clean builded files"
	// Init is a init command help
	Init = "initialize new goat project"
	// AddDep is a deps:add command help
	AddDep = "[path, repo_url, [branch, [revision]] Add new static dependency like golang vendor or js node module"
	// AddGODep is a deps:add:go command help
	AddGODep = "[go_src_path, [branch, [revision]] Add new golang dependency like 'github.com/goatcms/goatcore'"
	// AddGOImportsDep is a deps:add:go:imports command help
	AddGOImportsDep = "import commands from project"
	// SetSecretValueDep is a secrets:set command help
	SetSecretValueDep = "[key, value] adds or updates an secret with a specified key and value"
	// GetSecretValueDep is a secrets:set command help
	GetSecretValueDep = "[key] display a specified element from a secrets"
	// SetPropertyValueDep is a properties:set command help
	SetPropertyValueDep = "[key, value] adds or updates an property with a specified key and value"
	// GetPropertyValueDep is a properties:get command help
	GetPropertyValueDep = "[key] display a specified element from a properties"
	// VCSClean is a vcs:clean command help
	VCSClean = "Clean vcs ignored files"
	// VCSIgnoredAdd is a vcs:ignored:add command help
	VCSIgnoredAdd = "Add new vcs ignored file [--path=file path to be ignored]"
	// VCSIgnoredRemove is a vcs:ignored:remove command help
	VCSIgnoredRemove = "Remove a vcs ignored file [--path=file path]"
	// VCSIgnoredList is a vcs:ignored:list command help
	VCSIgnoredList = "Show ignored files listing"
	// VCSGeneratedList is a vcs:generated:list command help
	VCSGeneratedList = "Show generated files listing"
	// CWDArg is a cwd argument description
	CWDArg = "set Current Working Dirctory"
)
