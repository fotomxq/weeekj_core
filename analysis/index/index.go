package AnalysisIndex

import (
	"errors"
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsGetIndexList 获取指标列表参数
type ArgsGetIndexList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetIndexList 获取指标列表
func GetIndexList(args *ArgsGetIndexList) (dataList []FieldsIndex, dataCount int64, err error) {
	//获取数据
	dataCount, err = indexDB.GetList().GetListSimple(&BaseSQLTools.ArgsGetListSimple{
		Pages:           args.Pages,
		ConditionFields: nil,
		IsRemove:        args.IsRemove,
		Search:          args.Search,
	}, &dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		var vData FieldsIndex
		err = indexDB.GetInfo().GetInfoByID(v.ID, &vData)
		if err != nil || vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	//反馈
	return
}

// DataGetIndexListByTop 获取指标列表顶部数据
type DataGetIndexListByTop struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//指标名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" index:"true" field_search:"true" field_list:"true"`
	//指标描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"100" field_search:"true" field_list:"true" empty:"true"`
	//指标决策建议
	Decision string `db:"decision" json:"decision" check:"des" min:"1" max:"-1" empty:"true" field_search:"true"`
}

// GetIndexListByTop 获取指标列表顶部
func GetIndexListByTop() (dataList []DataGetIndexListByTop, dataCount int64, err error) {
	//获取数据
	err = indexDB.GetClient().DB.GetPostgresql().Select(&dataList, "SELECT i.code as code, max(i.name) as name, max(i.description) as description, max(i.decision) as decision FROM analysis_index as i, analysis_index_relation as r WHERE i.delete_at < to_timestamp(100000) and i.id = r.index_id GROUP BY i.code ORDER BY i.code;")
	if err != nil || len(dataList) < 1 {
		return
	}
	//反馈
	return
}

// GetIndexByCode 通过编码获取指标
func GetIndexByCode(code string) (data FieldsIndex, err error) {
	//获取数据
	err = indexDB.GetInfo().GetInfoByFields(map[string]any{
		"code": code,
	}, true, &data)
	if err != nil {
		return
	}
	//反馈
	return
}

// GetIndexByID 通过ID获取指标
func GetIndexByID(id int64) (data FieldsIndex, err error) {
	//获取数据
	err = indexDB.GetInfo().GetInfoByID(id, &data)
	if err != nil {
		return
	}
	//反馈
	return
}

// GetIndexNameByID 获取指标名称
func GetIndexNameByID(id int64) (name string) {
	data, _ := GetIndexByID(id)
	return data.Name
}

// ArgsCreateIndex 创建新的指标参数
type ArgsCreateIndex struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//指标名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" index:"true" field_search:"true" field_list:"true"`
	//是否内置
	// 前端应拦截内置指标的删除操作，以免影响系统正常运行，启动重启后将自动恢复，所以删除操作是无法生效的
	// 通过接口应强制给予false
	IsSystem bool `db:"is_system" json:"isSystem" index:"true"`
	//指标描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"100" field_search:"true" field_list:"true" empty:"true"`
	//指标决策建议
	Decision string `db:"decision" json:"decision" check:"des" min:"1" max:"-1" empty:"true" field_search:"true"`
	//指标预警阈值
	// 0-100，归一化后的数据，超出此范围将可触发预警事件记录
	Threshold int64 `db:"threshold" json:"threshold" index:"true"`
	//是否启用
	// 关闭后将不对该指标进行汇总运算
	// 通过接口应强制给予false
	IsEnable bool `db:"is_enable" json:"isEnable" index:"true"`
}

// CreateIndex 创建新的指标
func CreateIndex(args *ArgsCreateIndex) (err error) {
	//检查code是否重复
	if b, _ := indexDB.GetInfo().CheckInfoByFields(map[string]any{
		"code": args.Code,
	}, true); b {
		err = errors.New("code is exist")
		return
	}
	//修改启动设置
	args.IsEnable = false
	//创建
	_, err = indexDB.GetInsert().InsertRow(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateIndex 更新指标定义参数
type ArgsUpdateIndex struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//指标名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" index:"true" field_search:"true" field_list:"true"`
	//指标描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"100" field_search:"true" field_list:"true" empty:"true"`
	//指标决策建议
	Decision string `db:"decision" json:"decision" check:"des" min:"1" max:"-1" empty:"true" field_search:"true"`
	//指标预警阈值
	// 0-100，归一化后的数据，超出此范围将可触发预警事件记录
	Threshold int64 `db:"threshold" json:"threshold" index:"true"`
	//是否启用
	// 关闭后将不对该指标进行汇总运算
	IsEnable bool `db:"is_enable" json:"isEnable" index:"true"`
}

// UpdateIndex 更新指标定义
func UpdateIndex(args *ArgsUpdateIndex) (err error) {
	//修改数据
	err = indexDB.GetUpdate().UpdateByID(args)
	if err != nil {
		return
	}
	//反馈
	return
}

// DeleteIndex 删除指标
func DeleteIndex(id int64) (err error) {
	//删除数据
	err = indexDB.GetDelete().DeleteByID(id)
	if err != nil {
		return
	}
	//删除指标的参数
	_ = indexParamDB.GetDelete().DeleteByField("index_id", id)
	//删除指标的关系
	_ = indexRelationDB.GetDelete().DeleteByField("index_id", id)
	_ = indexRelationDB.GetDelete().DeleteByField("relation_index_id", id)
	//反馈
	return
}
