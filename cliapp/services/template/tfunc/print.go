package tfunc

import "html/template"

// Noescape return noescaped content
func Noescape(value string) template.HTML {
	return template.HTML(value)
}
