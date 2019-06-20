package fury

// QueryOption is return type for every main query and execute method
type QueryOption func(db *DB) (*DB, error)

// Query base struct
type Query struct {
	connectionPool
	tableName       interface{}
	columns         []interface{}
	whereConditions []interface{}
	limit           int
	offset          int
	groups          []interface{}
	orders          []interface{}
}

// First method return first record ordered by primary key
func (q *Query) First(out interface{}, opts ...QueryOption) error {
	return nil
}

// Find method return all record queried with specified conditions
func (q *Query) Find(out interface{}, opts ...QueryOption) error {
	return nil
}

// Count method will replace the query with number of rows returned
func (q *Query) Count(out interface{}, opts ...QueryOption) error {
	return nil
}

// Where function is used to add new query condition
func Where(conditions interface{}) QueryOption {
	return func(db *DB) (*DB, error) {
		return nil, nil
	}
}

// Select function is used to specify columns to be selected
func Select(columns ...interface{}) QueryOption {
	return func(db *DB) (*DB, error) {
		return nil, nil
	}
}

// Limit function is used to add limit query
func Limit(limit interface{}) QueryOption {
	return func(db *DB) (*DB, error) {
		return nil, nil
	}
}

// Offset function is used to add offset query
func Offset(limit interface{}) QueryOption {
	return func(db *DB) (*DB, error) {
		return nil, nil
	}
}

// GroupBy function is used to add limit query
func GroupBy(columns ...interface{}) QueryOption {
	return func(db *DB) (*DB, error) {
		return nil, nil
	}
}

// OrderBy function is used to add limit query
func OrderBy(columns ...interface{}) QueryOption {
	return func(db *DB) (*DB, error) {
		return nil, nil
	}
}
