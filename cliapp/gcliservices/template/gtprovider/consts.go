package gtprovider

const (
	// TemplateExtension is a default template file extension
	TemplateExtension = ".tmpl"
	// CtrlTemplateExtension is a default controll block definiction
	CtrlTemplateExtension = ".ctrl"
	// OnceTemplateExtension is a default single template file extension (it is build once)
	OnceTemplateExtension = ".once"
	// DefTemplateExtension is a default single definition file extension (it is not build but can be use in template)
	DefTemplateExtension = ".def"
	// RenderTemplateExtension is a default every time template file extension (it is rebuild each time when build application)
	RenderTemplateExtension = ".render"
)
