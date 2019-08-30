package entitymodel

import (
	"testing"
)

func TestNewRelations(t *testing.T) {
	var (
		data      map[string]string
		relations Relations
		err       error
		ok        bool
	)
	t.Parallel()
	data = map[string]string{
		"prefix.1.name": "user",
		"prefix.1.to":   "owner",
		"prefix.2.name": "receiver.email",
		"prefix.2.to":   "email",
	}
	if relations, err = NewRelationsFromPlainmap("prefix", data); err != nil {
		t.Error(err)
		return
	}
	if _, ok = relations["User"]; !ok {
		t.Errorf("expected User relation")
		return
	}
	if relations["User"].To != "owner" {
		t.Errorf("expected 'User.To' must be equals to 'owner'")
		return
	}
	if _, ok = relations["ReceiverEmail"]; !ok {
		t.Errorf("expected ReceiverEmail relation")
		return
	}
	if relations["ReceiverEmail"].To != "email" {
		t.Errorf("expected 'ReceiverEmail.To' must be equals to 'email'")
		return
	}
}

func TestRelationNameIsRequired(t *testing.T) {
	var (
		data map[string]string
		err  error
	)
	t.Parallel()
	data = map[string]string{
		"prefix.somestring.name": "",
		"prefix.somestring.to":   "text",
	}
	if _, err = NewRelationsFromPlainmap("prefix", data); err == nil {
		t.Errorf("NewRelations should return error. Name field is required")
		return
	}
}

func TestRelationToIsRequired(t *testing.T) {
	var (
		data map[string]string
		err  error
	)
	t.Parallel()
	data = map[string]string{
		"prefix.somestring.entity": "name",
		"prefix.somestring.type":   "",
	}
	if _, err = NewRelationsFromPlainmap("prefix", data); err == nil {
		t.Errorf("NewRelations should return error. To field is required")
		return
	}
}
