package initc

// ReadmeFile is content of .goat/build/templates/app/readme.render file
const ReadmeFile = `# {{index $ctx.Properties.Project "app.name"}}
{{- template "mdquote" (dict "text" (index $ctx.Properties.Project "app.slogan")) }}
{{index $ctx.Properties.Project "app.description"}}

# Authors
{{- range $index, $key := (keys $ctx.PlainData "appauthor.") -}}
{{- $firstname := (index $ctx.PlainData (print "appauthor." $key ".firstname")) -}}
{{- $lastname := (index $ctx.PlainData (print "appauthor." $key ".lastname")) -}}
{{- $email := (index $ctx.PlainData (print "appauthor." $key ".email")) }}
* {{camelcaseuf $firstname}} {{camelcaseuf $lastname}} <{{$email}}>
{{- end -}}`
