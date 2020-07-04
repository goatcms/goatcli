package gclivarutil

import (
	"regexp"
	"strings"

	"github.com/goatcms/goatcore/varutil/goaterr"
)

var nameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// ToTree conavert a plain string map to a multi level tree
func ToTree(source map[string]string) (root map[string]interface{}, err error) {
	var (
		node         map[string]interface{}
		end          int
		propertyName string
	)
	root = newTreeNode("")
	for spath, value := range source {
		if spath == "" {
			return nil, goaterr.Errorf("ToTree: empty path is no allowd")
		}
		end = strings.LastIndex(spath, ".")
		if end != -1 {
			if node, err = toTreeCreateNode(root, spath[:end]); err != nil {
				return nil, err
			}
			propertyName = spath[end+1:]
		} else {
			node = root
			propertyName = spath
		}
		if err = toTreeValidName(spath, propertyName); err != nil {
			return nil, err
		}
		appendTreeValue(node, propertyName, value)
	}
	return root, nil
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
			child := newTreeNode(spath[:end])
			appendTreeChildNode(node, child, name)
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

func newTreeNode(path string) (node map[string]interface{}) {
	node = make(map[string]interface{})
	node["__PATH"] = path
	node["__NODES"] = make(map[string]interface{}, 0)
	node["__VALUES"] = make(map[string]interface{}, 0)
	return node
}

func appendTreeChildNode(parent, child map[string]interface{}, name string) {
	nodes := parent["__NODES"].(map[string]interface{})
	nodes[name] = child
	parent[name] = child
}

func appendTreeValue(node map[string]interface{}, name string, value interface{}) {
	values := node["__VALUES"].(map[string]interface{})
	values[name] = value
	node[name] = value
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
