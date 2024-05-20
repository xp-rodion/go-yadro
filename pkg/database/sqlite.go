package database

import (
	"github.com/jmoiron/sqlx"
)

type SQLite struct {
	DSN string
	DB  sqlx.DB
}

func (d *SQLite) Init(DSN string) {
	d.DSN = DSN
}

func (d *SQLite) Open() bool {
	db, err := sqlx.Open("sqlite3", d.DSN)

	d.DB = *db
	
	if err != nil {
		return false
	}

	return true
}
