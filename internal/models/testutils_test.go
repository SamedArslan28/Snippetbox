package models

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Connect using the pgx driver
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/test_snippetbox?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Run setup SQL script
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(script)); err != nil {
		t.Fatal(err)
	}

	// Automatically clean up after test
	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(string(script)); err != nil {
			t.Fatal(err)
		}
		db.Close()
	})

	return db
}
