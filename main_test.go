package fury_test

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/nandaryanizar/fury"
)

var db *fury.DB

type Account struct {
	UserID    int `fury:"primary_key,auto_increment"`
	Username  string
	Password  string
	Email     string
	CreatedOn time.Time
	LastLogin time.Time
}

func TestMain(m *testing.M) {
	var err error
	db, err = fury.Connect(&fury.Configuration{
		Username: "postgres",
		Password: "pgadmin123",
		Host:     "postgres_db",
		Port:     "5432",
		DBName:   "testdb",
	})

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	runMigrationAndSeeder()

	flag.Parse()
	exitCode := m.Run()

	dropMigration()

	db.Close()

	// Exit
	os.Exit(exitCode)
}

func runMigrationAndSeeder() {
	dropMigration()

	query, err := ioutil.ReadFile("./migrations/create_account_table.sql")
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}

	seed, err := ioutil.ReadFile("./migrations/seed_account_table.sql")
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(string(seed)); err != nil {
		panic(err)
	}
}

func dropMigration() {
	query, err := ioutil.ReadFile("./migrations/drop_account_table.sql")
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
}
