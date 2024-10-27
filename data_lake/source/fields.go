package DataLakeSource

import (
	"errors"
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsGetFieldsList 查看列表参数
type ArgsGetFieldsList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//表ID
	TableID int64 `db:"table_id" json:"tableId" index:"true"`
	//字段名
	FieldName string `db:"field_name" json:"fieldName" field_search:"true"`
	//字段表单类型
	InputType string `db:"input_type" json:"inputType" field_search:"true"`
	//字段数据类型
	DataType string `db:"data_type" json:"dataType" field_search:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetFieldsList 查看列表
func GetFieldsList(args *ArgsGetFieldsList) (dataList []FieldsFields, dataCount int64, err error) {
	dataCount, err = fieldsDB.GetList().GetListSimple(&BaseSQLTools.ArgsGetListSimple{
		Pages: args.Pages,
		ConditionFields: []BaseSQLTools.ArgsGetListSimpleConditionID{
			{
				Name: "field_name",
				Val:  args.FieldName,
			},
			{
				Name: "input_type",
				Val:  args.InputType,
			},
			{
				Name: "data_type",
				Val:  args.DataType,
			},
		},
		IsRemove: args.IsRemove,
		Search:   args.Search,
	}, &dataList)
	if err != nil {
		return
	}
	if len(dataList) == 0 {
		err = errors.New("not found")
		return
	}
	return
}

// GetFieldsDetail 查看表详情
func GetFieldsDetail(id int64) (data FieldsFields, err error) {
	err = fieldsDB.GetInfo().GetInfoByID(id, &data)
	if err != nil {
		return
	}
	return
}

// GetFieldsListByTableID 获取表的所有列
func GetFieldsListByTableID(tableID int64) (dataList []FieldsFields, dataCount int64, err error) {
	dataCount, err = fieldsDB.GetList().GetAll(&BaseSQLTools.ArgsGetAll{
		ConditionFields: []BaseSQLTools.ArgsGetListSimpleConditionID{
			{
				Name: "table_id",
				Val:  tableID,
			},
		},
		IsRemove: false,
	}, &dataList)
	if err != nil {
		return
	}
	if len(dataList) == 0 {
		err = errors.New("not found")
		return
	}
	return
}

// GetFieldsDetailByTableIDAndFieldName 根据表ID和字段名称的信息
func GetFieldsDetailByTableIDAndFieldName(tableID int64, fieldName string) (data FieldsFields, err error) {
	var dataList []FieldsFields
	_, err = fieldsDB.GetList().GetAll(&BaseSQLTools.ArgsGetAll{
		ConditionFields: []BaseSQLTools.ArgsGetListSimpleConditionID{
			{
				Name: "table_id",
				Val:  tableID,
			},
			{
				Name: "field_name",
				Val:  fieldName,
			},
		},
		IsRemove: false,
	}, &dataList)
	if err != nil {
		return
	}
	if len(dataList) == 0 {
		err = errors.New("not found")
		return
	}
	data = dataList[0]
	return
}

// ArgsCreateFields 创建表参数
type ArgsCreateFields struct {
	//表ID
	TableID int64 `db:"table_id" json:"tableId" index:"true"`
	//字段名
	FieldName string `db:"field_name" json:"fieldName" field_search:"true"`
	//提示名称
	FieldLabel string `db:"field_label" json:"fieldLabel" field_search:"true"`
	//字段表单类型
	// string/number/textarea/checkbox/radio/select/date/datetime
	InputType string `db:"input_type" json:"inputType" field_search:"true"`
	//字段表单长度
	InputLength int `db:"input_length" json:"inputLength"`
	//字段表单默认值
	InputDefault string `db:"input_default" json:"inputDefault"`
	//字段表单是否必填
	InputRequired bool `db:"input_required" json:"inputRequired"`
	//字段表单正则表达式
	InputPattern string `db:"input_pattern" json:"inputPattern"`
	//是否为主键
	IsPrimary bool `db:"is_primary" json:"isPrimary"`
	//字段是否为索引
	IsIndex bool `db:"is_index" json:"isIndex"`
	//是否为系统内置字段
	// id/create_at/update_at/delete_at
	IsSystem bool `db:"is_system" json:"isSystem"`
	//是否支持搜索
	IsSearch bool `db:"is_search" json:"isSearch"`
	//字段数据类型
	// int/int64/float/text/bool/date/datetime
	DataType string `db:"data_type" json:"dataType" field_search:"true"`
	//字段描述
	FieldDesc string `db:"field_desc" json:"fieldDesc" field_search:"true"`
}

// CreateFields 创建表
func CreateFields(args *ArgsCreateFields) (newID int64, err error) {
	//检查数据表是否存在
	var tableData FieldsTable
	tableData, err = GetTableDetail(args.TableID)
	if err != nil || tableData.ID < 1 || CoreFilter.CheckHaveTime(tableData.DeleteAt) {
		err = errors.New("table not found")
		return
	}
	//检查字段名称
	var fieldData FieldsFields
	fieldData, err = GetFieldsDetailByTableIDAndFieldName(args.TableID, args.FieldName)
	if err == nil && fieldData.ID > 0 && !CoreFilter.CheckHaveTime(fieldData.DeleteAt) {
		err = errors.New("field is exist")
		return
	}
	//写入数据
	newID, err = fieldsDB.GetInsert().InsertRow(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateFields 修改表参数
type ArgsUpdateFields struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//字段名
	FieldName string `db:"field_name" json:"fieldName" field_search:"true"`
	//提示名称
	FieldLabel string `db:"field_label" json:"fieldLabel" field_search:"true"`
	//字段表单类型
	// input/textarea/select/radio/checkbox/date/datetime
	InputType string `db:"input_type" json:"inputType" field_search:"true"`
	//字段表单长度
	InputLength int `db:"input_length" json:"inputLength"`
	//字段表单默认值
	InputDefault string `db:"input_default" json:"inputDefault"`
	//字段表单是否必填
	InputRequired bool `db:"input_required" json:"inputRequired"`
	//字段表单正则表达式
	InputPattern string `db:"input_pattern" json:"inputPattern"`
	//是否为主键
	IsPrimary bool `db:"is_primary" json:"isPrimary"`
	//字段是否为索引
	IsIndex bool `db:"is_index" json:"isIndex"`
	//是否支持搜索
	IsSearch bool `db:"is_search" json:"isSearch"`
	//字段描述
	FieldDesc string `db:"field_desc" json:"fieldDesc" field_search:"true"`
}

// UpdateFields 修改表
func UpdateFields(args *ArgsUpdateFields) (err error) {
	err = fieldsDB.GetUpdate().UpdateByID(args)
	if err != nil {
		return
	}
	return
}

// DeleteFields 删除表
// 必须先删除表数据和结构后才能删除，否则报错
func DeleteFields(id int64) (err error) {
	err = fieldsDB.GetDelete().DeleteByID(id)
	if err != nil {
		return
	}
	return
}

// ClearFields 清理表字段
func ClearFields(tableID int64) (err error) {
	err = fieldsDB.GetDelete().DeleteByField("table_id", tableID)
	if err != nil {
		return
	}
	return
}
