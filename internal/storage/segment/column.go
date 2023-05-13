package segment

import (
	"github.com/spaghettifunk/norman/internal/common/types"
)

type Column struct {
	Name      string          `json:"-"`
	FieldType types.FieldType `json:"-"`
}

func NewColumn(name string, ft types.FieldType) *Column {
	return &Column{
		Name:      name,
		FieldType: ft,
	}
}
