package CorePostgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type Client struct {
	//DB 数据库对象
	DB *sqlx.DB
	//InstallDir 配置文件默认路径
	InstallDir string
	//最大连接数量
	MaxConnect int
	//连接超时时间秒
	ConnectExpireSec int
}

// Init 初始化
// eg url: host=%s port=%d user=%s password=%s dbname=%s sslmode=disable
func (t *Client) Init(url string, installDir string, timeZone string) (err error) {
	t.DB, err = sqlx.Connect("postgres", url)
	if err != nil {
		return
	}
	_, err = t.DB.Exec(fmt.Sprint("set time zone \"", timeZone, "\";"))
	if err != nil {
		err = errors.New("init exec sql, " + err.Error())
	}
	//设置超时时间
	if t.ConnectExpireSec < 1 {
		t.ConnectExpireSec = 10
	}
	t.DB.SetConnMaxLifetime(time.Duration(t.ConnectExpireSec) * time.Second)
	if t.MaxConnect < 1 {
		t.MaxConnect = 0
	}
	t.DB.SetMaxOpenConns(t.MaxConnect)
	//设置安装包目录
	if installDir != "" {
		t.InstallDir = installDir
	}
	//反馈
	return
}

// Get using this DB.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (t *Client) Get(dest interface{}, query string, args ...interface{}) (err error) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
			return
		}
	}()
	return t.DB.Get(dest, query, args...)
}

// Select using this DB.
// Any placeholder parameters are replaced with supplied args.
func (t *Client) Select(dest interface{}, query string, args ...interface{}) (err error) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
			return
		}
	}()
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
func (t *Client) Exec(query string, args ...any) (s sql.Result, err error) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
			return
		}
	}()
	return t.DB.Exec(query, args...)
}

// PrepareNamed returns an sqlx.NamedStmt
func (t *Client) PrepareNamed(query string) (s *sqlx.NamedStmt, err error) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
			return
		}
	}()
	return t.DB.PrepareNamed(query)
}

// NamedExec using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (t *Client) NamedExec(query string, arg interface{}) (s sql.Result, err error) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
			return
		}
	}()
	return t.DB.NamedExec(query, arg)
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
//
// Ping uses context.Background internally; to specify the context, use
// PingContext.
func (t *Client) Ping() (err error) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
			return
		}
	}()
	return t.DB.Ping()
}
