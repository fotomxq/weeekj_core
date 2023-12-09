package CoreSQL2

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// Get using this DB.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (t *ClientCtx) Get(dest interface{}, query string, args ...interface{}) error {
	switch t.client.DB.mode {
	case modePostgresql:
		return t.client.DB.dbPostgresql.Get(dest, query, args...)
	case modeMsSQL:
		return t.client.DB.dbMssql.Get(dest, query, args...)
	default:
		return t.client.DB.sqlDB.Get(dest, query, args...)
	}
}

// Select using this DB.
// Any placeholder parameters are replaced with supplied args.
func (t *ClientCtx) Select(dest interface{}, query string, args ...interface{}) error {
	switch t.client.DB.mode {
	case modePostgresql:
		return t.client.DB.dbPostgresql.Select(dest, query, args...)
	case modeMsSQL:
		return t.client.DB.dbMssql.Select(dest, query, args...)
	default:
		return t.client.DB.sqlDB.Select(dest, query, args...)
	}
}

// MustBegin starts a transaction, and panics on error.  Returns an *sqlx.Tx instead
// of an *sql.Tx.
func (t *ClientCtx) MustBegin() *sqlx.Tx {
	switch t.client.DB.mode {
	case modePostgresql:
		return t.client.DB.dbPostgresql.MustBegin()
	case modeMsSQL:
		return t.client.DB.dbMssql.MustBegin()
	default:
		tx, _ := t.client.DB.sqlDB.Beginx()
		return tx
	}
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
//
// Exec uses context.Background internally; to specify the context, use
// ExecContext.
func (t *ClientCtx) Exec(query string, args ...any) (sql.Result, error) {
	switch t.client.DB.mode {
	case modePostgresql:
		return t.client.DB.dbPostgresql.Exec(query, args...)
	case modeMsSQL:
		return t.client.DB.dbMssql.Exec(query, args...)
	default:
		return t.client.DB.sqlDB.Exec(query, args...)
	}
}

// PrepareNamed returns an sqlx.NamedStmt
func (t *ClientCtx) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	switch t.client.DB.mode {
	case modePostgresql:
		return t.client.DB.dbPostgresql.PrepareNamed(query)
	case modeMsSQL:
		return t.client.DB.dbMssql.PrepareNamed(query)
	default:
		return t.client.DB.sqlDB.PrepareNamed(query)
	}
}

// NamedExec using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (t *ClientCtx) NamedExec(query string, arg interface{}) (sql.Result, error) {
	switch t.client.DB.mode {
	case modePostgresql:
		return t.client.DB.dbPostgresql.NamedExec(query, arg)
	case modeMsSQL:
		return t.client.DB.dbMssql.NamedExec(query, arg)
	default:
		return t.client.DB.sqlDB.NamedExec(query, arg)
	}
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
//
// Ping uses context.Background internally; to specify the context, use
// PingContext.
func (t *ClientCtx) Ping() error {
	switch t.client.DB.mode {
	case modePostgresql:
		return t.client.DB.dbPostgresql.Ping()
	case modeMsSQL:
		return t.client.DB.dbMssql.Ping()
	default:
		return t.client.DB.sqlDB.Ping()
	}
}
