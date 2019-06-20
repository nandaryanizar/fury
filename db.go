package fury

// DB object consists of DB connection pool and configuration
// 	Use Connect(config *Configuration) to create new instance of this struct.
type DB struct {
	connectionPool
	config *Configuration
	*Query
}

// Connect to database, instantiate DB struct for querying to database.
func Connect(config *Configuration) (*DB, error) {
	config.initialize()

	newConnPool, err := newConnectionPool(config)

	if err != nil {
		return nil, err
	}

	db := &DB{
		connectionPool: newConnPool,
		config:         config,
	}

	return db, nil
}

// Model function to set table from struct
func (db *DB) Model(model interface{}) {
	return
}
