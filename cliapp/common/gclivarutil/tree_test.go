package gclivarutil

import (
	"testing"
)

func TestToTree(t *testing.T) {
	t.Parallel()
	var (
		err    error
		result map[string]interface{}
		node   map[string]interface{}
		ok     bool
		value  interface{}
		svalue string
	)
	if result, err = ToTree(map[string]string{
		"dep.key":  "a",
		"simpekey": "b",
	}); err != nil {
		t.Error(err)
		return
	}
	// check simpekey
	if value, ok = result["simpekey"]; !ok {
		t.Errorf("expected simpekey")
		return
	}
	if svalue, ok = value.(string); !ok {
		t.Errorf("expected simpekey as string type")
		return
	}
	if svalue != "b" {
		t.Errorf("expected dep.key equals to 'b'")
		return
	}
	// check dep.key
	if value, ok = result["dep"]; !ok {
		t.Errorf("expected dep node")
		return
	}
	if node, ok = value.(map[string]interface{}); !ok {
		t.Errorf("expected dep as map[string]interface{} type")
		return
	}
	if value, ok = node["key"]; !ok {
		t.Errorf("expected dep.key value")
		return
	}
	if svalue, ok = value.(string); !ok {
		t.Errorf("expected dep.key as string value")
		return
	}
	if svalue != "a" {
		t.Errorf("expected dep.key equals to 'a'")
		return
	}
}

func TestToTreeMagicPathValue(t *testing.T) {
	t.Parallel()
	var (
		err          error
		result       map[string]interface{}
		node1, node2 map[string]interface{}
		ok           bool
		value        interface{}
		svalue       string
	)
	if result, err = ToTree(map[string]string{
		"node1.node2.key": "value",
	}); err != nil {
		t.Error(err)
		return
	}
	// check node1
	if value, ok = result["node1"]; !ok {
		t.Errorf("expected node1 %s\n%v", value, result)
		return
	}
	if node1, ok = value.(map[string]interface{}); !ok {
		t.Errorf("expected node1 as map[string]interface{} type %s\n%v", value, result)
		return
	}
	if value, ok = node1["__PATH"]; !ok {
		t.Errorf("expected node1.__PATH value %s\n%v", value, result)
		return
	}
	if svalue, ok = value.(string); !ok {
		t.Errorf("expected node1.__PATH as string %s\n%v", value, result)
		return
	}
	if svalue != "node1" {
		t.Errorf("expected node1.__PATH value equals to 'node1' and take %s\n%v", value, result)
		return
	}
	if value, ok = node1["node2"]; !ok {
		t.Errorf("expected node2 %s\n%v", value, result)
		return
	}
	if node2, ok = value.(map[string]interface{}); !ok {
		t.Errorf("expected node2 as map[string]interface{} type %s\n%v", value, result)
		return
	}
	if value, ok = node2["__PATH"]; !ok {
		t.Errorf("expected node1.node2.__PATH value %s\n%v", value, result)
		return
	}
	if svalue, ok = value.(string); !ok {
		t.Errorf("expected node1.node2.__PATH as string %s\n%v", value, result)
		return
	}
	if svalue != "node1.node2" {
		t.Errorf("expected node1.node2.__PATH value equals to 'node1.node2' and take %s\n%v", value, result)
		return
	}
	if value, ok = node2["key"]; !ok {
		t.Errorf("expected node1.node2.key value %s\n%v", value, result)
		return
	}
	if svalue, ok = value.(string); !ok {
		t.Errorf("expected node1.node2..key as string value %s\n%v", value, result)
		return
	}
	if svalue != "value" {
		t.Errorf("expected node1.node2.key equals to 'value' %s\n%v", value, result)
		return
	}
}

func TestToTreeMagicNodesValue(t *testing.T) {
	t.Parallel()
	var (
		err    error
		result map[string]interface{}
		node1  map[string]interface{}
		ok     bool
		value  interface{}
		lvalue map[string]interface{}
	)
	if result, err = ToTree(map[string]string{
		"root1.node1.key": "value",
		"root1.node2.key": "value",
		"root2.node2.key": "value",
	}); err != nil {
		t.Error(err)
		return
	}
	// check root __NODES
	if value, ok = result["__NODES"]; !ok {
		t.Errorf("expected __NODES %s\n%v", value, result)
		return
	}
	if lvalue, ok = value.(map[string]interface{}); !ok {
		t.Errorf("expected __NODES as map[string]interface{} type %s\n%v", value, result)
		return
	}
	if len(lvalue) != 2 {
		t.Errorf("expected __NODES contains two elements %s\n%v", value, result)
		return
	}
	// check root1
	if value, ok = result["root1"]; !ok {
		t.Errorf("expected root1 %s\n%v", value, result)
		return
	}
	if node1, ok = value.(map[string]interface{}); !ok {
		t.Errorf("expected root1 as map[string]interface{} type %s\n%v", value, result)
		return
	}
	if value, ok = node1["__NODES"]; !ok {
		t.Errorf("expected root1.__NODES value %s\n%v", value, result)
		return
	}
	if lvalue, ok = value.(map[string]interface{}); !ok {
		t.Errorf("expected root1.__NODES as map[string]interface{} %s\n%v", value, result)
		return
	}
	if len(lvalue) != 2 {
		t.Errorf("expected root1.__NODES contains two elements %s\n%v", value, result)
		return
	}
}

func TestToTreeMagicValues(t *testing.T) {
	t.Parallel()
	var (
		err    error
		result map[string]interface{}
		node1  map[string]interface{}
		ok     bool
		value  interface{}
		lvalue map[string]interface{}
	)
	if result, err = ToTree(map[string]string{
		"node.key1": "value",
		"key1":      "value1",
		"key2":      "value2",
	}); err != nil {
		t.Error(err)
		return
	}
	// check root __VALUES
	if value, ok = result["__VALUES"]; !ok {
		t.Errorf("expected __VALUES %s\n(from %v)", value, result)
		return
	}
	if lvalue, ok = value.(map[string]interface{}); !ok {
		t.Errorf("expected __VALUES as map[string]interface{} type %s\n(from %v)", value, result)
		return
	}
	if len(lvalue) != 2 {
		t.Errorf("expected __VALUES contains two elements and take %s\n(from %v)", value, result)
		return
	}
	// check node
	if value, ok = result["node"]; !ok {
		t.Errorf("expected node %s\n%v", value, result)
		return
	}
	if node1, ok = value.(map[string]interface{}); !ok {
		t.Errorf("expected node as map[string]interface{} type %s\n%v", value, result)
		return
	}
	if value, ok = node1["__VALUES"]; !ok {
		t.Errorf("expected node.__VALUES value %s\n%v", value, result)
		return
	}
	if lvalue, ok = value.(map[string]interface{}); !ok {
		t.Errorf("expected node.__VALUES as map[string]interface{} %s\n%v", value, result)
		return
	}
	if len(lvalue) != 1 {
		t.Errorf("expected node.__VALUES contains one element %s\n%v", value, result)
		return
	}
}
