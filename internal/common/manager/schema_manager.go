package manager

import (
	"github.com/spaghettifunk/norman/internal/common/schema"
	"github.com/spaghettifunk/norman/pkg/consul"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	// bring the first letter of each field to capital
	// this is necessary to export the struct in go to access
	// the values when we dynamically read the field during
	// the ingestion job
	for _, dimension := range s.DimensionFieldSpecs {
		dimName := cases.Title(language.Und, cases.NoLower).String(dimension.Name)
		dimension.Name = dimName
	}

	for _, metric := range s.MetricFieldSpecs {
		metrName := cases.Title(language.Und, cases.NoLower).String(metric.Name)
		metric.Name = metrName
	}

	for _, dt := range s.DateTimeFieldSpecs {
		dtName := cases.Title(language.Und, cases.NoLower).String(dt.Name)
		dt.Name = dtName
	}

	return sm.consul.PutSchemaConfiguration(s)
}

func (sm *SchemaManager) Shutdown() error {
	return nil
}
