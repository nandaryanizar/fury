package fury_test

import (
	"testing"

	"github.com/nandaryanizar/fury"
)

func TestExpressionToStringIfEmpty(t *testing.T) {
	want := ""
	have, _ := (&fury.Expression{}).ToString()
	if have != want {
		t.Errorf("Error expected %v, found %v", want, have)
	}
}

func TestExpressionToStringIfMissingOp(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"", "val", "key > val"},
		{"key", nil, "key > val"},
	}

	for _, tc := range cases {
		_, err := fury.IsEqualsTo(tc.operand1, tc.operand2).ToString()
		if err == nil {
			t.Errorf("Error expected found nil")
		}
	}
}

func TestIsGreaterThanExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key > 2"},
		{"key", "val", "key > val"},
	}

	for _, tc := range cases {
		have, _ := fury.IsGreaterThan(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error expected %v, found %v", tc.want, have)
		}
	}
}

func TestIsGreaterThanOrEqualsToExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key >= 2"},
		{"key", "val", "key >= val"},
	}

	for _, tc := range cases {
		have, _ := fury.IsGreaterThanOrEqualsTo(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error expected %v, found %v", tc.want, have)
		}
	}
}

func TestIsLessThanExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key < 2"},
		{"key", "val", "key < val"},
	}

	for _, tc := range cases {
		have, _ := fury.IsLessThan(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error expected %v, found %v", tc.want, have)
		}
	}
}

func TestIsLessThanOrEqualsToExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key <= 2"},
		{"key", "val", "key <= val"},
	}

	for _, tc := range cases {
		have, _ := fury.IsLessThanOrEqualsTo(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error expected %v, found %v", tc.want, have)
		}
	}
}

func TestIsEqualsToExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key = 2"},
		{"key", "val", "key = val"},
	}

	for _, tc := range cases {
		have, _ := fury.IsEqualsTo(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error expected %v, found %v", tc.want, have)
		}
	}
}

func TestIsNotEqualsToExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key <> 2"},
		{"key", "val", "key <> val"},
	}

	for _, tc := range cases {
		have, _ := fury.IsNotEqualsTo(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error expected %v, found %v", tc.want, have)
		}
	}
}

func TestLogicalExpressionToStringEmptyExp(t *testing.T) {
	cases := []struct {
		expressions interface{}
		want        string
	}{
		{nil, ""},
	}

	for _, tc := range cases {
		have, _ := fury.And(tc.expressions).ToString()
		if have != tc.want {
			t.Errorf("Error expected %s, found %s", tc.want, have)
		}
	}
}

func TestLogicalExpressionToString(t *testing.T) {
	cases := []struct {
		expression1 interface{}
		expression2 interface{}
		want        string
	}{
		{fury.IsEqualsTo("key", 1), nil, "key = 1"},
		{fury.IsEqualsTo("key", 1), fury.IsEqualsTo("key", 2), "key = 1 OR key = 2"},
		{fury.IsEqualsTo("key", 1), fury.And(fury.IsEqualsTo("key", 2)), "key = 1 OR key = 2"},
		{fury.IsEqualsTo("key", 1), fury.And(fury.IsEqualsTo("key", 2), fury.IsEqualsTo("key", 3)), "key = 1 OR (key = 2 AND key = 3)"},
	}

	for _, tc := range cases {
		have, _ := fury.Or(tc.expression1, tc.expression2).ToString()
		if have != tc.want {
			t.Errorf("Error expected %s, found %s", tc.want, have)
		}
	}
}

// func TestLogicalExpressionToStringError(t *testing.T) {
// 	cases := []struct {
// 		expressions interface{}
// 	}{
// 		{[]*fury.Expression{fury.IsEqualsTo("key", 1), fury.IsEqualsTo("key", nil)}},
// 		{[]*fury.Expression{fury.IsEqualsTo("key", 1), fury.IsEqualsTo("key", 2)}},
// 	}

// 	for _, tc := range cases {
// 		_, err := fury.And(tc.expressions).ToString()
// 		if err == nil {
// 			t.Errorf("Error expected")
// 		}
// 	}
// }

// var logicalExpCases = []struct {
// 	logicalOperator string
// 	expressions     interface{}
// 	want            *fury.LogicalExpression
// }{
// 	{"AND", []*fury.Expression{fury.IsEqualsTo("key", 1), fury.IsEqualsTo("key", 2)}},
// }

// func TestLogicalExpressionToString(t *testing.T) {
// 	for _, tc := range expCases {
// 		have, _ := fury.newExpression(tc.operator, tc.operand1, tc.operand2).ToString()
// 		if have != tc.want {
// 			t.Errorf("Error expected %v, found %v", tc.want, have)
// 		}
// 	}
// }
