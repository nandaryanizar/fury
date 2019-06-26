package model

import (
	"reflect"
	"strings"
)

// Field struct
type Field struct {
	Properties      reflect.StructField
	Value           reflect.Value
	IsPrimaryKey    bool
	IsAutoIncrement bool
	IsIgnored       bool
}

// NewField create new field literal
func NewField(prop reflect.StructField, val reflect.Value) *Field {
	field := &Field{
		Properties: prop,
		Value:      val,
	}

	field.processTagString()

	return field
}

func (f *Field) processTagString() {
	if tag := f.Properties.Tag.Get("fury"); tag != "" {
		tags := strings.Split(tag, ",")
		for _, val := range tags {
			if strings.ToLower(val) == "primary_key" {
				f.IsPrimaryKey = true
			}

			if strings.ToLower(val) == "auto_increment" {
				f.IsAutoIncrement = true
			}
		}
	}
}

// CheckIfZeroValue check if value of the field is the zero value of the type of the field
func (f *Field) CheckIfZeroValue() bool {
	return reflect.Zero(f.Value.Type()).Interface() == f.Value.Interface()
}
