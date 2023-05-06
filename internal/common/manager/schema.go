package manager

type SchemaManager struct {
}

func NewSchemaManager() *SchemaManager {
	return &SchemaManager{}
}

func (sm *SchemaManager) Initialize() error {
	return nil
}

func (sm *SchemaManager) Start() error {
	return nil
}

func (sm *SchemaManager) Shutdown() error {
	return nil
}
