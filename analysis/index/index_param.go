package AnalysisIndex

import (
	"errors"
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsGetIndexParamList 获取指标列表参数
type ArgsGetIndexParamList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//指标ID
	IndexID int64 `json:"indexID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetIndexParamList 获取指标列表
func GetIndexParamList(args *ArgsGetIndexParamList) (dataList []FieldsIndexParam, dataCount int64, err error) {
	//组合条件
	var conditionFields = []BaseSQLTools.ArgsGetListSimpleConditionID{
		{
			Name: "index_id",
			Val:  args.IndexID,
		},
	}
	//获取数据
	dataCount, err = indexParamDB.GetList().GetListSimple(&BaseSQLTools.ArgsGetListSimple{
		Pages:           args.Pages,
		ConditionFields: conditionFields,
		IsRemove:        args.IsRemove,
		Search:          args.Search,
	}, &dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		var vData FieldsIndexParam
		err = indexParamDB.GetInfo().GetInfoByID(v.ID, &vData)
		if err != nil || vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	//反馈
	return
}

// GetIndexParamByIndexCode 获取指定编码参数
func GetIndexParamByIndexCode(indexCode string, code string) (data FieldsIndexParam, err error) {
	//获取指标
	var indexData FieldsIndex
	indexData, err = GetIndexByCode(indexCode)
	if err != nil {
		return
	}
	//获取数据
	err = indexParamDB.GetInfo().GetInfoByFields(map[string]any{
		"index_id": indexData.ID,
		"code":     code,
	}, true, &data)
	if err != nil {
		return
	}
	//反馈
	return
}

// GetIndexParamByCode 获取指定编码参数
func GetIndexParamByCode(indexID int64, code string) (data FieldsIndexParam, err error) {
	//获取数据
	err = indexParamDB.GetInfo().GetInfoByFields(map[string]any{
		"index_id": indexID,
		"code":     code,
	}, true, &data)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsCreateIndexParam 创建新的指标参数
type ArgsCreateIndexParam struct {
	//指标ID
	IndexID int64 `db:"index_id" json:"indexID" check:"id" index:"true"`
	//参数名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"50" field_search:"true" field_list:"true"`
	//参数编码
	// 用于程序内识别内置指标的参数
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" field_search:"true" field_list:"true"`
	//参数值
	ParamVal string `db:"param_val" json:"paramVal"`
}

// CreateIndexParam 创建新的指标
func CreateIndexParam(args *ArgsCreateIndexParam) (err error) {
	//指标必须存在
	indexData, _ := GetIndexByID(args.IndexID)
	if indexData.ID < 1 {
		err = errors.New("index not exist")
		return
	}
	//检查code是否重复
	if b, _ := indexParamDB.GetInfo().CheckInfoByFields(map[string]any{
		"index_id": args.IndexID,
		"code":     args.Code,
	}, true); b {
		err = errors.New("code is exist")
		return
	}
	//创建
	_, err = indexParamDB.GetInsert().InsertRow(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateIndexParam 更新指标定义参数
type ArgsUpdateIndexParam struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//参数值
	ParamVal string `db:"param_val" json:"paramVal"`
}

// UpdateIndexParam 更新指标定义
func UpdateIndexParam(args *ArgsUpdateIndexParam) (err error) {
	//修改数据
	err = indexParamDB.GetUpdate().UpdateByID(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// DeleteIndexParam 删除指标参数
func DeleteIndexParam(id int64) (err error) {
	//删除数据
	err = indexParamDB.GetDelete().DeleteByID(id)
	if err != nil {
		return
	}
	//反馈
	return
}
