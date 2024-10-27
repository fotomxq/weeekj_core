package DataLakeSource

import (
	"errors"
	"fmt"
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsGetTableList 查看列表参数
type ArgsGetTableList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTableList 查看列表
func GetTableList(args *ArgsGetTableList) (dataList []FieldsTable, dataCount int64, err error) {
	dataCount, err = tableDB.GetList().GetListSimple(&BaseSQLTools.ArgsGetListSimple{
		Pages:           args.Pages,
		ConditionFields: []BaseSQLTools.ArgsGetListSimpleConditionID{},
		IsRemove:        args.IsRemove,
		Search:          args.Search,
	}, &dataList)
	if err != nil {
		return
	}
	return
}

// GetTableDetail 查看表详情
func GetTableDetail(id int64) (data FieldsTable, err error) {
	err = tableDB.GetInfo().GetInfoByID(id, &data)
	if err != nil {
		return
	}
	return
}

// GetTableDetailByName 找到表名称
func GetTableDetailByName(tableName string) (data FieldsTable, err error) {
	err = tableDB.GetInfo().GetInfoByField("table_name", tableName, true, &data)
	if err != nil {
		return
	}
	return
}

// ArgsCreateTable 创建表参数
type ArgsCreateTable struct {
	//表名称
	TableName string `db:"table_name" json:"tableName" field_search:"true"`
	//表描述
	TableDesc string `db:"table_desc" json:"tableDesc" field_search:"true"`
	//提示名称
	TipName string `db:"tip_name" json:"tipName" field_search:"true"`
	//数据唯一渠道名称
	// 如果是多处来源，应拆分表
	ChannelName string `db:"channel_name" json:"channelName" field_search:"true"`
	//数据唯一渠道提示名称
	ChannelTipName string `db:"channel_tip_name" json:"channelTipName" field_search:"true"`
}

// CreateTable 创建表
func CreateTable(args *ArgsCreateTable) (newID int64, err error) {
	//检查是否存在重复表
	var tableData FieldsTable
	tableData, err = GetTableDetailByName(args.TableName)
	if err == nil && tableData.ID > 0 && !CoreFilter.CheckHaveTime(tableData.DeleteAt) {
		err = errors.New(fmt.Sprint("table is exist, table name: ", args.TableName, ", table id: ", tableData.ID))
		return
	}
	//写入数据
	newID, err = tableDB.GetInsert().InsertRow(args)
	if err != nil {
		err = errors.New(fmt.Sprint("create table error: ", err))
		return
	}
	//初始化建立系统内置字段
	// id/create_at/update_at/delete_at
	_, err = CreateFields(&ArgsCreateFields{
		TableID:       newID,
		FieldName:     "id",
		FieldLabel:    "ID",
		InputType:     "number",
		InputLength:   0,
		InputDefault:  "id",
		InputRequired: true,
		InputPattern:  "",
		IsPrimary:     true,
		IsIndex:       true,
		IsSystem:      true,
		IsSearch:      false,
		DataType:      "int64",
		FieldDesc:     "ID",
	})
	if err != nil {
		err = errors.New(fmt.Sprint("create id field error: ", err))
	}
	_, err = CreateFields(&ArgsCreateFields{
		TableID:       newID,
		FieldName:     "create_at",
		FieldLabel:    "创建时间",
		InputType:     "datetime",
		InputLength:   0,
		InputDefault:  "now",
		InputRequired: true,
		InputPattern:  "",
		IsPrimary:     false,
		IsIndex:       false,
		IsSystem:      true,
		IsSearch:      false,
		DataType:      "datetime",
		FieldDesc:     "创建时间",
	})
	if err != nil {
		err = errors.New(fmt.Sprint("create create_at field error: ", err))
		return
	}
	_, err = CreateFields(&ArgsCreateFields{
		TableID:       newID,
		FieldName:     "update_at",
		FieldLabel:    "更新时间",
		InputType:     "datetime",
		InputLength:   0,
		InputDefault:  "now",
		InputRequired: true,
		InputPattern:  "",
		IsPrimary:     false,
		IsIndex:       false,
		IsSystem:      true,
		IsSearch:      false,
		DataType:      "datetime",
		FieldDesc:     "更新时间",
	})
	if err != nil {
		err = errors.New(fmt.Sprint("create update_at field error: ", err))
		return
	}
	_, err = CreateFields(&ArgsCreateFields{
		TableID:       newID,
		FieldName:     "delete_at",
		FieldLabel:    "删除时间",
		InputType:     "datetime",
		InputLength:   0,
		InputDefault:  "",
		InputRequired: true,
		InputPattern:  "",
		IsPrimary:     false,
		IsIndex:       false,
		IsSystem:      true,
		IsSearch:      false,
		DataType:      "datetime",
		FieldDesc:     "删除时间",
	})
	if err != nil {
		err = errors.New(fmt.Sprint("create delete_at field error: ", err))
		return
	}
	//反馈
	return
}

// ArgsUpdateTable 修改表参数
type ArgsUpdateTable struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//表名称
	TableName string `db:"table_name" json:"tableName" field_search:"true"`
	//表描述
	TableDesc string `db:"table_desc" json:"tableDesc" field_search:"true"`
	//提示名称
	TipName string `db:"tip_name" json:"tipName" field_search:"true"`
	//数据唯一渠道名称
	// 如果是多处来源，应拆分表
	ChannelName string `db:"channel_name" json:"channelName" field_search:"true"`
	//数据唯一渠道提示名称
	ChannelTipName string `db:"channel_tip_name" json:"channelTipName" field_search:"true"`
}

// UpdateTable 修改表
func UpdateTable(args *ArgsUpdateTable) (err error) {
	err = tableDB.GetUpdate().UpdateByID(args)
	if err != nil {
		return
	}
	return
}

// DeleteTable 删除表
// 必须先删除表数据和结构后才能删除，否则报错
func DeleteTable(id int64) (err error) {
	//检查是否存在列记录
	var dataCount int64
	_, dataCount, err = GetFieldsListByTableID(id)
	if err == nil && dataCount > 0 {
		err = errors.New("has fields")
		return
	}
	//执行删除
	err = tableDB.GetDelete().DeleteByID(id)
	if err != nil {
		err = errors.New(fmt.Sprint("delete table error: ", err, ", id: ", id))
		return
	}
	//反馈
	return
}
