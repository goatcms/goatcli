package simpletf

// Wrap return text wraped with '{{' and end with '}}'
// Deprecated: It is now required.
func Wrap(value string) string {
	return "{{" + value + "}}"
}

// Noescape return noescaped content
// Deprecated: It is now required.
func Noescape(value string) string {
	return value
}

// SafeHTMLAttr return noescaped content
// Deprecated: It is now required.
func SafeHTMLAttr(s string) string {
	return s
}

// SafeHTML return noescaped content
// Deprecated: It is now required.
func SafeHTML(s string) string {
	return s
}

// SafeURL return noescaped content
// Deprecated: It is now required.
func SafeURL(s string) string {
	return s
}
