package initc

// VncGitignoreHookFile is content of .goat/build/templates/hook/vcs/git/.gitignore.render file
const VncGitignoreHookFile = `{{- $vcsData := $ctx.DotData -}}
{{- $generated := $vcsData.VCSGeneratedFiles -}}
{{- $persisted := $vcsData.VCSPersistedFiles -}}

############################
# Builded/Output files
############################

# Add your output files here

############################
# Project files
############################
data

############################
# GoatCLI - private data
############################
.goat/secrets.json
.goat/vcs/generated
.goat/secrets.json
.goat/secrets.enc

############################
# Other system and ide files
############################
.DS_Store
main
npm-debug.log
__debug_bin
*.test
*~
*.exe
debug
*.db
*.db-journal
.idea/
*.log
configuration.json
node_modules/**/*
.expo/*
npm-debug.*
*.jks
*.p8
*.p12
*.key
*.mobileprovision
*.orig.*

############################
# Genereted files
############################
{{- range $index, $row := $generated.All }}
{{- if and (not ($persisted.ContainsPath $row.Path)) (ne $row.Path "/.gitignore") (ne $row.Path "/readme.md") (ne $row.Path "") }}
{{$row.Path}}
{{- end }}
{{- end }}
`
