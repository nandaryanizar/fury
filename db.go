package fury

import (
	"time"
)

// DB object consists of DB connection pool and configuration
// 	Use Connect(config *Configuration) to create new instance of this struct.
type DB struct {
	ConnectionPooler
	config *Configuration
	query  *Query
}

// Connect to database, instantiate DB struct for querying to database.
func Connect(configFileName string) (*DB, error) {
	config, err := LoadConfiguration(configFileName)
	if err != nil {
		return nil, err
	}

	newConnPool, err := NewConnectionPool(config)
	if err != nil {
		return nil, err
	}

	db := &DB{
		ConnectionPooler: newConnPool,
		config:           config,
	}

	return db, nil
}

// ConnectMock to database, instantiate DB struct for querying to database.
func ConnectMock(mockPool ConnectionPooler) (*DB, error) {
	config := &Configuration{
		MaxRetries:      2,
		ConnMaxLifetime: time.Hour,
		MaxIdleConns:    0,
		MaxOpenConns:    0,
	}

	db := &DB{
		ConnectionPooler: mockPool,
		config:           config,
	}

	return db, nil
}

func (db *DB) clone(model interface{}) (*DB, error) {
	q, err := NewQuery(model)
	if err != nil {
		return nil, err
	}

	return &DB{
		ConnectionPooler: db.ConnectionPooler,
		config:           db.config,
		query:            q,
	}, nil
}

// First method return first record ordered by primary key
func (db *DB) First(model interface{}, opts ...QueryOption) error {
	opts = append(opts, Limit(1))
	return db.Find(model, opts...)
}

// Find method return all record queried with specified conditions
func (db *DB) Find(model interface{}, opts ...QueryOption) error {
	newDB, err := db.clone(model)
	if err != nil {
		return err
	}

	for _, opt := range opts {
		_, err := opt(newDB.query)
		if err != nil {
			return err
		}
	}

	return newDB.executeSelectQuery()
}

// Insert query method
func (db *DB) Insert(model interface{}, opts ...QueryOption) error {
	newDB, err := db.clone(model)
	if err != nil {
		return err
	}

	for _, opt := range opts {
		_, err := opt(newDB.query)
		if err != nil {
			return err
		}
	}

	return newDB.executeInsertQuery()
}

// Update query method
func (db *DB) Update(model interface{}, opts ...QueryOption) error {
	newDB, err := db.clone(model)
	if err != nil {
		return err
	}

	for _, opt := range opts {
		_, err := opt(newDB.query)
		if err != nil {
			return err
		}
	}

	return newDB.executeUpdateQuery()
}

// Delete query method
func (db *DB) Delete(model interface{}, opts ...QueryOption) error {
	newDB, err := db.clone(model)
	if err != nil {
		return err
	}

	for _, opt := range opts {
		_, err := opt(newDB.query)
		if err != nil {
			return err
		}
	}

	return newDB.executeDeleteQuery()
}

func (db *DB) executeSelectQuery() error {
	if err := db.query.prepareSelectQuery(); err != nil {
		return err
	}

	rows, err := db.Query(db.query.SQL, db.query.args...)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		mPtr, err := db.query.nextOrCreateModel()
		if err != nil {
			return err
		}

		if mPtr == nil {
			break
		}

		pointers := db.query.modelPtr.GetScanPtrByColumnNames(columns)
		if err := rows.Scan(pointers...); err != nil {
			return err
		}
	}
	rows.Close()

	return nil
}

func (db *DB) executeInsertQuery() error {
	for db.query.nextModel() != nil {
		if err := db.query.prepareInsertQuery(); err != nil {
			return err
		}

		_, err := db.Exec(db.query.SQL, db.query.args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) executeUpdateQuery() error {
	for db.query.nextModel() != nil {
		if err := db.query.prepareUpdateQuery(); err != nil {
			return err
		}

		_, err := db.Exec(db.query.SQL, db.query.args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) executeDeleteQuery() error {
	for db.query.nextModel() != nil {
		if err := db.query.prepareDeleteQuery(); err != nil {
			return err
		}

		_, err := db.Exec(db.query.SQL, db.query.args...)
		if err != nil {
			return err
		}
	}

	return nil
}
