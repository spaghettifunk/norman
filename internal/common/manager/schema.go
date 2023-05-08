package manager

import "github.com/spaghettifunk/norman/internal/common/model"

type SchemaManager struct {
	Schemas []*model.Schema
}

func NewSchemaManager() *SchemaManager {
	return &SchemaManager{}
}

func (sm *SchemaManager) Initialize() error {
	return nil
}

func (sm *SchemaManager) CreateSchema(config []byte) error {
	s, err := model.NewSchema(config)
	if err != nil {
		return err
	}
	sm.Schemas = append(sm.Schemas, s)

	// TODO: notify Aqua that a new Schema has been created
	// ...

	return nil
}

func (sm *SchemaManager) Start() error {
	return nil
}

func (sm *SchemaManager) Shutdown() error {
	return nil
}
