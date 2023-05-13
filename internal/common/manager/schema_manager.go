package manager

import "github.com/spaghettifunk/norman/internal/common/schema"

type SchemaManager struct {
	Schemas []*schema.Schema
	// Aqua gRPC client
}

func NewSchemaManager() *SchemaManager {
	return &SchemaManager{}
}

func (sm *SchemaManager) Initialize() error {
	return nil
}

func (sm *SchemaManager) Execute(config []byte) error {
	s, err := schema.NewSchema(config)
	if err != nil {
		return err
	}
	sm.Schemas = append(sm.Schemas, s)

	// TODO: notify Aqua that a new Schema has been created
	// ...

	return nil
}

func (sm *SchemaManager) Shutdown() error {
	return nil
}
