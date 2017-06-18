package config

import "github.com/goatcms/goatcli/cliapp/common"

type TestStringInjector struct{}

func NewTestStringInjector() common.StringInjector {
	return &TestStringInjector{}
}

func (ti *TestStringInjector) InjectToString(s string) (string, error) {
	return s, nil
}
