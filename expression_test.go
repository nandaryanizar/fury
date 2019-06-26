package fury_test

import (
	"testing"

	"github.com/nandaryanizar/fury"
)

func TestExpressionToStringIfEmpty(t *testing.T) {
	want := ""
	have, _, _ := (&fury.Expression{}).ToString()
	if have != want {
		t.Errorf("Error: expected %v, found %v", want, have)
	}
}

func TestExpressionToStringIfMissingOp(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
	}{
		{"", "val"},
		{"key", nil},
	}

	for _, tc := range cases {
		_, _, err := fury.IsEqualsTo(tc.operand1, tc.operand2).ToString()
		if err == nil {
			t.Error("Expected error found nil")
		}
	}
}

func TestIsGreaterThanExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key > ?"},
		{"key", "2", "key > ?"},
	}

	for _, tc := range cases {
		have, _, _ := fury.IsGreaterThan(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, have)
		}
	}
}

func TestIsGreaterThanOrEqualsToExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key >= ?"},
		{"key", "2", "key >= ?"},
	}

	for _, tc := range cases {
		have, _, _ := fury.IsGreaterThanOrEqualsTo(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, have)
		}
	}
}

func TestIsLessThanExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key < ?"},
		{"key", "2", "key < ?"},
	}

	for _, tc := range cases {
		have, _, _ := fury.IsLessThan(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, have)
		}
	}
}

func TestIsLessThanOrEqualsToExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key <= ?"},
		{"key", "2", "key <= ?"},
	}

	for _, tc := range cases {
		have, _, _ := fury.IsLessThanOrEqualsTo(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, have)
		}
	}
}

func TestIsEqualsToExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key = ?"},
		{"key", "2", "key = ?"},
	}

	for _, tc := range cases {
		have, _, _ := fury.IsEqualsTo(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, have)
		}
	}
}

func TestIsNotEqualsToExpression(t *testing.T) {
	cases := []struct {
		operand1 string
		operand2 interface{}
		want     string
	}{
		{"key", 2, "key <> ?"},
		{"key", "2", "key <> ?"},
	}

	for _, tc := range cases {
		have, _, _ := fury.IsNotEqualsTo(tc.operand1, tc.operand2).ToString()
		if have != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, have)
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
		have, _, _ := fury.And(tc.expressions).ToString()
		if have != tc.want {
			t.Errorf("Error: expected %s, found %s", tc.want, have)
		}
	}
}

func TestLogicalExpressionToString(t *testing.T) {
	cases := []struct {
		expression1 interface{}
		expression2 interface{}
		want        string
	}{
		{fury.IsEqualsTo("key", "1"), nil, "key = ?"},
		{fury.IsEqualsTo("key", 1), fury.IsEqualsTo("key", 2), "(key = ? OR key = ?)"},
		{fury.IsEqualsTo("key", 1), fury.And(fury.IsEqualsTo("key", 2)), "(key = ? OR key = ?)"},
		{fury.IsEqualsTo("key", 1), fury.And(fury.IsEqualsTo("key", 2), fury.IsEqualsTo("key", 3)), "(key = ? OR (key = ? AND key = ?))"},
	}

	for _, tc := range cases {
		have, _, _ := fury.Or(tc.expression1, tc.expression2).ToString()
		if have != tc.want {
			t.Errorf("Error: expected %s, found %s", tc.want, have)
		}
	}
}
