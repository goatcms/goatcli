package initc

// goatStructureSnapshot define goat empty project structure
var goatStructureSnapshot = map[string]string{
	// define build and scaffolding
	".goat/build.def.json": BildDefFile,
	// helpers
	".goat/build/helpers/md/mdquote.tmpl": MdQuoteHelperFile,
	// app templete
	".goat/build/templates/app/readme.md.render": ReadmeFile,
	// vnc - hook: .git structure
	".goat/build/templates/hook/vcs/git/.gitignore.render": VncGitignoreHookFile,
	// data definition
	".goat/data.def/appauthor.json": AuthorDataDefFile,
	// secrets definition
	".goat/secrets.def/00_main.json": MainSecretsFile,
	// properties definition
	".goat/properties.def/00_main.json": PropertiesDefFile,
}
