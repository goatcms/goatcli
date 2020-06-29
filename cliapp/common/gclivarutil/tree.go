package gclivarutil

import (
	"regexp"
	"strings"

	"github.com/goatcms/goatcore/varutil/goaterr"
)

var nameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// ToTree conavert a plain string map to a multi level tree
func ToTree(source map[string]string) (out map[string]interface{}, err error) {
	var (
		node         map[string]interface{}
		end          int
		propertyName string
	)
	out = make(map[string]interface{})
	for spath, value := range source {
		if spath == "" {
			return nil, goaterr.Errorf("ToTree: empty path is no allowd")
		}
		end = strings.LastIndex(spath, ".")
		if end != -1 {
			if node, err = toTreeCreateNode(out, spath[:end]); err != nil {
				return nil, err
			}
			propertyName = spath[end+1:]
		} else {
			node = out
			propertyName = spath
		}
		if err = toTreeValidName(spath, propertyName); err != nil {
			return nil, err
		}
		node[propertyName] = value
	}
	return out, nil
}

func toTreeCreateNode(rmap map[string]interface{}, spath string) (node map[string]interface{}, err error) {
	var (
		v     interface{}
		ok    bool
		start int
		end   int
		name  string
		stop  bool
	)
	node = rmap
	end = 0
	if spath[0] == '.' {
		return nil, goaterr.Errorf("empty name is not allowed %s", spath)
	}
	for ; !stop; end++ {
		start = end
		end = dotPos(spath, end)
		if end == -1 {
			stop = true
			end = len(spath)
		}
		name = spath[start:end]
		if v, ok = node[name]; !ok {
			if err = toTreeValidName(spath, name); err != nil {
				return nil, err
			}
			child := make(map[string]interface{})
			child["__PATH"] = spath[:end]
			node[name] = child
			node = child
			continue
		}
		switch child := v.(type) {
		case map[string]interface{}:
			node = child
			continue
		default:
			return nil, goaterr.Errorf("ToTree: %v is not a node", spath[:end])
		}
	}
	return node, nil
}

func toTreeValidName(spath, name string) error {
	if name == "" {
		return goaterr.Errorf("Empty name is not allowed ('%s' for '%s')", name, spath)
	}
	if strings.HasPrefix(name, "__") {
		return goaterr.Errorf("Prefix '__' is reserved for 'MAGIC' properties like __PATH ('%s' for '%s')", name, spath)
	}
	if !nameRegexp.MatchString(name) {
		return goaterr.Errorf("Key can contains a-z, A-Z, 0-9 and '_' character  ('%s' for '%s')", name, spath)
	}
	return nil
}

func dotPos(str string, start int) int {
	for start < len(str) && str[start] != '.' {
		start++
	}
	if start == len(str) {
		return -1
	}
	return start
}
