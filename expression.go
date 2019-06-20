package fury

import "fmt"

// Expression struct to store query expression
type Expression struct {
	operator string
	operand1 string
	operand2 interface{}
}

// ToString method convert Expression struct to string
func (e *Expression) ToString() (string, error) {

	op2 := fmt.Sprintf("%v", e.operand2)
	if e.operator == "" || e.operand1 == "" || e.operand2 == nil || op2 == "" {
		return "", fmt.Errorf("Error creating expression: missing operator or operand")
	}

	return fmt.Sprintf("%s %s %s", e.operand1, e.operator, op2), nil
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

// ToString method convert the LogicalExpression struct to string
func (le *LogicalExpression) ToString() (out string, err error) {
	if len(le.expressions) == 0 {
		return "", nil
	}

	if err = walk(&out, le); err != nil {
		return "", err
	}

	return out, nil
}

func walk(out *string, lExp *LogicalExpression) error {
	// for _, in := range lExp.expressions {

	// }
	return nil
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
