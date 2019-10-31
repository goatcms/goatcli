package executor

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// Context contains template dot object with data and APIs
type Context struct {
	AM         interface{}
	DotData    interface{}
	PlainData  map[string]string
	Properties GlobalProperties

	// internal variables
	secretsMU sync.RWMutex
	secrets   []*config.Property
	executor  *SecretsExecutor
}

func newContext(executor *SecretsExecutor, sharedData SharedData) *Context {
	return &Context{
		AM:        sharedData.AppData.AM,
		DotData:   sharedData.DotData,
		PlainData: sharedData.AppData.Plain,
		Properties: GlobalProperties{
			Project: sharedData.Properties.Project,
		},
		executor: executor,
		secrets:  []*config.Property{},
	}
}

// Error return an execution error
func (c *Context) Error(msg string) (err error) {
	err = fmt.Errorf(msg)
	c.executor.Scope().AppendError(err)
	return err
}

// Error return an execution error
func (c *Context) AddSecret(dict map[string]interface{}) (err error) {
	var (
		min, max int
		prompt   string
		pattern  *regexp.Regexp
	)
	if _, ok := dict["Key"]; !ok || dict["Key"] == "" {
		return goaterr.Errorf("cliapp/services/secrets/executor/context.go.AddSecret: Key is required")
	}
	if _, ok := dict["Type"]; !ok || dict["Type"] == "" {
		return goaterr.Errorf("cliapp/services/secrets/executor/context.go.AddSecret: Type is required")
	}
	if _, ok := dict["Min"]; ok {
		max = dict["Min"].(int)
	}
	if _, ok := dict["Max"]; ok {
		max = dict["Max"].(int)
	} else {
		max = 100000
	}
	if _, ok := dict["Prompt"]; ok {
		prompt = dict["Prompt"].(string)
	}
	if _, ok := dict["Pattern"]; ok {
		if pattern, err = regexp.Compile(dict["Pattern"].(string)); err != nil {
			return err
		}
	}
	c.secretsMU.Lock()
	defer c.secretsMU.Unlock()
	c.secrets = append(c.secrets, &config.Property{
		Key:     dict["Key"].(string),
		Type:    dict["Type"].(string),
		Prompt:  prompt,
		Min:     min,
		Max:     max,
		Pattern: pattern,
	})
	return nil
}
