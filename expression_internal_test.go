package fury

import (
	"errors"
	"reflect"
	"testing"
)

func TestLogicalExpressionWalkStringReturn(t *testing.T) {
	cases := []struct {
		logicalOperator string
		expression1     interface{}
		expression2     interface{}
		want            string
	}{
		{"OR", IsEqualsTo("key", 1), nil, "key = ?"},
		{"OR", "key = 1", nil, "key = 1"},
		{"OR", IsEqualsTo("key", "1"), nil, "key = ?"},
		{"OR", "key = '1'", nil, "key = '1'"},
		{"OR", IsEqualsTo("key", 1), IsEqualsTo("key", 2), "(key = ? OR key = ?)"},
		{"OR", IsEqualsTo("key", 1), "key = 2", "(key = ? OR key = 2)"},
		{"OR", IsEqualsTo("key", 1), And(IsEqualsTo("key", 2)), "(key = ? OR key = ?)"},
		{"OR", IsEqualsTo("key", 1), And(IsEqualsTo("key", 2), IsEqualsTo("key", 3)), "(key = ? OR (key = ? AND key = ?))"},
	}

	for _, tc := range cases {
		have, _, _ := newLogicalExpression(tc.logicalOperator, tc.expression1, tc.expression2).ToString()
		if have != tc.want {
			t.Errorf("Error: expected %s, found %s", tc.want, have)
		}
	}
}

func TestLogicalExpressionWalkArgsReturn(t *testing.T) {
	cases := []struct {
		logicalOperator string
		expression1     interface{}
		expression2     interface{}
		want            []interface{}
	}{
		{"OR", IsEqualsTo("key", 1), nil, []interface{}{1}},
		{"OR", IsEqualsTo("key", "1"), nil, []interface{}{"1"}},
		{"OR", IsEqualsTo("key", 1), IsEqualsTo("key", 2), []interface{}{1, 2}},
		{"OR", IsEqualsTo("key", 1), And(IsEqualsTo("key", "2")), []interface{}{1, "2"}},
		{"OR", IsEqualsTo("key", 1), And(IsEqualsTo("key", 2), IsEqualsTo("key", true)), []interface{}{1, 2, true}},
	}

	for _, tc := range cases {
		_, have, _ := newLogicalExpression(tc.logicalOperator, tc.expression1, tc.expression2).ToString()
		if !reflect.DeepEqual(have, tc.want) {
			t.Errorf("Error: expected %v, found %v", tc.want, have)
		}
	}
}

type testStruct struct{}

type testStructWrongMethod struct{}

func (ts *testStructWrongMethod) ToString() {}

type testStructMethodErr struct{}

func (ts *testStructMethodErr) ToString() (int, int, error) { return 0, 0, errors.New("Test") }

func TestLogicalExpressionWalkError(t *testing.T) {
	cases := []struct {
		have *LogicalExpression
	}{
		{&LogicalExpression{expressions: []interface{}{&testStruct{}}}},
		{&LogicalExpression{expressions: []interface{}{&testStructWrongMethod{}}}},
		{&LogicalExpression{expressions: []interface{}{&testStructMethodErr{}}}},
	}

	for _, tc := range cases {
		_, _, err := tc.have.ToString()
		if err == nil {
			t.Error("Expected error found nil")
		}
	}
}
