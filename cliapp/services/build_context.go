package services

/*
type GeneratorExecutor interface {
	ExecuteTask(task GeneratorTask) (err error)
}

// GeneratorTask is single task for generator
type GeneratorTask struct {
	// Context helper values
	//From, To string
	// Template to run
	Template TemplateHandler
	// external data
	DotData interface{}
	// External properties
	BuildProperties map[string]string
	// destination filesystem
	FS     filesystem.Filespace
	FSPath string
}

// TemplateHandler describe template and block to run
type TemplateHandler struct {
	Layout string
	Path   string
	Name   string
}

/*
// BuildContext is context object
type BuildContext interface {
	NewChildContext(templateData TemplateData) (instance BuildContext, err error)
	AppendError(err error)
	Errors(err error) []error
	ToError() error
	ExecuteHook(name string) (err error)
	SelfExecute() (err error)
	Execute() (err error)
	RenderOnce(destPath, templateName, templatePath string, templateData interface{}) (err error)
}

// TemplateData is data for template to execution
type TemplateData struct {
	From, To     string
	Layout       string
	TemplatePath string
	TemplateName string
	CurrentData  interface{}
}

// ContextData context data
type ContextData struct {
	// From, To   string
	Data       map[string]string
	Properties ContextProperties
}

// ContextProperties contains context properties
type ContextProperties struct {
	Build   map[string]string
	Project map[string]string
	Secrets map[string]string
}
*/
