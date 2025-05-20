package main

import (
	"database/sql"
	"testing"
)

func Test_openDB(t *testing.T) {
	type args struct {
		dsn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestNormalConnection",
			args: args{
				dsn: "postgres://postgres:postgres@localhost:5432/test_snippetbox?sslmode=disable",
			},
			wantErr: false,
		},
		{
			name: "TestInvalidConnection",
			args: args{
				dsn: "postgres://postgres:postgres@localhost:8080/test_snippetbox?sslmode=disable",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := openDB(tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("openDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && db != nil {
				defer func(db *sql.DB) {
					err := db.Close()
					if err != nil {
						t.Fatal(err)
					}
				}(db)
			}
		})
	}
}
