//go:generate go-enum -f=$GOFILE --nocase

package sql

import "strings"

// Driver is an enumeration of drivers.
/*
ENUM(
None // not existed
Postgres
Mysql
Sqlite3
Mssql
Redshift
Tidb
Clickhouse
)
*/
type Driver int32

func (x Driver) LowerString() string {
	return strings.ToLower(x.String())
}
