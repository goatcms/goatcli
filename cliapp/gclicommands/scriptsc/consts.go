package scriptsc

import (
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipservices"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipservices/namespaces"
)

const (
	// FirstKeyParameterIsRequired is error message for lost name
	FirstKeyParameterIsRequired = "First parameter name is required"
)

var (
	// defaultNamespace is default namespace for main task
	defaultNamespace = namespaces.NewNamespaces(pipservices.NamasepacesParams{
		Task: "",
		Lock: "",
	})
)
