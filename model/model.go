package model

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Model struct
type Model struct {
	Name        string
	Fields      map[string]*Field
	FieldSlice  []*Field
	PrimaryKeys []*Field
	Type        reflect.Type
	ScanAddr    interface{}
}

// GetColumnNamesAndValues return names and values as slice
func (m *Model) GetColumnNamesAndValues(includeAutoInc bool) ([]string, []interface{}) {
	cols := []string{}
	args := []interface{}{}

	for _, f := range m.FieldSlice {
		if f.IsIgnored || (f.IsPrimaryKey && f.CheckIfZeroValue()) || (f.IsAutoIncrement && !includeAutoInc) {
			continue
		}

		cols = append(cols, strings.ToLower(f.Properties.Name))
		args = append(args, f.Value.Interface())
	}

	return cols, args
}

// GetScanPtrByColumnNames return scanner pointers ordered as specifed in the input slice.
func (m *Model) GetScanPtrByColumnNames(columns []string) []interface{} {
	var pointers []interface{}

	for _, col := range columns {
		if f, ok := m.Fields[col]; ok {
			pointers = append(pointers, f.Value.Addr().Interface())
		}
	}

	return pointers
}

// NewModels creates new Model literal
//  Return slice of pointer to models, first pointer to model, and error
//  If slice of pointer to models is empty then the second parameter return newly created pointer to model.
//	Use the second parameter as base type or shortcut to first element of slice
func NewModels(modelInterface interface{}) ([]*Model, *Model, error) {
	// Check if model is valid
	reflectVal := reflect.ValueOf(modelInterface)
	if !reflectVal.IsValid() {
		return nil, nil, errors.New("Error: invalid model")
	}

	// Get model and name of base type
	modelType := reflectVal.Type()
	// Delete later
	// name := modelType.Name()

	// Check if model is pointer, if pointer then get the pointer points to
	if reflectVal.Kind() == reflect.Ptr {
		reflectVal = reflectVal.Elem()

		if reflectVal.Kind() == reflect.Struct {
			m, err := newSingleModel(reflectVal, modelType, modelInterface)
			if err != nil {
				return nil, nil, err
			}

			return []*Model{m}, m, nil
		}
	}

	// Check if the value is slice, if slice then iterate slice or create one model from slice
	if (reflectVal.Kind() == reflect.Slice || reflectVal.Kind() == reflect.Array) && reflectVal.Type().Elem().Kind() == reflect.Ptr {
		models := []*Model{}
		mod := &Model{}
		modelType = reflectVal.Type().Elem()

		for i := 0; i < reflectVal.Len(); i++ {
			val := reflectVal.Index(i)

			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			m, err := newSingleModel(val, modelType, modelInterface)
			if err != nil {
				return nil, nil, err
			}

			if i == 0 {
				mod = m
			}

			models = append(models, m)
		}

		if reflectVal.Len() == 0 {
			newModelType := modelType
			if newModelType.Kind() == reflect.Ptr {
				newModelType = modelType.Elem()
			}

			val := reflect.New(newModelType)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			m, err := newSingleModel(val, modelType, modelInterface)
			if err != nil {
				return nil, nil, err
			}
			mod = m
		}

		return models, mod, nil
	}

	// If model is not struct then return error
	return nil, nil, fmt.Errorf("Error: expected pointer to struct, slice of pointer to struct or pointer to slice of pointer to struct, found %v", reflectVal.Kind())
}

func newSingleModel(structVal reflect.Value, modelType reflect.Type, modelInterface interface{}) (*Model, error) {
	if structVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Error: expected struct, found %v", structVal.Kind())
	}

	// Create model literal
	m := &Model{
		Name:     strings.ToLower(structVal.Type().Name()),
		Fields:   make(map[string]*Field),
		Type:     modelType,
		ScanAddr: modelInterface,
	}

	// Iterate through struct fields
	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)

		if !field.IsValid() {
			continue
		}

		fieldProperties := structVal.Type().Field(i)
		furyField := NewField(fieldProperties, field)
		fieldPropertiesName := strings.ToLower(fieldProperties.Name)

		if _, ok := m.Fields[fieldPropertiesName]; !ok {
			m.Fields[fieldPropertiesName] = furyField
			m.FieldSlice = append(m.FieldSlice, furyField)
		}

		if furyField.IsPrimaryKey {
			m.PrimaryKeys = append(m.PrimaryKeys, furyField)
		}
	}

	return m, nil
}
