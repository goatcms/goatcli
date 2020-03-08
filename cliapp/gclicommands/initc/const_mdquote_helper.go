package initc

// MdQuoteHelperFile is content of .goat/build/helpers/md/mdquote_helper.tmpl file
const MdQuoteHelperFile = `{{- define "mdquote" -}}
{{- if ne .text "" }}
*"{{.text}}"*
{{- if .author }}
 - {{.author}}
{{- end }}
{{ end }}
{{- end -}}`
