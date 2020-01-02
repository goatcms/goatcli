package initc

const (
	// InitBildDefJSON is content of init build.def.json file
	InitBildDefJSON = `[{
    "to":".",
    "template":"main",
    "properties": {
      "title":"Example goatcli project",
      "text":"Hello world"
    },
    "afterBuild": "echo 'some after build command'",
}]`

	// InitMainBuildHelper is content of init .goat/build/helpers/main.tmpl file
	InitMainBuildHelper = `{{- define "main" -}}{{- end -}}`

	// InitMainTemplateMainTMPL is content of init .goat/build/helpers/templates/main/main.tmpl file
	InitMainTemplateMainTMPL = `{{- $ctx := . -}}

{{- $path := (print $ctx.To "/readme.md") -}}
{{- if not ($ctx.Filesystem.IsFile $path) -}}
	{{- $ctx.Out.File $path -}}
		{{- template "main.readme_md" $ctx -}}
	{{- $ctx.Out.EOF -}}
{{- end -}}

{{- $path := (print $ctx.To "/.gitignore") -}}
{{- if not ($ctx.Filesystem.IsFile $path) -}}
	{{- $ctx.Out.File $path -}}
		{{- template "main.gitignore" $ctx -}}
	{{- $ctx.Out.EOF -}}
{{- end -}}`

	// InitMainTemplateReadmeTMPL is content of init .goat/build/helpers/templates/main/readme.tmpl file
	InitMainTemplateReadmeTMPL = `{{- define "main.readme_md" -}}
{{- $ctx := . -}}

{{index $ctx.Properties.Project "app.name"}}
=====
{{index $ctx.Properties.Project "app.description"}}


{{- end -}}
`

	// InitMainTemplateGitignoreTMPL is content of init .goat/build/helpers/templates/main/gitignore.tmpl file
	InitMainTemplateGitignoreTMPL = `{{- define "main.gitignore" -}}
{{- $ctx := . -}}
# .goat - personal data
/.goat/secrets.json

# logs files
*.log

# database
*.db
*.db-journal

# npm
npm-debug.log

# IDE
.DS_STORE
.idea/

# defaults
*.test
*~
*.exe
*.log
debug
{{- end -}}
`

	// InitDataDefJSON is content of init .goat/data.def.json file
	InitDataDefJSON = `[]`

	// InitDependenciesDefJSON is content of init .goat/dependencies.def.json file
	InitDependenciesDefJSON = `[]`

	// InitPropertiesDefJSON is content of init .goat/properties.def.json file
	InitPropertiesDefJSON = `[{
	  "prompt":"Insert application userfriendly name",
	  "key":"app.name",
	  "type":"line",
	  "min":1,
	  "max":50
	},{
	  "prompt":"Insert application description",
	  "key":"app.description",
	  "type":"line",
	  "min":1,
	  "max":50
	},{
	  "prompt":"Insert application default web port",
	  "key":"app.port",
	  "type":"numeric",
	  "min":1,
	  "max":4
	},{
	  "prompt":"Insert application id name",
	  "key":"app.id",
	  "type":"line",
	  "min":1,
	  "max":50,
	  "pattern": "^[a-zA-Z]+[a-zA-Z0-9_]*$"
	},{
	  "prompt":"Insert application default language (ISO 3166-1 ALFA-2 like pl, en...)",
	  "key":"app.lang.default",
	  "type":"line",
	  "min":1,
	  "max":4,
	  "pattern": "^[a-zA-Z]+$"
	},{
	  "prompt":"Insert application base URL (like 'www.mysilt.pl')",
	  "key":"app.baseURL",
	  "type":"line",
	  "min":1,
	  "max":500
	}]`

	// InitReplacesDefJSON is content of init .goat/replaces.def.json file
	InitReplacesDefJSON = `[]`

	// InitSecretsDefJSON is content of init .goat/secrets.def.json file
	InitSecretsDefJSON = `[]`
)
