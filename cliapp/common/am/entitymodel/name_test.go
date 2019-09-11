package entitymodel

import (
	"testing"
)

func TestNewName(t *testing.T) {
	var (
		name   Name
		expect string
		err    error
	)
	t.Parallel()
	if name, err = NewName("test.new.name"); err != nil {
		t.Error(err)
		return
	}
	expect = "test.new.name"
	if name.Plain != "test.new.name" {
		t.Errorf("name.Plain: value is incorrect. Take %s and expect %s", name.Plain, expect)
		return
	}
	expect = "TestNewName"
	if name.CamelCaseUF != expect {
		t.Errorf("name.CamelCaseUF: value is incorrect. Take %s and expect %s", name.CamelCaseUF, expect)
		return
	}
	expect = "testNewName"
	if name.CamelCaseLF != expect {
		t.Errorf("name.CamelCaseLF: value is incorrect. Take %s and expect %s", name.CamelCaseLF, expect)
		return
	}
	expect = "test_new_name"
	if name.Underscore != expect {
		t.Errorf("name.Underscore: value is incorrect. Take %s and expect %s", name.Underscore, expect)
		return
	}
	expect = "TESTNEWNAME"
	if name.Upper != expect {
		t.Errorf("name.Upper: value is incorrect. Take %s and expect %s", name.Upper, expect)
		return
	}
	expect = "testnewname"
	if name.Lower != expect {
		t.Errorf("name.Lower: value is incorrect. Take %s and expect %s", name.Lower, expect)
		return
	}
}
