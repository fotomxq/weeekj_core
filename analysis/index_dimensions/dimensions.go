package AnalysisIndexDimensions

import (
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsGetDimensionsList 获取维度列表参数
type ArgsGetDimensionsList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetDimensionsList 获取列表
func GetDimensionsList(args *ArgsGetDimensionsList) (dataList []FieldsDimensions, dataCount int64, err error) {
	//获取数据
	dataCount, err = dimensionsDB.GetList().GetListSimple(&BaseSQLTools.ArgsGetListSimple{
		Pages:           args.Pages,
		ConditionFields: nil,
		IsRemove:        args.IsRemove,
		Search:          args.Search,
	}, &dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	//反馈
	return
}

// GetDimensionsAll 获取全部维度
func GetDimensionsAll() (dataList []FieldsDimensions, err error) {
	//获取所有维度数据
	_, err = dimensionsDB.GetList().GetAll(&BaseSQLTools.ArgsGetAll{
		ConditionFields: nil,
		IsRemove:        false,
	}, &dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	//反馈
	return
}

// GetDimensionsByID 通过ID获取维度
func GetDimensionsByID(id int64) (data FieldsDimensions, err error) {
	//获取数据
	err = dimensionsDB.GetInfo().GetInfoByID(id, &data)
	if err != nil || data.ID < 1 {
		return
	}
	//反馈
	return
}

// GetDimensionsByCode 通过编码查询指定维度内容
func GetDimensionsByCode(code string) (data FieldsDimensions, err error) {
	//获取数据
	err = dimensionsDB.GetInfo().GetInfoByField("code", code, true, &data)
	if err != nil || data.ID < 1 {
		return
	}
	//反馈
	return
}

// ArgsCreateDimensions 创建维度参数
type ArgsCreateDimensions struct {
	//编码
	// 维度编码，用于程序内部识别
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true" field_search:"true" field_list:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"0" max:"500" empty:"true" field_search:"true"`
	//约定Extend字段名称
	// 约定指标、指标值
	// 例如：extend1
	ExtendIndex string `db:"extend_index" json:"extendIndex" index:"true"`
	//所属表名称
	TableName string `db:"table_name" json:"tableName" check:"des" min:"1" max:"50"`
	//字段名称
	FieldName string `db:"field_name" json:"fieldName" check:"des" min:"1" max:"50"`
}

// CreateDimensions 创建维度
func CreateDimensions(args *ArgsCreateDimensions) (id int64, err error) {
	//创建数据
	id, err = dimensionsDB.GetInsert().InsertRow(args)
	if err != nil || id < 1 {
		return
	}
	//反馈
	return
}

// ArgsUpdateDimensions 更新维度
type ArgsUpdateDimensions struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true" field_search:"true" field_list:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"0" max:"500" empty:"true" field_search:"true"`
	//约定Extend字段名称
	// 约定指标、指标值
	// 例如：extend1
	ExtendIndex string `db:"extend_index" json:"extendIndex" index:"true"`
	//所属表名称
	TableName string `db:"table_name" json:"tableName" check:"des" min:"1" max:"50"`
	//字段名称
	FieldName string `db:"field_name" json:"fieldName" check:"des" min:"1" max:"50"`
}

func UpdateDimensions(args *ArgsUpdateDimensions) (err error) {
	//更新数据
	err = dimensionsDB.GetUpdate().UpdateByID(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// RemoveDimensions 删除维度
func RemoveDimensions(id int64) (err error) {
	//删除数据
	err = dimensionsDB.GetDelete().DeleteByID(id)
	if err != nil {
		return
	}
	//反馈
	return
}
