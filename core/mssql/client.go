package CoreMsSQL

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	"github.com/jmoiron/sqlx"
	"time"
)

type Client struct {
	//DB 数据库对象
	DB *sqlx.DB
	//OpenEncrypt &encrypt=disable，官方握手BUG个别sqlserver下会造成TLS链接失败，启动后可规避该问题
	OpenEncrypt bool
}

func (t *Client) Init(dataSource, database, user, password string, port int) (err error) {
	dbDNS := fmt.Sprint("server=", dataSource, ";port=", port, ";user id=", user, ";password=", password, ";database=", database)
	if t.OpenEncrypt {
		dbDNS = fmt.Sprint(dbDNS, ";encrypt=disable")
	}
	CoreLog.Info("connect mssql dns: ", dbDNS)
	t.DB, err = sqlx.Open("mssql", dbDNS)
	if err != nil {
		return
	}
	//设置超时时间
	t.DB.SetConnMaxLifetime(60 * time.Second)
	return
}

// Get using this DB.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (t *Client) Get(dest interface{}, query string, args ...interface{}) error {
	return t.DB.Get(dest, query, args...)
}

// Select using this DB.
// Any placeholder parameters are replaced with supplied args.
func (t *Client) Select(dest interface{}, query string, args ...interface{}) error {
	return t.DB.Select(dest, query, args...)
}

// MustBegin starts a transaction, and panics on error.  Returns an *sqlx.Tx instead
// of an *sql.Tx.
func (t *Client) MustBegin() *sqlx.Tx {
	tx, err := t.DB.Beginx()
	if err != nil {
		panic(err)
	}
	return tx
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
//
// Exec uses context.Background internally; to specify the context, use
// ExecContext.
func (t *Client) Exec(query string, args ...any) (sql.Result, error) {
	return t.DB.Exec(query, args...)
}

// PrepareNamed returns an sqlx.NamedStmt
func (t *Client) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	return t.DB.PrepareNamed(query)
}

// NamedExec using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (t *Client) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return t.DB.NamedExec(query, arg)
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
//
// Ping uses context.Background internally; to specify the context, use
// PingContext.
func (t *Client) Ping() error {
	return t.DB.Ping()
}
