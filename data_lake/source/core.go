package DataLakeSource

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

/**
数据湖模块
1. 该模块可快速导入Excel/csv，并建立专门的数据表结构
2. 该模块也可以在前端构建手动建立表结构，方便后续导入数据或集成外部数据
*/

var (
	//DB
	tableDB  BaseSQLTools.Quick
	fieldsDB BaseSQLTools.Quick
)

func Init() (err error) {
	//初始化数据库
	if err = tableDB.Init("data_lake_source_table", &FieldsTable{}); err != nil {
		return
	}
	if err = fieldsDB.Init("data_lake_source_fields", &FieldsFields{}); err != nil {
		return
	}
	//反馈
	return
}
