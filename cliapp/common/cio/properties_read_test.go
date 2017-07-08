package cio

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
)

func TestPropertiesRead(t *testing.T) {
	var (
		err       error
		in        app.Input
		out       app.Output
		isChanged bool
	)
	t.Parallel()
	in = gio.NewInput(strings.NewReader(" value1\n value2"))
	out = gio.NewOutput(new(bytes.Buffer))
	data := map[string]string{}
	if isChanged, err = ReadProperties("basekey.", in, out, testProperties, data, map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if isChanged != true {
		t.Errorf("expected return isChanged flag equals to true")
		return
	}
	if len(data) != 2 {
		t.Errorf("result data should contains two elements and it have %d", len(data))
		return
	}
	if data["basekey.key1"] != "value1" {
		t.Errorf("expected data[basekey.key1] equals to value1 and it is %v", data["basekey.key1"])
		return
	}
	if data["basekey.key2"] != "value2" {
		t.Errorf("expected data[basekey.key2] equals to value2 and it is %v", data["basekey.key2"])
		return
	}
}

func TestPropertiesReadEOF(t *testing.T) {
	var (
		err       error
		in        app.Input
		out       app.Output
		isChanged bool
	)
	t.Parallel()
	in = gio.NewInput(strings.NewReader(""))
	out = gio.NewOutput(new(bytes.Buffer))
	data := map[string]string{}
	if isChanged, err = ReadProperties("", in, out, testProperties, data, map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if isChanged != true {
		t.Errorf("Expected return isChanged flag equals to true")
		return
	}
	if len(data) != 2 {
		t.Errorf("Result data should contains two elements and it have %d", len(data))
		return
	}
	if data["key1"] == "" {
		t.Errorf("Expected data[key1] should have generated value. It is equals to %v", data["key1"])
		return
	}
	if data["key2"] == "" {
		t.Errorf("Expected data[key2] should have generated value. It is equals to %v", data["key2"])
		return
	}
}
