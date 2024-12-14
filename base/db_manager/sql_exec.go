package BaseDBManager

import (
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsGetSQLList 获取SQL列表参数
type ArgsGetSQLList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//来源系统
	// 例如: analysis
	FromSystem string `db:"from_system" json:"fromSystem" index:"true" field_list:"true"`
	//来源模块
	// 例如: index_sql
	FromModule string `db:"from_module" json:"fromModule" index:"true" field_list:"true"`
	//内部标识码
	// 可用于标记内部识别标识码，例如Index中的维度值，或一组维度值组合后的标识码
	FromCode string `db:"from_code" json:"fromCode" index:"true" field_list:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetSQLList 获取SQL列表
func GetSQLList(args *ArgsGetSQLList) (dataList []FieldsSQL, dataCount int64, err error) {
	//构建参数
	var conditionFields []BaseSQLTools.ArgsGetListSimpleConditionID
	conditionFields = append(conditionFields, BaseSQLTools.ArgsGetListSimpleConditionID{
		Name: "from_system",
		Val:  args.FromSystem,
	})
	conditionFields = append(conditionFields, BaseSQLTools.ArgsGetListSimpleConditionID{
		Name: "from_module",
		Val:  args.FromModule,
	})
	conditionFields = append(conditionFields, BaseSQLTools.ArgsGetListSimpleConditionID{
		Name: "from_code",
		Val:  args.FromCode,
	})
	//获取数据
	dataCount, err = sqlDB.GetList().GetListSimple(&BaseSQLTools.ArgsGetListSimple{
		Pages:           args.Pages,
		ConditionFields: conditionFields,
		IsRemove:        args.IsRemove,
		Search:          args.Search,
	}, &dataList)
	if err != nil {
		return
	}
	//反馈
	return
}

// GetSQLByID 通过ID获取指定SQL
func GetSQLByID(id int64) (data FieldsSQL, err error) {
	//获取数据
	err = sqlDB.GetInfo().GetInfoByID(id, &data)
	if err != nil {
		return
	}
	//反馈
	return
}

// GetSQLByCode 通过Code获取指定SQL
func GetSQLByCode(fromSystem, fromModule, fromCode string) (data FieldsSQL, err error) {
	//获取数据
	err = sqlDB.GetInfo().GetInfoByFields(map[string]any{
		"from_system": fromSystem,
		"from_module": fromModule,
		"from_code":   fromCode,
	}, true, &data)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsCreateSQL 创建SQL参数
type ArgsCreateSQL struct {
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 来源
	// 如果存在值，尤其是带有FromCode时，应确保数据唯一性
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//来源系统
	// 例如: analysis
	FromSystem string `db:"from_system" json:"fromSystem" index:"true" field_list:"true"`
	//来源模块
	// 例如: index_sql
	FromModule string `db:"from_module" json:"fromModule" index:"true" field_list:"true"`
	//内部标识码
	// 可用于标记内部识别标识码，例如Index中的维度值，或一组维度值组合后的标识码
	FromCode string `db:"from_code" json:"fromCode" index:"true" field_list:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 基础设置
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//定时器设置Carbon编码
	// 例如: 15s
	CarbonCode string `db:"carbon_code" json:"carbonCode" index:"true" field_list:"true"`
	//开始运行时通知中间件地址
	// 用于通知需发起该SQL，将SQL和来源信息传递给对应的中间件
	PostURL string `db:"post_url" json:"postURL" index:"true" field_list:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 数据
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//SQL内容
	SQLData string `db:"sql_data" json:"sqlData"`
}

// CreateSQL 创建SQL
func CreateSQL(args *ArgsCreateSQL) (err error) {
	//执行创建
	_, err = sqlDB.GetInsert().InsertRow(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateSQL 修改SQL参数
type ArgsUpdateSQL struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 基础设置
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//定时器设置Carbon编码
	// 例如: 15s
	CarbonCode string `db:"carbon_code" json:"carbonCode" index:"true" field_list:"true"`
	//开始运行时通知中间件地址
	// 用于通知需发起该SQL，将SQL和来源信息传递给对应的中间件
	PostURL string `db:"post_url" json:"postURL" index:"true" field_list:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 数据
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//SQL内容
	SQLData string `db:"sql_data" json:"sqlData"`
}

// UpdateSQL 修改SQL
func UpdateSQL(args *ArgsUpdateSQL) (err error) {
	//执行更新
	err = sqlDB.GetUpdate().UpdateByID(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// DeleteSQLByID 删除SQL
func DeleteSQLByID(id int64) (err error) {
	//执行删除
	err = sqlDB.GetDelete().DeleteByID(id)
	if err != nil {
		return
	}
	//反馈
	return
}
