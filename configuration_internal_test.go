package fury

import (
	"reflect"
	"testing"
	"time"
)

func TestInitializeDefault(t *testing.T) {
	cases := []struct{ have, want *Configuration }{
		{
			&Configuration{
				Host:     "localhost",
				Port:     "5432",
				Username: "admin",
				Password: "admin",
			},
			&Configuration{
				Host:            "localhost",
				Port:            "5432",
				Username:        "admin",
				Password:        "admin",
				MaxRetries:      1,
				ConnMaxLifetime: 0, // No maximum lifetime
				MaxOpenConns:    0, // Unlimited open connections
				MaxIdleConns:    2,
			},
		},
		{
			&Configuration{
				Host:            "localhost",
				Port:            "5432",
				Username:        "admin",
				Password:        "admin",
				SSLMode:         false,
				MaxRetries:      -1,
				ConnMaxLifetime: -1,
				MaxOpenConns:    -1,
				MaxIdleConns:    -1,
			},
			&Configuration{
				Host:            "localhost",
				Port:            "5432",
				Username:        "admin",
				Password:        "admin",
				SSLMode:         false,
				MaxRetries:      1,
				ConnMaxLifetime: 0, // No maximum lifetime
				MaxOpenConns:    0, // Unlimited open connections
				MaxIdleConns:    2,
			},
		},
		{
			&Configuration{
				Host:            "localhost",
				Port:            "5432",
				Username:        "admin",
				Password:        "admin",
				SSLMode:         true,
				MaxRetries:      3,
				ConnMaxLifetime: time.Hour,
				MaxOpenConns:    10,
				MaxIdleConns:    5,
			},
			&Configuration{
				Host:            "localhost",
				Port:            "5432",
				Username:        "admin",
				Password:        "admin",
				SSLMode:         true,
				MaxRetries:      3,
				ConnMaxLifetime: time.Hour,
				MaxOpenConns:    10,
				MaxIdleConns:    5,
			},
		},
	}

	for _, tc := range cases {
		tc.have.initialize()
		if !reflect.DeepEqual(tc.have, tc.want) {
			t.Errorf("Error expected %v, found %v", tc.want, tc.have)
		}
	}
}
