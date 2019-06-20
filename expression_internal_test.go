package fury

import (
	"testing"
)

func TestLogicalExpressionWalk(t *testing.T) {
	cases := []struct {
		logicalOperator string
		expression1     interface{}
		expression2     interface{}
		want            string
	}{
		{"OR", IsEqualsTo("key", 1), nil, "key = 1"},
		{"OR", IsEqualsTo("key", "1"), nil, "key = '1'"},
		{"OR", IsEqualsTo("key", 1), IsEqualsTo("key", 2), "(key = 1 OR key = 2)"},
		{"OR", IsEqualsTo("key", 1), And(IsEqualsTo("key", 2)), "(key = 1 OR key = 2)"},
		{"OR", IsEqualsTo("key", 1), And(IsEqualsTo("key", 2), IsEqualsTo("key", 3)), "(key = 1 OR (key = 2 AND key = 3))"},
	}

	for _, tc := range cases {
		have, _ := newLogicalExpression(tc.logicalOperator, tc.expression1, tc.expression2).ToString()
		if have != tc.want {
			t.Errorf("Error expected %s, found %s", tc.want, have)
		}
	}
}
