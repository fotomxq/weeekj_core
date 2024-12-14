package BaseDBManager

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

//数据库表管理
/**
用途：
1. 提供数据库表管理功能
2. 提供SQL存储和执行功能，该功能可触发中间件，需通过中间件接收和处理
	该内容也可以用于直接执行SQL，通过SQL自带能力直接写入对应的表内，避免二次代码开发工作
*/

var (
	//SQL Exec
	sqlDB BaseSQLTools.Quick
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() (err error) {
	//初始化指标定义
	if err = sqlDB.Init("base_db_manager_sql", &FieldsSQL{}); err != nil {
		return
	}
	if OpenSub {
		//消息列队
		//subNats()
	}
	//反馈
	return
}
