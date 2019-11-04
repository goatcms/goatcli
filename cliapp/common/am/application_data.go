package am

import (
	"github.com/goatcms/goatcli/cliapp/services"
)

// NewApplicationData create a new application data
func NewApplicationData(plain map[string]string) services.ApplicationData {
	return services.ApplicationData{
		AM:    NewApplicationModel(plain),
		Plain: plain,
	}
}
