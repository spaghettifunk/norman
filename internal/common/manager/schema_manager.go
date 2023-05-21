package manager

import (
	"github.com/spaghettifunk/norman/internal/common/schema"
	"github.com/spaghettifunk/norman/pkg/consul"
)

type SchemaManager struct {
	consul *consul.Consul
}

func NewSchemaManager(c *consul.Consul) *SchemaManager {
	return &SchemaManager{
		consul: c,
	}
}

func (sm *SchemaManager) Initialize() error {
	return nil
}

func (sm *SchemaManager) Execute(s *schema.Schema) error {
	return sm.consul.PutSchemaConfiguration(s)
}

func (sm *SchemaManager) Shutdown() error {
	return nil
}
