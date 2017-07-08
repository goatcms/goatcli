package cio

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
)

func TestReadPropertisMap(t *testing.T) {
	var (
		err       error
		in        app.Input
		out       app.Output
		isChanged bool
	)
	t.Parallel()
	in = gio.NewInput(strings.NewReader("bk1\nvalue1\nvalue2\nbk2\nv1\nv2\n\n"))
	out = gio.NewOutput(new(bytes.Buffer))
	data := map[string]string{}
	if isChanged, err = ReadPropertiesMap("basekey.", in, out, testProperties, data); err != nil {
		t.Error(err)
		return
	}
	if isChanged != true {
		t.Errorf("expected return isChanged flag equals to true")
		return
	}
	if len(data) != 4 {
		t.Errorf("result data should contains two elements and it have %d", len(data))
		return
	}
	if data["basekey.bk1.key1"] != "value1" {
		t.Errorf("expected data[basekey.bk1.key1] equals to value1 and it is %v", data["basekey.bk1.key1"])
		return
	}
	if data["basekey.bk1.key2"] != "value2" {
		t.Errorf("expected data[basekey.bk1.key2] equals to value2 and it is %v", data["basekey.bk1.key2"])
		return
	}
	if data["basekey.bk2.key1"] != "v1" {
		t.Errorf("expected data[basekey.bk2.key1] equals to value1 and it is %v", data["basekey.bk2.key1"])
		return
	}
	if data["basekey.bk2.key2"] != "v2" {
		t.Errorf("expected data[basekey.bk2.key2] equals to value2 and it is %v", data["basekey.bk2.key2"])
		return
	}
}
