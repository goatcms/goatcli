package godependencies

import "regexp"

const (
	// MaxImportDepth is default value for max depths during recursive import of dependency
	MaxImportDepth = 404
)

var (
	// AlwaysIgnored is set of ignored strings
	AlwaysIgnored = []*regexp.Regexp{
		// Top lvl domains
		regexp.MustCompile("^[^/]+/?$"),
		// Top lvl directories
		regexp.MustCompile("^[^/]+/[^/]+/?$"),
	}
)

var (
	openInlineImportReg     = regexp.MustCompile("(\\n|\\A)[\\t\\s]*import[\\t\\s]?\\\"")
	closeInlineImportReg    = regexp.MustCompile("\\\"[\\t\\s]*(\\n|\\z)")
	openMultilineImportReg  = regexp.MustCompile("(\\n|\\A)[\\t\\s]*import[\\t\\s\\n]?\\(")
	closeMultilineImportReg = regexp.MustCompile("[\\t\\s\\n]*\\)")
	inlineCommentReg        = regexp.MustCompile("//[^\\n]*(\\n|\\z)")
	multilineCommentReg     = regexp.MustCompile("\\/\\*([^\\*].*)*\\*/")
	multilineStrings        = regexp.MustCompile("\\`[^\\`]*(\\`|\\z)")
)
