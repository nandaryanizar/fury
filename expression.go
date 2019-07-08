package fury

import (
	"errors"
	"fmt"
	"reflect"
)

// Expression struct to store query expression
type Expression struct {
	operator string
	operand1 string
	operand2 interface{}
}

// ToString method convert Expression struct to string and slice of arguments
func (e *Expression) ToString() (string, []interface{}, error) {
	if e.operator == "" || e.operand1 == "" || e.operand2 == nil {
		return "", nil, errors.New("Error creating expression: missing operator or operand")
	}
	args := []interface{}{e.operand2}

	return fmt.Sprintf("%s %s ?", e.operand1, e.operator), args, nil
}

// newExpression as factory function for Expression struct
func newExpression(operator, operand1 string, operand2 interface{}) *Expression {
	return &Expression{
		operator: operator,
		operand1: operand1,
		operand2: operand2,
	}
}

// IsGreaterThan expression
// 	This function will generate expression equivalent to 'operand1 > operand2'
func IsGreaterThan(operand1 string, operand2 interface{}) *Expression {
	return newExpression(">", operand1, operand2)
}

// IsGreaterThanOrEqualsTo expression
// 	This function will generate expression equivalent to 'operand1 >= operand2'
func IsGreaterThanOrEqualsTo(operand1 string, operand2 interface{}) *Expression {
	return newExpression(">=", operand1, operand2)
}

// IsLessThan expression
// 	This function will generate expression equivalent to 'operand1 < operand2'
func IsLessThan(operand1 string, operand2 interface{}) *Expression {
	return newExpression("<", operand1, operand2)
}

// IsLessThanOrEqualsTo expression
// 	This function will generate expression equivalent to 'operand1 <= operand2'
func IsLessThanOrEqualsTo(operand1 string, operand2 interface{}) *Expression {
	return newExpression("<=", operand1, operand2)
}

// IsEqualsTo expression
// 	This function will generate expression equivalent to 'operand1 = operand2'
func IsEqualsTo(operand1 string, operand2 interface{}) *Expression {
	return newExpression("=", operand1, operand2)
}

// IsNotEqualsTo expression
// 	This function will generate expression equivalent to 'operand1 <> operand2'
func IsNotEqualsTo(operand1 string, operand2 interface{}) *Expression {
	return newExpression("<>", operand1, operand2)
}

// LogicalExpression struct to store expression with logical condition as tree
type LogicalExpression struct {
	logicalOperator string
	expressions     []interface{}
}

// ToString method convert the LogicalExpression struct to string and slice of arguments
func (le *LogicalExpression) ToString() (string, []interface{}, error) {
	if len(le.expressions) == 0 {
		return "", nil, nil
	}

	out, args, err := walk(le)

	if err != nil {
		return "", nil, err
	}

	return out, args, nil
}

// Traverse LogicalExpression struct to convert it to string
func walk(lExp *LogicalExpression) (string, []interface{}, error) {
	out := ""
	args := []interface{}{}

	countNonNil := 0
	for i, val := range lExp.expressions {
		if val == nil {
			continue
		}

		countNonNil++
		if i > 0 {
			out += fmt.Sprintf(" %s ", lExp.logicalOperator)
		}

		if str, ok := val.(string); ok {
			out += str
			continue
		}

		exp := reflect.ValueOf(val)
		methodName := "ToString"

		toStringMethod := exp.MethodByName(methodName)
		if !toStringMethod.IsValid() {
			return "", nil, fmt.Errorf("Error: Couldn't find method %s in interface %v", methodName, exp)
		}

		ret := toStringMethod.Call([]reflect.Value{})
		// fmt.Println(ret)
		if len(ret) != 3 {
			return "", nil, fmt.Errorf("Error: Insufficient return value, expected 3 found %d", len(ret))
		}

		if err, ok := ret[2].Interface().(error); ok && err != nil {
			return "", nil, err
		}

		if retArgs, ok := ret[1].Interface().([]interface{}); ok {
			args = append(args, retArgs...)
		}

		if retStr, ok := ret[0].Interface().(string); ok {
			out += retStr
		}
	}

	if countNonNil > 1 {
		out = fmt.Sprintf("(%s)", out)
	}

	return out, args, nil
}

// newLogicalExpression is the factory function for LogicalExpression struct
func newLogicalExpression(logicalOp string, operands ...interface{}) *LogicalExpression {
	return &LogicalExpression{
		logicalOperator: logicalOp,
		expressions:     operands,
	}
}

// And expression
// 	Return logical expression with AND operator, equivalent to 'operands[0] AND operands[1] AND ...'
func And(operands ...interface{}) *LogicalExpression {
	return newLogicalExpression("AND", operands...)
}

// Or expression
// 	Return logical expression with OR operator, equivalent to 'operands[0] OR operands[1] OR ...'
func Or(operands ...interface{}) *LogicalExpression {
	return newLogicalExpression("OR", operands...)
}
