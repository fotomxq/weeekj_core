package AnalysisIndexEvent

import (
	"errors"
	"fmt"
	BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsGetEventList 获取预警事件列表参数
type ArgsGetEventList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//触发时间范围
	// 例如：2021-01-01
	BetweenAt CoreSQL2.ArgsTimeBetween `json:"betweenAt"`
	//预警等级
	// 根据项目需求划定等级
	// 例如：0 低风险; 1 中风险; 2 高风险
	Level int `db:"level" json:"level" index:"true" field_list:"true"`
	//来源指标值的系统和ID
	// 避免重复触发预警
	FromSystem string `db:"from_system" json:"fromSystem" check:"des" min:"1" max:"50" empty:"true" index:"true" field_list:"true"`
	FromID     int64  `db:"from_id" json:"fromID" check:"id" empty:"true" index:"true" field_list:"true"`
	//触发类型
	// 根据项目需求划定类型，可以留空
	FromType string `db:"from_type" json:"fromType" check:"des" min:"1" max:"50" empty:"true" index:"true" field_list:"true"`
	//扩展维度1
	// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
	Extend1 string `db:"extend1" json:"extend1" index:"true" field_list:"true"`
	//扩展维度2
	Extend2 string `db:"extend2" json:"extend2" index:"true" field_list:"true"`
	//扩展维度3
	Extend3 string `db:"extend3" json:"extend3" index:"true" field_list:"true"`
	//扩展维度4
	Extend4 string `db:"extend4" json:"extend4" index:"true" field_list:"true"`
	//扩展维度5
	Extend5 string `db:"extend5" json:"extend5" index:"true" field_list:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetEventList 获取预警事件列表
func GetEventList(args *ArgsGetEventList) (dataList []FieldsEvent, dataCount int64, err error) {
	//构建筛选条件
	ctx := eventDB.GetList().GetListPrev(&BaseSQLTools.ArgsGetListPrev{
		Pages:    args.Pages,
		IsRemove: args.IsRemove,
		Search:   args.Search,
	})
	ctx = ctx.SetStringQuery("code", args.Code)
	if args.BetweenAt.MinTime != "" || args.BetweenAt.MaxTime != "" {
		ctx = ctx.SetTimeBetweenByArgQuery("year_md", args.BetweenAt)
	}
	ctx = ctx.SetIntQuery("level", args.Level)
	ctx = ctx.SetStringQuery("from_system", args.FromSystem)
	ctx = ctx.SetIDQuery("from_id", args.FromID)
	ctx = ctx.SetStringQuery("from_type", args.FromType)
	ctx = ctx.SetStringQuery("extend1", args.Extend1)
	ctx = ctx.SetStringQuery("extend2", args.Extend2)
	ctx = ctx.SetStringQuery("extend3", args.Extend3)
	ctx = ctx.SetStringQuery("extend4", args.Extend4)
	ctx = ctx.SetStringQuery("extend5", args.Extend5)
	//获取数据
	dataCount, err = eventDB.GetList().GetListDo(ctx, &dataList)
	if err != nil || len(dataList) < 1 {
		err = errors.New(fmt.Sprint("event get list error:", err))
		return
	}
	//反馈
	return
}

// 通过ID获取风险详情
func GetEventByID(id int64) (data FieldsEvent, err error) {
	//获取数据
	err = eventDB.GetInfo().GetInfoByID(id, &data)
	if err != nil {
		return
	}
	//反馈
	return
}

// getEventBySystem 通过系统来源获取数据
func getEventBySystem(fromSystem string, fromID int64, fromType string) (data FieldsEvent, err error) {
	//获取数据
	err = eventDB.GetInfo().GetInfoByFields(map[string]any{
		"from_system": fromSystem,
		"from_id":     fromID,
		"from_type":   fromType,
	}, true, &data)
	if err != nil {
		return
	}
	//反馈
	return
}
