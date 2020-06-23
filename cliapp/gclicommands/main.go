package gclicommands

const (
	// Clone is a clone command help
	Clone = "[git_url, destination_path, [--rev=rev]] clone project (and modules)"
	// DataAdd is a data:add command help
	DataAdd = "[type_name, data_name] add new data set to project scaffolding data"
	// Build is a build command help
	Build = "build goat project in current directory"
	// Rebuild is a build command help
	Rebuild = "clean build and run new build"
	// Clean is a clean command help
	Clean = "clean builded files and dependencies"
	// CleanDependencies is a clean command help
	CleanDependencies = "clean dependencies files"
	// CleanBuild is a clean command help
	CleanBuild = "clean build files"
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
	VCSClean = "Clean vcs persisted files"
	// VCSPersistedAdd is a vcs:persisted:add command help
	VCSPersistedAdd = "Add new vcs persisted file [--path=file path to be persisted]"
	// VCSPersistedRemove is a vcs:persisted:remove command help
	VCSPersistedRemove = "Remove a vcs persisted file [--path=file path]"
	// VCSPersistedList is a vcs:persisted:list command help
	VCSPersistedList = "Show persisted files listing"
	// VCSScan is a vcs:scan command help
	VCSScan = "Scan files for changes (and add it to vcs persisted files)"
	// VCSGeneratedList is a vcs:generated:list command help
	VCSGeneratedList = "Show generated files listing"
	// ScriptsRun is a scripts:run command help
	ScriptsRun = "run script by name"
	// ScriptsEnvs is a scripts:envs command help
	ScriptsEnvs = "list scripts environments"
	// CWDArg is a cwd argument description
	CWDArg = "set Current Working Dirctory"
)
