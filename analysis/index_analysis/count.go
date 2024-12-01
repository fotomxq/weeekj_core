package AnalysisIndexAnalysis

import (
	AnalysisIndex "github.com/fotomxq/weeekj_core/v5/analysis/index"
	AnalysisIndexFilter "github.com/fotomxq/weeekj_core/v5/analysis/index_filter"
	AnalysisIndexVal "github.com/fotomxq/weeekj_core/v5/analysis/index_val"
	"github.com/lib/pq"
)

// ArgsGetAnalysisIndexCount 获取所有指标的总体值参数
type ArgsGetAnalysisIndexCount struct {
	//指定编码
	CodeList pq.StringArray `db:"code_list" json:"codeList"`
	//时间范围
	StartAt string `db:"start_at" json:"startAt"`
	EndAt   string `db:"end_at" json:"endAt"`
}

// DataGetAnalysisIndexCount 获取所有指标的总体值参数
type DataGetAnalysisIndexCount struct {
	//指标编码
	Code string `db:"code" json:"code"`
	//数据量
	DataCount int64 `db:"data_count" json:"dataCount"`
	//最小时间
	MinTime string `db:"min_time" json:"minTime"`
	//最大时间
	MaxTime string `db:"max_time" json:"maxTime"`
	//数据最小值
	MinVal float64 `db:"min_val" json:"minVal"`
	//数据最大值
	MaxVal float64 `db:"max_val" json:"maxVal"`
	//是否为筛选指标
	// 筛选类指标部分情况下，不包含时间维度
	IsFilter bool `db:"is_filter" json:"isFilter"`
}

// GetAnalysisIndexCount 获取所有指标的总体值
// 不含维度筛选
func GetAnalysisIndexCount(args *ArgsGetAnalysisIndexCount) (dataList []DataGetAnalysisIndexCount) {
	//如果给空，则获取所有指标数据
	if len(args.CodeList) == 0 {
		//获取所有指标
		indexList := AnalysisIndex.GetIndexAllNoStruct()
		for _, v := range indexList {
			args.CodeList = append(args.CodeList, v.Code)
		}
	}
	//获取指标val的数据量
	valTotalList, _ := AnalysisIndexVal.GetAnalysisIndexValTotalAll(&AnalysisIndexVal.ArgsGetAnalysisIndexValTotalAll{
		CodeList: args.CodeList,
		StartAt:  args.StartAt,
		EndAt:    args.EndAt,
	})
	for _, v := range valTotalList {
		dataList = append(dataList, DataGetAnalysisIndexCount{
			Code:      v.Code,
			DataCount: v.DataCount,
			MinTime:   v.MinTime,
			MaxTime:   v.MaxTime,
			MinVal:    v.MinVal,
			MaxVal:    v.MaxVal,
			IsFilter:  false,
		})
	}
	//获取指标filter的数据量
	for _, v := range args.CodeList {
		isFind := false
		for k2, v2 := range dataList {
			if v == v2.Code && v2.DataCount > 0 {
				isFind = true
				continue
			}
			dataList[k2].DataCount = AnalysisIndexFilter.GetCount(v)
			dataList[k2].IsFilter = true
			break
		}
		if !isFind {
			dataList = append(dataList, DataGetAnalysisIndexCount{
				Code:      v,
				DataCount: AnalysisIndexFilter.GetCount(v),
				MinTime:   "",
				MaxTime:   "",
				MinVal:    0,
				MaxVal:    0,
				IsFilter:  true,
			})
		}
	}
	//反馈
	return
}
