package am

import (
	"testing"
)

func TestNewRelation(t *testing.T) {
	var (
		data     map[string]string
		relation *Relation
		err      error
	)
	t.Parallel()
	data = map[string]string{
		"prefix.name":     "owner",
		"prefix.to":       "user",
		"prefix.system":   "y",
		"prefix.unique":   "y",
		"prefix.required": "y",
	}
	if relation, err = NewRelation("prefix", data); err != nil {
		t.Error(err)
		return
	}
	if relation.Name.CamelCaseUF != "Owner" {
		t.Errorf("relation.Name.CamelCaseUF should be equals to 'Title' and take %s", relation.Name.CamelCaseUF)
		return
	}
	if relation.To != "user" {
		t.Errorf("relation.To should be equals to 'user' and take %s", relation.Name.CamelCaseUF)
		return
	}
	if relation.Flags.Required != true {
		t.Errorf("relation.Flags.Required should be equals to true")
		return
	}
	if relation.Flags.System != true {
		t.Errorf("relation.Flags.System should be equals to true")
		return
	}
	if relation.Flags.Unique != true {
		t.Errorf("relation.Flags.Unique should be equals to true")
		return
	}
}
