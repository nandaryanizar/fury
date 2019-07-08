package model_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/nandaryanizar/fury/model"
)

type Account struct {
	UserID  int `fury:"primary_key"`
	Counter int `fury:"auto_increment"`
}

type Account2 struct {
	UserID  int  `fury:"primary_key"`
	Counter *int `fury:"auto_increment"`
}

type Account3 struct {
	UserID    int `fury:"primary_key,auto_increment"`
	Counter   *int
	Name      string
	IsActive  bool
	LastLogin time.Time
}

func TestGetColumnsAndValues(t *testing.T) {
	num := 1
	tm := time.Now()
	cases := []struct {
		have        interface{}
		wantColumns []string
		wantArgs    []interface{}
	}{
		{
			&Account3{UserID: 123, Counter: &num, Name: "Test", IsActive: true, LastLogin: tm},
			[]string{"counter", "name", "isactive", "lastlogin"},
			[]interface{}{&num, "Test", true, tm},
		},
	}

	for _, tc := range cases {
		models, _, err := model.NewModels(tc.have)
		if err != nil {
			t.Error(err)
		}

		if len(models) < 1 {
			t.Error("Models should have length more than zero")
		}

		m := models[0]

		cols, args := m.GetColumnNamesAndValues(false)

		if !reflect.DeepEqual(tc.wantColumns, cols) || !reflect.DeepEqual(tc.wantArgs, args) {
			t.Errorf("Error: expected %v and %v, found %v and %v", tc.wantColumns, tc.wantArgs, cols, args)
		}
	}
}

func TestGetScanner(t *testing.T) {
	testInt := 2
	cases := []struct {
		have    interface{}
		columns []string
		want    []interface{}
	}{
		{
			&[]*Account{
				&Account{
					UserID:  123,
					Counter: 1,
				},
			},
			[]string{"userid", "counter"},
			[]interface{}{},
		},
		{
			&[]*Account2{
				&Account2{
					UserID:  123,
					Counter: &testInt,
				},
			},
			[]string{"userid", "counter"},
			[]interface{}{},
		},
	}

	for _, tc := range cases {
		_, m, err := model.NewModels(tc.have)
		if err != nil {
			t.Error(err)
		}

		scanner := m.GetScanPtrByColumnNames(tc.columns)

		f := reflect.ValueOf(tc.have).Elem().Index(0)
		if f.Kind() == reflect.Ptr {
			f = f.Elem()
		}

		field1 := f.Field(0)
		field2 := f.Field(1)
		tc.want = []interface{}{field1.Addr().Interface(), field2.Addr().Interface()}

		if !reflect.DeepEqual(scanner, tc.want) {
			t.Errorf("Error: expected %v, found %v", tc.want, scanner)
		}
	}
}

func TestNewModelStruct(t *testing.T) {
	cases := []struct {
		have interface{}
		want []*model.Model
	}{
		{
			&Account{
				UserID:  123,
				Counter: 1,
			},
			[]*model.Model{
				&model.Model{
					Name:   "account",
					Fields: make(map[string]*model.Field),
				},
			},
		},
	}

	for _, tc := range cases {
		models, m, err := model.NewModels(tc.have)
		if err != nil {
			t.Error(err)
		}

		tc.want[0].Fields["userid"] = m.Fields["userid"]
		tc.want[0].Fields["counter"] = m.Fields["counter"]
		tc.want[0].FieldSlice = m.FieldSlice

		tc.want[0].PrimaryKeys = append(tc.want[0].PrimaryKeys, m.Fields["userid"])
		tc.want[0].ScanAddr = tc.have
		tc.want[0].Type = reflect.ValueOf(tc.have).Type()

		if !reflect.DeepEqual(tc.want, models) {
			t.Errorf("Error: expected %v, found %v", tc.want, models)
		}
	}
}

func TestNewModelSlice(t *testing.T) {
	cases := []struct {
		have interface{}
		want []*model.Model
	}{
		{
			&[]*Account{
				&Account{
					UserID:  123,
					Counter: 1,
				},
			},
			[]*model.Model{
				&model.Model{
					Name:   "account",
					Fields: make(map[string]*model.Field),
				},
			},
		},
	}

	for _, tc := range cases {
		models, m, err := model.NewModels(tc.have)
		if err != nil {
			t.Error(err)
		}

		tc.want[0].Fields["userid"] = m.Fields["userid"]
		tc.want[0].Fields["counter"] = m.Fields["counter"]
		tc.want[0].FieldSlice = m.FieldSlice

		tc.want[0].PrimaryKeys = append(tc.want[0].PrimaryKeys, m.Fields["userid"])
		tc.want[0].ScanAddr = tc.have

		tc.want[0].Type = reflect.ValueOf(tc.have).Elem().Index(0).Type()

		if !reflect.DeepEqual(tc.want, models) {
			t.Errorf("Error: expected %v, found %v", tc.want, models)
		}
	}
}

func TestNewModelEmptySlice(t *testing.T) {
	cases := []struct {
		have interface{}
		want []*model.Model
	}{
		{
			&[]*Account{&Account{}},
			[]*model.Model{
				&model.Model{
					Name:   "account",
					Fields: make(map[string]*model.Field),
				},
			},
		},
	}

	for _, tc := range cases {
		models, m, err := model.NewModels(tc.have)
		if err != nil {
			t.Error(err)
		}

		tc.want[0].Fields["userid"] = m.Fields["userid"]
		tc.want[0].Fields["counter"] = m.Fields["counter"]
		tc.want[0].FieldSlice = m.FieldSlice

		tc.want[0].PrimaryKeys = append(tc.want[0].PrimaryKeys, m.Fields["userid"])
		tc.want[0].ScanAddr = tc.have

		tc.want[0].Type = reflect.ValueOf(tc.have).Elem().Type().Elem()

		if !reflect.DeepEqual(tc.want, models) {
			t.Errorf("Error: expected %v, found %v", tc.want, models)
		}
	}
}
