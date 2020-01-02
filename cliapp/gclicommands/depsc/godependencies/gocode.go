package godependencies

import (
	"fmt"
	"strings"
)

// FindImports return all imports paths from golang code
func FindImports(code string) (imports []string, err error) {
	var tmp []string
	code = reduceMultilineStrings(removeComments(code))
	if imports, err = findInlineImports(code); err != nil {
		return nil, err
	}
	if tmp, err = findMultilineImports(code); err != nil {
		return nil, err
	}
	return append(imports, tmp...), nil
}

func findInlineImports(code string) (bodies []string, err error) {
	opens := openInlineImportReg.FindAllStringIndex(code, -1)
	bodies = make([]string, len(opens))
	for i, row := range opens {
		close := closeInlineImportReg.FindStringIndex(code[row[1]:])
		if close == nil {
			return nil, fmt.Errorf("godependencies.FindInlineImportBodies: unclose inline import at:\n %v", code[row[0]:])
		}
		bodies[i] = code[row[1] : row[1]+close[0]]
	}
	return bodies, nil
}

func findMultilineImports(code string) (bodies []string, err error) {
	opens := openMultilineImportReg.FindAllStringIndex(code, -1)
	bodies = make([]string, len(opens))
	for i, row := range opens {
		close := closeMultilineImportReg.FindStringIndex(code[row[1]:])
		if close == nil {
			return nil, fmt.Errorf("godependencies.findMultilineImportBodies: unclose multiline import at:\n %v", code[row[0]:])
		}
		bodies[i] = code[row[1] : row[1]+close[0]]
	}
	return reduceMultilineImportBodies(bodies), nil
}

func reduceMultilineImportBodies(bodies []string) (imports []string) {
	imports = []string{}
	for _, body := range bodies {
		rows := strings.Split(body, "\"")
		for i, impo := range rows {
			if i%2 == 0 {
				continue
			}
			imports = append(imports, impo)
		}
	}
	return imports
}

func removeComments(code string) string {
	code = inlineCommentReg.ReplaceAllString(code, "\n")
	return multilineCommentReg.ReplaceAllString(code, "")
}

func reduceMultilineStrings(code string) string {
	return multilineStrings.ReplaceAllString(code, "``")
}
