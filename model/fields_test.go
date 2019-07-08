package model_test

import (
	"reflect"
	"testing"

	"github.com/nandaryanizar/fury/model"
)

func TestNewField(t *testing.T) {
	cases := []struct {
		have *Account
		want []*model.Field
	}{
		{
			&Account{UserID: 123, Counter: 1},
			[]*model.Field{
				&model.Field{
					IsPrimaryKey:    true,
					IsAutoIncrement: false,
					IsIgnored:       false,
				},
				&model.Field{
					IsPrimaryKey:    false,
					IsAutoIncrement: true,
					IsIgnored:       false,
				},
			},
		},
	}

	for _, tc := range cases {
		field1 := reflect.ValueOf(tc.have).Elem().Field(0)
		fieldProp1 := reflect.ValueOf(tc.have).Elem().Type().Field(0)
		field2 := reflect.ValueOf(tc.have).Elem().Field(1)
		fieldProp2 := reflect.ValueOf(tc.have).Elem().Type().Field(1)

		tc.want[0].Properties = fieldProp1
		tc.want[0].Value = field1
		tc.want[1].Properties = fieldProp2
		tc.want[1].Value = field2

		haveFields := []*model.Field{
			model.NewField(fieldProp1, field1),
			model.NewField(fieldProp2, field2),
		}

		if !reflect.DeepEqual(tc.want, haveFields) {
			t.Errorf("Error: expected %v, found %v", tc.want, haveFields)
		}
	}
}

func TestCheckIfZeroValue(t *testing.T) {
	cases := []struct {
		have interface{}
		want bool
	}{
		{1, false},
		{true, false},
		{"", true},
	}

	for _, tc := range cases {
		f := &model.Field{Value: reflect.ValueOf(tc.have)}

		if f.CheckIfZeroValue() != tc.want {
			t.Errorf("Error: value detected zero %v when checking %v", f.CheckIfZeroValue(), tc.have)
		}
	}
}
