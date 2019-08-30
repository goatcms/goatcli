package am

import (
	"sync"

	"github.com/goatcms/goatcli/cliapp/common/am/entitymodel"
)

// ApplicationModel contains application data
type ApplicationModel struct {
	data       map[string]string
	entities   entitymodel.Entities
	entitiesMU sync.Mutex
	options    entitymodel.Entities
	optionsMU  sync.Mutex
	dto        entitymodel.Entities
	dtoMU      sync.Mutex
}

// NewApplicationModel create new ApplicationModel instance
func NewApplicationModel(data map[string]string) *ApplicationModel {
	return &ApplicationModel{
		data: data,
	}
}

// Entities return application entities model
func (am *ApplicationModel) Entities() (model entitymodel.Entities, err error) {
	if am.entities != nil {
		return am.entities, nil
	}
	am.entitiesMU.Lock()
	defer am.entitiesMU.Unlock()
	if am.entities != nil {
		return am.entities, nil
	}
	if am.entities, err = entitymodel.NewEntitiesFromPlainmap("model", am.data); err != nil {
		return nil, err
	}
	return am.entities, nil
}

// Options return application options model
func (am *ApplicationModel) Options() (model entitymodel.Entities, err error) {
	if am.options != nil {
		return am.options, nil
	}
	am.optionsMU.Lock()
	defer am.optionsMU.Unlock()
	if am.options != nil {
		return am.options, nil
	}
	if am.options, err = entitymodel.NewEntitiesFromPlainmap("options", am.data); err != nil {
		return nil, err
	}
	return am.options, nil
}

// DTO return application dto model
func (am *ApplicationModel) DTO() (model entitymodel.Entities, err error) {
	if am.dto != nil {
		return am.dto, nil
	}
	am.dtoMU.Lock()
	defer am.dtoMU.Unlock()
	if am.dto != nil {
		return am.dto, nil
	}
	if am.dto, err = entitymodel.NewEntitiesFromPlainmap("dto", am.data); err != nil {
		return nil, err
	}
	return am.dto, nil
}
