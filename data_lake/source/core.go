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
	//表单字段类型枚举值
	// input/textarea/select/radio/checkbox/date/datetime
	FIELDS_INPUT_TYPE_ENUM_INPUT    = "input"
	FIELDS_INPUT_TYPE_ENUM_TEXTAREA = "textarea"
	FIELDS_INPUT_TYPE_ENUM_SELECT   = "select"
	FIELDS_INPUT_TYPE_ENUM_RADIO    = "radio"
	FIELDS_INPUT_TYPE_ENUM_CHECKBOX = "checkbox"
	FIELDS_INPUT_TYPE_ENUM_DATE     = "date"
	FIELDS_INPUT_TYPE_ENUM_DATETIME = "datetime"
	//字段数据类型
	// int/int64/float/text/bool/date/datetime
	FIELDS_DATA_TYPE_ENUM_INT      = "int"
	FIELDS_DATA_TYPE_ENUM_INT64    = "int64"
	FIELDS_DATA_TYPE_ENUM_FLOAT    = "float"
	FIELDS_DATA_TYPE_ENUM_TEXT     = "text"
	FIELDS_DATA_TYPE_ENUM_BOOL     = "bool"
	FIELDS_DATA_TYPE_ENUM_DATE     = "date"
	FIELDS_DATA_TYPE_ENUM_DATETIME = "datetime"
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
