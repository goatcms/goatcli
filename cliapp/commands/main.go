package commands

const (
	// Clone is a clone command help
	Clone = "[git_url, destination_path, [--rev=rev]] clone project (and modules)"
	// DataAdd is a data:add command help
	DataAdd = "[type_name, data_name] add new data set to project scaffolding data"
	// Build is a build command help
	Build = "build goat project in current directory"
	// Init is a init command help
	Init = "initialize new goat project"
	// AddDep is a deps:add command help
	AddDep = "[path, repo_url, [branch, [revision]] Add new static dependency like golang vendor or js node module"
	// AddGODep is a deps:add:go command help
	AddGODep = "[go_src_path, [branch, [revision]] Add new golang dependency like 'github.com/goatcms/goatcore'"
	// AddGOImportsDep is a deps:add:go:imports command help
	AddGOImportsDep = "import commands from project"
	// CWDArg is a cwd argument description
	CWDArg = "set Current Working Dirctory"
)
