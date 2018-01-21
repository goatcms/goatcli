package tfunc

import "html/template"

// Noescape return noescaped content
func Noescape(value string) template.HTML {
	return template.HTML(value)
}

// SafeHTMLAttr return noescaped content
func SafeHTMLAttr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

// SafeHTML return noescaped content
func SafeHTML(s string) template.HTML {
	return template.HTML(s)
}

// SafeURL return noescaped content
func SafeURL(s string) template.URL {
	return template.URL(s)
}
