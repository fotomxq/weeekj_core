package CoreSQL2

import (
	CoreMsSQL "github.com/fotomxq/weeekj_core/v5/core/mssql"
	CorePostgres "github.com/fotomxq/weeekj_core/v5/core/postgres"
	"github.com/jmoiron/sqlx"
)

//第二代sql生成工具

var (
	//模式常量
	modePostgresql = "postgresql"
	modeMsSQL      = "mssql"
	//OpenDebug debug
	OpenDebug = false
)

// SQLClient SQL操作控制器核心
// 初始化连接后需使用本结构构建控制对象，写入基础数据集
type SQLClient struct {
	//数据库模式
	// postgresql; mssql
	mode string
	//数据库句柄
	dbPostgresql *CorePostgres.Client
	dbMssql      *CoreMsSQL.Client
	sqlDB        *sqlx.DB
}

// InitSQLDB 初始化
func (t *SQLClient) InitSQLDB(mode string, mainDB *sqlx.DB) {
	t.mode = mode
	t.sqlDB = mainDB
}

// InitPostgresql 初始化
func (t *SQLClient) InitPostgresql(mainDB *CorePostgres.Client) {
	t.mode = modePostgresql
	t.dbPostgresql = mainDB
}

// InitMssql 初始化
func (t *SQLClient) InitMssql(mainDB *CoreMsSQL.Client) {
	t.mode = modeMsSQL
	t.dbMssql = mainDB
}

// GetPostgresql 获取postgresql句柄
func (t *SQLClient) GetPostgresql() *CorePostgres.Client {
	return t.dbPostgresql
}
