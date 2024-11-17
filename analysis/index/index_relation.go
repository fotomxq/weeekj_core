package AnalysisIndex

import (
	"errors"
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsGetIndexRelationList 获取指标关系列表参数
type ArgsGetIndexRelationList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//指标ID
	IndexID int64 `json:"indexID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetIndexRelationList 获取指标列表
func GetIndexRelationList(args *ArgsGetIndexRelationList) (dataList []FieldsIndexRelation, dataCount int64, err error) {
	//组合条件
	var conditionFields = []BaseSQLTools.ArgsGetListSimpleConditionID{
		{
			Name: "index_id",
			Val:  args.IndexID,
		},
	}
	//获取数据
	dataCount, err = indexRelationDB.GetList().GetListSimple(&BaseSQLTools.ArgsGetListSimple{
		Pages:           args.Pages,
		ConditionFields: conditionFields,
		IsRemove:        args.IsRemove,
		Search:          args.Search,
	}, &dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		var vData FieldsIndexRelation
		err = indexRelationDB.GetInfo().GetInfoByID(v.ID, &vData)
		if err != nil || vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	//反馈
	return
}

// ArgsCreateIndexRelation 创建指标关系参数
type ArgsCreateIndexRelation struct {
	//指标ID
	// 上级指标
	IndexID int64 `db:"index_id" json:"indexID" check:"id"`
	//关联指标
	// 禁出现嵌套关系，系统将检查并报错
	RelationIndexID int64 `db:"relation_index_id" json:"relationIndexID" check:"id"`
	//指标权重占比
	// 同一个indexID下，所有指标的权重总和必须为1，否则计算的结果将出现异常
	Weight int64 `db:"weight" json:"weight" check:"int64Than0"`
	//算法自动权重
	// 同一个indexID下，所有指标的权重总和必须为1，否则计算的结果将出现异常
	// 接口不能设置该参数，必须程序化内置实现
	AutoWeight int64 `db:"auto_weight" json:"autoWeight" check:"int64Than0"`
	//是否启动算法自动权重
	// 接口不能设置该参数，必须程序化内置实现
	IsAutoWeight bool `db:"is_auto_weight" json:"isAutoWeight"`
}

// CreateIndexRelation 创建指标关系
func CreateIndexRelation(args *ArgsCreateIndexRelation) (err error) {
	//指标必须存在
	indexData, _ := GetIndexByID(args.IndexID)
	if indexData.ID < 1 {
		err = errors.New("index not exist")
		return
	}
	relationIndexData, _ := GetIndexByID(args.RelationIndexID)
	if relationIndexData.ID < 1 {
		err = errors.New("relation index not exist")
		return
	}
	//检查是否存在该关系
	if b, _ := indexRelationDB.GetInfo().CheckInfoByFields(map[string]any{
		"index_id":          args.IndexID,
		"relation_index_id": args.RelationIndexID,
	}, true); b {
		err = errors.New("relation is exist")
		return
	}
	//创建
	_, err = indexRelationDB.GetInsert().InsertRow(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateIndexRelation 更新指标定义参数
type ArgsUpdateIndexRelation struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//指标权重占比
	// 同一个indexID下，所有指标的权重总和必须为1，否则计算的结果将出现异常
	Weight int64 `db:"weight" json:"weight" check:"int64Than0"`
	//算法自动权重
	// 同一个indexID下，所有指标的权重总和必须为1，否则计算的结果将出现异常
	// 接口不能设置该参数，必须程序化内置实现
	AutoWeight int64 `db:"auto_weight" json:"autoWeight" check:"int64Than0"`
	//是否启动算法自动权重
	// 接口不能设置该参数，必须程序化内置实现
	IsAutoWeight bool `db:"is_auto_weight" json:"isAutoWeight"`
}

// UpdateIndexRelation 更新指标定义
func UpdateIndexRelation(args *ArgsUpdateIndexRelation) (err error) {
	//修改数据
	err = indexRelationDB.GetUpdate().UpdateByID(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// DeleteIndexRelation 删除指标
func DeleteIndexRelation(id int64) (err error) {
	//删除数据
	err = indexRelationDB.GetDelete().DeleteByID(id)
	if err != nil {
		return
	}
	//反馈
	return
}
