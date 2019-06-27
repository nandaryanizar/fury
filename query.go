package fury

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/nandaryanizar/fury/model"
)

// QueryOption is return type for every main query and execute method
type QueryOption func(q *Query) (*Query, error)

// Query base struct
type Query struct {
	SQL             string
	tableName       string
	models          []*model.Model
	columns         []interface{}
	scanTo          interface{}
	whereConditions []interface{}
	args            []interface{}
	limit           int
	offset          int
	groups          []interface{}
	orders          []interface{}
	useModelAsCond  bool
	modelPtr        *model.Model
	modelPtrCtr     int
}

// NewQuery return new Query literal
func NewQuery(modelInterface interface{}) (*Query, error) {
	m, mPtr, err := model.NewModels(modelInterface)
	if err != nil {
		return nil, err
	}

	q := &Query{
		models:         m,
		useModelAsCond: true,
		scanTo:         modelInterface,
		modelPtr:       mPtr,
		modelPtrCtr:    -1,
	}

	return q, nil
}

// Table specify which table to query
//	using this function means that query will not use model to specify query condition
func Table(tableName string) QueryOption {
	return func(q *Query) (*Query, error) {
		q.tableName = tableName
		q.useModelAsCond = false
		return q, nil
	}
}

// Where function is used to add new query condition
// 	Supported expression condition type: Expression, LogicalExpression, string
//  Use Expression and LogicalExpression for type safety
func Where(conditions interface{}) QueryOption {
	return func(q *Query) (*Query, error) {
		switch conditions.(type) {
		case *Expression, *LogicalExpression, string:
			q.whereConditions = append(q.whereConditions, conditions)
		default:
			return nil, errors.New("Unsupported expression conditions type")
		}

		return q, nil
	}
}

// Select function is used to specify columns in query
func Select(columns ...interface{}) QueryOption {
	return func(q *Query) (*Query, error) {
		q.columns = append(q.columns, columns...)
		return q, nil
	}
}

// Limit function is used to add limit query
func Limit(limit int) QueryOption {
	return func(q *Query) (*Query, error) {
		if limit < 0 {
			return nil, errors.New("Limit cannot be negative number")
		}
		q.limit = limit
		return q, nil
	}
}

// Offset function is used to add offset query
func Offset(offset int) QueryOption {
	return func(q *Query) (*Query, error) {
		if offset < 0 {
			return nil, errors.New("Offset cannot be negative number")
		}
		q.offset = offset
		return q, nil
	}
}

// GroupBy function is used to add group by query
func GroupBy(columns ...interface{}) QueryOption {
	return func(q *Query) (*Query, error) {
		q.groups = append(q.groups, columns...)
		return q, nil
	}
}

// OrderBy function is used to add order by query
func OrderBy(columns ...interface{}) QueryOption {
	return func(q *Query) (*Query, error) {
		q.orders = append(q.orders, columns...)
		return q, nil
	}
}

// Create new context
func (q *Query) clone() *Query {
	return &Query{
		SQL:             "",
		tableName:       q.tableName,
		models:          q.models,
		columns:         q.columns,
		scanTo:          q.scanTo,
		whereConditions: q.whereConditions,
		args:            []interface{}{},
		limit:           q.limit,
		offset:          q.offset,
		groups:          q.groups,
		orders:          q.orders,
		useModelAsCond:  q.useModelAsCond,
		modelPtr:        q.modelPtr,
		modelPtrCtr:     q.modelPtrCtr,
	}
}

// NextModel shift modelPtr to next model if available, if empty return nil and modelPtr not shifted
//	When model is struct or pointer to struct, then this method should return nil and modelPtr not shifted too
func (q *Query) nextModel() *model.Model {
	reflectVal := reflect.ValueOf(q.scanTo)
	if reflectVal.Kind() == reflect.Ptr {
		reflectVal = reflectVal.Elem()
	}

	if reflectVal.Kind() == reflect.Slice || reflectVal.Kind() == reflect.Struct {
		q.modelPtrCtr++
		if q.modelPtrCtr < len(q.models) {
			q.modelPtr = q.models[q.modelPtrCtr]
			return q.modelPtr
		}
		q.modelPtrCtr--
	}
	return nil
}

// NextOrCreateModel shift modelPtr to next model or create new and append it to slice
//	When model is struct or pointer to struct, then this method should return nil, modelPtr not shifted, and not creating new instance
func (q *Query) nextOrCreateModel() (*model.Model, error) {
	reflectVal := reflect.ValueOf(q.scanTo)
	if reflectVal.Kind() == reflect.Ptr {
		reflectVal = reflectVal.Elem()
	}

	if q.modelPtrCtr++; q.modelPtrCtr < len(q.models) {
		q.modelPtr = q.models[q.modelPtrCtr]
		return q.modelPtr, nil
	}

	if reflectVal.Kind() != reflect.Slice {
		q.modelPtrCtr--
		return nil, nil
	}

	rType := q.modelPtr.Type
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}

	m := reflect.New(rType)
	models, _, err := model.NewModels(m.Interface())
	if err != nil {
		return nil, err
	}

	q.models = append(q.models, models...)
	reflectVal.Set(reflect.Append(reflectVal, m))

	if q.modelPtrCtr < len(q.models) {
		q.modelPtr = q.models[q.modelPtrCtr]
	}

	return q.modelPtr, nil
}

func (q *Query) addPKWhereConditions() error {
	if q.modelPtr == nil {
		return nil
	}

	whereConds := []interface{}{}
	for _, f := range q.modelPtr.PrimaryKeys {
		if f.CheckIfZeroValue() {
			continue
		}

		val := f.Value
		if f.Value.Kind() == reflect.Ptr {
			val = f.Value.Elem()
		}

		key := fmt.Sprintf("%s.%s", q.modelPtr.Name, strings.ToLower(f.Properties.Name))
		whereConds = append(whereConds, IsEqualsTo(key, val.Interface()))
	}

	if len(whereConds) > 0 {
		q.whereConditions = append(q.whereConditions, And(whereConds...))
	}

	return nil
}

func (q *Query) addAllPKWhereConditions() error {
	if q.models == nil || len(q.models) < 1 {
		return nil
	}

	allWhereConds := []interface{}{}
	for i := 0; i < len(q.models); i++ {
		if q.models[i] == nil {
			continue
		}
		m := q.models[i]

		whereConds := []interface{}{}
		for _, f := range m.PrimaryKeys {
			if f.CheckIfZeroValue() {
				continue
			}

			val := f.Value
			if f.Value.Kind() == reflect.Ptr {
				val = f.Value.Elem()
			}

			key := fmt.Sprintf("%s.%s", m.Name, strings.ToLower(f.Properties.Name))
			whereConds = append(whereConds, IsEqualsTo(key, val.Interface()))
		}

		if len(whereConds) > 0 {
			allWhereConds = append(allWhereConds, And(whereConds...))
		}
	}

	if len(allWhereConds) > 0 {
		q.whereConditions = append(q.whereConditions, Or(allWhereConds...))
	}

	return nil
}

func (q *Query) prepareSelectColumn() (string, error) {
	out := ""
	for i, col := range q.columns {
		if colStr, ok := col.(string); ok {
			if i > 0 {
				out += ", "
			}
			out += colStr
		}
	}

	if len(q.columns) == 0 {
		tableName, err := q.getTableName()
		if err != nil {
			return "", err
		}

		out = fmt.Sprintf("%s.*", tableName)
	}

	return fmt.Sprintf("SELECT %s", out), nil
}

func (q *Query) getTableName() (string, error) {
	if q.tableName != "" {
		return q.tableName, nil
	}

	if q.modelPtr != nil {
		return q.modelPtr.Name, nil
	}

	return "", errors.New("Error: unspecified table name")
}

func (q *Query) prepareWhereQuery() (string, error) {
	out := ""
	for i, cond := range q.whereConditions {
		if i > 0 {
			out += " AND "
		}

		if exp, ok := cond.(*Expression); ok {
			expStr, args, err := exp.ToString()
			if err != nil {
				return "", err
			}

			q.args = append(q.args, args...)
			out += expStr
		} else if exp, ok := cond.(*LogicalExpression); ok {
			expStr, args, err := exp.ToString()
			if err != nil {
				return "", err
			}

			q.args = append(q.args, args...)
			out += expStr
		} else if exp, ok := cond.(string); ok {
			out += exp
		}
	}

	return fmt.Sprintf(" WHERE %s", out), nil
}

func (q *Query) prepareLimitOffsetQuery() string {
	out := ""

	if q.limit > 0 {
		out += fmt.Sprintf(" LIMIT %d", q.limit)
	}

	if q.offset > 0 {
		out += fmt.Sprintf(" OFFSET %d", q.offset)
	}

	return out
}

func (q *Query) prepareGroupByQuery() string {
	out := ""
	for i, col := range q.groups {
		if colStr, ok := col.(string); ok {
			if i > 0 {
				out += ", "
			}
			out += colStr
		}
	}

	if len(out) > 0 {
		out = fmt.Sprintf(" GROUP BY %s", out)
	}

	return out
}

func (q *Query) prepareOrderByQuery() string {
	out := ""
	for i, col := range q.orders {
		if colStr, ok := col.(string); ok {
			if i > 0 {
				out += ", "
			}
			out += colStr
		}
	}

	if len(out) > 0 {
		out = fmt.Sprintf(" ORDER BY %s", out)
	}

	return out
}

func (q *Query) replaceSQLPlaceholder() {
	for i := 0; i < len(q.args); i++ {
		q.SQL = strings.Replace(q.SQL, "?", fmt.Sprintf("$%d", i+1), 1)
	}
}

func (q *Query) prepareSelectQuery() error {
	selectColumn, err := q.prepareSelectColumn()
	if err != nil {
		return err
	}

	tableName, err := q.getTableName()
	if err != nil {
		return err
	}

	if q.useModelAsCond {
		if err := q.addAllPKWhereConditions(); err != nil {
			return err
		}
	}

	whereQuery, err := q.prepareWhereQuery()
	if err != nil {
		return err
	}

	limitOffsetQuery := q.prepareLimitOffsetQuery()
	groupByQuery := q.prepareGroupByQuery()
	orderByQuery := q.prepareOrderByQuery()

	q.SQL = fmt.Sprintf("%s FROM %s%s%s%s%s;", selectColumn, tableName, whereQuery, groupByQuery, orderByQuery, limitOffsetQuery)
	q.replaceSQLPlaceholder()

	return nil
}

func (q *Query) getColumnsNamesAndValues(includeAutoInc bool) ([]string, []interface{}) {
	return q.modelPtr.GetColumnNamesAndValues(includeAutoInc)
}

func (q *Query) prepareInsertQuery() error {
	query := q.clone()
	cols, args := query.getColumnsNamesAndValues(false)
	if len(cols) != len(args) {
		return errors.New("Columns and argument length not match")
	}
	query.args = append(query.args, args...)
	if len(cols) < 1 || len(args) < 1 {
		return errors.New("Columns or argument slice cannot be empty")
	}

	tableName, err := query.getTableName()
	if err != nil {
		return err
	}

	columnQuery := ""
	valueQuery := ""
	for i := 0; i < len(args); i++ {
		if i != 0 {
			columnQuery += ","
			valueQuery += ","
		}
		columnQuery += cols[i]
		valueQuery += fmt.Sprintf("$%d", i+1)
	}

	q.SQL = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s);", tableName, columnQuery, valueQuery)
	q.args = query.args

	return nil
}

func (q *Query) prepareUpdateQuery() error {
	query := q.clone()
	cols, args := query.getColumnsNamesAndValues(true)
	if len(cols) != len(args) {
		return errors.New("Columns and argument length not match")
	}
	query.args = append(args, query.args...)

	if len(cols) < 1 || len(args) < 1 {
		return errors.New("Columns or argument slice cannot be empty")
	}

	if query.useModelAsCond {
		if err := query.addPKWhereConditions(); err != nil {
			return err
		}
	}

	tableName, err := query.getTableName()
	if err != nil {
		return err
	}

	columnQuery := ""
	valueQuery := ""
	for i := 0; i < len(args); i++ {
		if i != 0 {
			columnQuery += ","
			valueQuery += ","
		}
		columnQuery += cols[i]
		valueQuery += "?"
	}

	whereQuery, err := query.prepareWhereQuery()
	if err != nil {
		return err
	}

	query.SQL = fmt.Sprintf("UPDATE %s SET (%s) = (%s)%s;", tableName, columnQuery, valueQuery, whereQuery)
	query.replaceSQLPlaceholder()

	q.SQL = query.SQL
	q.args = query.args

	return nil
}

func (q *Query) prepareDeleteQuery() error {
	query := q.clone()
	if query.useModelAsCond {
		if err := query.addPKWhereConditions(); err != nil {
			return err
		}
	}

	tableName, err := query.getTableName()
	if err != nil {
		return err
	}

	if len(query.whereConditions) < 1 {
		return errors.New("Unsupported delete without filter")
	}

	whereQuery, err := query.prepareWhereQuery()
	if err != nil {
		return err
	}

	query.SQL = fmt.Sprintf("DELETE FROM %s%s;", tableName, whereQuery)
	query.replaceSQLPlaceholder()

	q.SQL = query.SQL
	q.args = query.args

	return nil
}
