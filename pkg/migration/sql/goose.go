package sql

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func NewGooseMigration(driver Driver, dir, dbString string) Migration {
	conn, err := goose.OpenDBWithDriver(driver.LowerString(), dbString)
	if err != nil {
		panic("sql connection failed")
	}

	return &gooseMigration{
		conn: conn,
		dir:  dir,
	}
}

type gooseMigration struct {
	dir  string
	conn *sql.DB
}

func (m *gooseMigration) Close() error {
	return m.conn.Close()
}

func (m *gooseMigration) Up() error {
	return goose.Up(m.conn, m.dir, goose.WithAllowMissing())
}

func (m *gooseMigration) Reset() error {
	return goose.Reset(m.conn, m.dir)
}
