package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/spaghettifunk/norman/internal/common/entities"
	"github.com/stretchr/testify/assert"
)

func generateSchemas(t *testing.T) (*entities.Schema, *arrow.Schema) {
	schema := &entities.Schema{
		Name: "test-schema",
		DimensionFieldSpecs: []*entities.DimensionFieldSpec{
			{
				Name:     "dimension-a",
				DataType: "LONG",
				Nullable: true,
			},
			{
				Name:     "dimension-b",
				DataType: "STRING",
				Nullable: false,
			},
		},
		DateTimeFieldSpecs: &entities.DateTimeFieldSpec{
			Name:        "timestamp",
			DataType:    "TIMESTAMP",
			Format:      "1:MILLISECONDS:EPOCH",
			Granularity: "5:MINUTE",
		},
	}

	evtSchema, err := schema.GetFullArrowSchema()
	if err != nil {
		t.Errorf("failed to create arrow schema with error: %s", err.Error())
	}

	return schema, evtSchema
}

func deleteOutputFolder(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to get current working directory with error: %s", err.Error())
	}
	if err := os.RemoveAll(fmt.Sprintf("%s/output", path)); err != nil {
		t.Errorf("failed to delete directory of table manager with error: %s", err.Error())
	}
}

func TestNewTableManager(t *testing.T) {
	schema, evtSchema := generateSchemas(t)
	table := &entities.Table{
		Name:        "test-table",
		Schema:      schema,
		EventSchema: evtSchema,
	}

	tm, err := NewTableManager(table, nil)
	if err != nil {
		t.Errorf("failed to create table manager with error: %s", err.Error())
	}
	assert.NotNil(t, tm)

	deleteOutputFolder(t)
}

func TestCreateNewSegment(t *testing.T) {
	schema, evtSchema := generateSchemas(t)
	table := &entities.Table{
		Name:        "test-table",
		Schema:      schema,
		EventSchema: evtSchema,
	}

	tm, err := NewTableManager(table, nil)
	if err != nil {
		t.Errorf("failed to create table manager with error: %s", err.Error())
	}
	assert.NotNil(t, tm)

	if err := tm.CreateNewSegment(); err != nil {
		t.Errorf("failed to create a new segment with error: %s", err.Error())
	}

	deleteOutputFolder(t)
}

func TestInsertData(t *testing.T) {
	schema, evtSchema := generateSchemas(t)
	table := &entities.Table{
		Name:        "test-table",
		Schema:      schema,
		EventSchema: evtSchema,
	}

	tm, err := NewTableManager(table, nil)
	if err != nil {
		t.Errorf("failed to create table manager with error: %s", err.Error())
	}
	assert.NotNil(t, tm)

	if err := tm.CreateNewSegment(); err != nil {
		t.Errorf("failed to create a new segment with error: %s", err.Error())
	}

	for _, data := range []map[string]interface{}{
		{"dimension-a": 205, "dimension-b": "Yuki", "timestamp": 1685992352000},
		{"dimension-a": 206, "dimension-b": "Lewis", "timestamp": 1685992352000},
		{"dimension-a": 207, "dimension-b": "Max", "timestamp": 1685992772000},
		{"dimension-a": 208, "dimension-b": "Charles", "timestamp": 1685992772000},
	} {
		d, err := json.Marshal(data)
		if err != nil {
			t.Errorf("failed to Marshal fake data with error: %s", err.Error())
		}
		if err := tm.InsertData(d); err != nil {
			t.Errorf("failed to insert data with error: %s", err.Error())
		}
	}

	if err := tm.FlushSegment(); err != nil {
		t.Errorf("failed to flush segment with error: %s", err.Error())
	}

	deleteOutputFolder(t)
}
