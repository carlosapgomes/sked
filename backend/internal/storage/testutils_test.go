package storage_test

import (
	"database/sql"
	"io/ioutil"
	"testing"
	// // blank inport posgres lib
	// _ "github.com/lib/pq"
)

// Shamelessly copied from "Let'sGo!" book
// by Alex Edwards (https://www.alexedwards.net/#book)
// Establish a sql.DB connection pool for our test database.
func newTestDB(t *testing.T) (*sql.DB, func()) {
	// Initialise a new connection pool
	db, err := sql.Open("postgres", "postgres://test_sked:test_sked@localhost/test_sked?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script from file and execute the statements.
	script, err := ioutil.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Return the connection pool and an anonymous function which reads and
	// executes the teardown script, and closes the connection pool. We can
	// assign this anonymous function and call it later once our test has
	// completed.
	return db, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	}
}
