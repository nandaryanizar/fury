package fury

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

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

// LoadConfiguration load environment variable and create new configuration struct based on the variable
func LoadConfiguration(fileName string) (*Configuration, error) {
	if err := godotenv.Load(fileName); err != nil {
		return nil, fmt.Errorf("File with name %s doesn't exist", fileName)
	}

	username, err := getRequiredEnv("DATABASE_USERNAME")
	if err != nil {
		return nil, err
	}

	password, err := getRequiredEnv("DATABASE_PASSWORD")
	if err != nil {
		return nil, err
	}

	host, err := getRequiredEnv("DATABASE_HOST")
	if err != nil {
		return nil, err
	}

	port, err := getRequiredEnv("DATABASE_PORT")
	if err != nil {
		return nil, err
	}

	dbname, err := getRequiredEnv("DATABASE_NAME")
	if err != nil {
		return nil, err
	}

	return &Configuration{
		Username:        username,
		Password:        password,
		Host:            host,
		Port:            port,
		DBName:          dbname,
		SSLMode:         getEnvAsBool("DATABASE_SSLMODE", false),
		MaxRetries:      getEnvAsInt("DATABASE_MAXRETRIES", 1),
		MaxIdleConns:    getEnvAsInt("DATABASE_MAXIDLECONNS", 2),
		MaxOpenConns:    getEnvAsInt("DATABASE_MAXOPENCONNS", 0),
		ConnMaxLifetime: getEnvAsTimeDuration("DATABASE_CONNMAXXLIFETIME", 0),
	}, nil
}

// Lookup env variable and return default value if not exists
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// Get required string environment environment variable
func getRequiredEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("Error: cannot find environment variable %s", key)
}

// Lookup env variable and return default value as int if not exists
func getEnvAsInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}
	}
	return defaultVal
}

// Lookup env variable and return default value as bool if not exists
func getEnvAsTimeDuration(key string, defaultVal time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if val, err := strconv.ParseInt(value, 10, 64); err == nil {
			return time.Duration(val)
		}
	}
	return defaultVal
}

// Lookup env variable and return default value as bool if not exists
func getEnvAsBool(key string, defaultVal bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if val, err := strconv.ParseBool(value); err == nil {
			return val
		}
	}
	return defaultVal
}
