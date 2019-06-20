package fury

import "time"

// Configuration struct for database object. Consists of connection and driver configurations.
type Configuration struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  bool

	MaxRetries      int
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

func (c *Configuration) initialize() {
	if c.MaxRetries < 1 {
		c.MaxRetries = 1
	}

	if c.MaxIdleConns < 1 {
		c.MaxIdleConns = 2
	}

	if c.MaxOpenConns < 0 {
		c.MaxOpenConns = 0
	}

	if c.ConnMaxLifetime < 0 {
		c.ConnMaxLifetime = 0
	}
}
