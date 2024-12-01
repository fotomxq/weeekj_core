package AnalysisIndexFilter

import (
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsGetList 获取数据列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取数据列表
func GetList(args *ArgsGetList) (dataList []FieldsFilter, dataCount int64, err error) {
	//构建筛选条件
	dataCount, err = filterDB.GetList().GetListSimple(&BaseSQLTools.ArgsGetListSimple{
		Pages: args.Pages,
		ConditionFields: []BaseSQLTools.ArgsGetListSimpleConditionID{
			{
				Name: "code",
				Val:  args.Code,
			},
		},
		IsRemove: args.IsRemove,
		Search:   args.Search,
	}, &dataList)
	if err != nil {
		return
	}
	for k, v := range dataList {
		var vData FieldsFilter
		_ = filterDB.GetInfo().GetInfoByID(v.ID, vData)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	//反馈
	return
}

// GetAll 获取指标所有数据
func GetAll(code string) (dataList []FieldsFilter) {
	_, _ = filterDB.GetList().GetAll(&BaseSQLTools.ArgsGetAll{
		ConditionFields: []BaseSQLTools.ArgsGetListSimpleConditionID{
			{
				Name: "code",
				Val:  code,
			},
		},
	}, &dataList)
	return
}
