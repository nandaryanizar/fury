package fury_test

import (
	"reflect"
	"testing"

	"github.com/nandaryanizar/fury"
)

func TestInitializeDefault(t *testing.T) {
	cases := []struct {
		have string
		want *fury.Configuration
	}{
		{
			"database.yaml",
			&fury.Configuration{
				Host:            "postgres_db",
				Port:            "5432",
				Username:        "postgres",
				Password:        "pgadmin123",
				DBName:          "testdb",
				SSLMode:         false,
				MaxRetries:      1,
				ConnMaxLifetime: 0,
				MaxOpenConns:    0,
				MaxIdleConns:    2,
			},
		},
	}

	for _, tc := range cases {
		config, err := fury.LoadConfiguration(tc.have)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(config, tc.want) {
			t.Errorf("Error: expected %v, found %v", tc.want, config)
		}
	}
}
