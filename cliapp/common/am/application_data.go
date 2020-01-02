package am

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
)

// NewApplicationData create a new application data
func NewApplicationData(plain map[string]string) gcliservices.ApplicationData {
	return gcliservices.ApplicationData{
		AM:    NewApplicationModel(plain),
		Plain: plain,
	}
}
