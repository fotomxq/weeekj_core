package BlogUserRead

import CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"

// ArgsMarge 聚合统计参数
type ArgsMarge struct {
	//多个标识码组
	Marks []ArgsMargeMark `json:"marks"`
	//商户ID
	OrgID int64 `json:"orgID"`
	//子商户ID
	ChildOrgID int64 `json:"childOrgID"`
}

type ArgsMargeMark struct {
	//标识码
	Mark string `json:"mark"`
	//指定渠道
	FromMark string `json:"fromMark"`
	//特定分类
	SortID int64 `json:"sortID"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//数量限制
	// 部分统计支持
	// 数据最多反馈1000条
	Limit int64 `json:"limit"`
}

// DataMarge 聚合统计反馈结构
type DataMarge struct {
	//数据结构
	Marks []DataMargeMark `json:"marks"`
}

type DataMargeMark struct {
	//标识码
	Mark string `json:"mark"`
	//指定渠道
	FromMark string `json:"fromMark"`
	//特定分类
	SortID int64 `json:"sortID"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//数量限制
	// 部分统计支持
	// 数据最多反馈1000条
	Limit int64 `json:"limit"`
	//数据集合
	Data interface{} `json:"data"`
}

// GetMarge 获取聚合统计
func GetMarge(args *ArgsMarge) (result DataMarge, err error) {
	//遍历聚合函数处理
	for _, vMark := range args.Marks {
		if vMark.Limit < 1 || vMark.Limit > 1000 {
			continue
		}
		vResult := DataMargeMark{
			Mark:        vMark.Mark,
			FromMark:    vMark.FromMark,
			SortID:      vMark.SortID,
			TimeBetween: vMark.TimeBetween,
			Limit:       vMark.Limit,
		}
		switch vMark.Mark {
		case "read_time":
			//阅读累计时间
			vResult.Data, _ = GetAnalysisTime(&ArgsGetAnalysisCount{
				OrgID:       args.OrgID,
				ChildOrgID:  args.ChildOrgID,
				UserID:      -1,
				FromMark:    vMark.FromMark,
				FromName:    "",
				IP:          "",
				SortID:      vMark.SortID,
				TimeBetween: vMark.TimeBetween,
			})
		case "read_count":
			//阅读总数
			vResult.Data, _ = GetAnalysisCount(&ArgsGetAnalysisCount{
				OrgID:       args.OrgID,
				ChildOrgID:  args.ChildOrgID,
				UserID:      -1,
				FromMark:    vMark.FromMark,
				FromName:    "",
				IP:          "",
				SortID:      vMark.SortID,
				TimeBetween: vMark.TimeBetween,
			})
		case "read_avg_time":
			//平均阅读时间
			vResult.Data, _ = GetAnalysisAvgReadTime(&ArgsGetAnalysisAvgReadTime{
				OrgID:       args.OrgID,
				ChildOrgID:  args.ChildOrgID,
				UserID:      -1,
				FromMark:    vMark.FromMark,
				FromName:    "",
				IP:          "",
				SortID:      vMark.SortID,
				ContentID:   -1,
				TimeBetween: vMark.TimeBetween,
			})
		case "read_sort_count":
			//不同分类下的统计数据
			vResult.Data, _ = GetAnalysisSortCount(&ArgsGetAnalysisSortCount{
				OrgID:       args.OrgID,
				ChildOrgID:  args.ChildOrgID,
				UserID:      -1,
				FromMark:    vMark.FromMark,
				FromName:    "",
				IP:          "",
				TimeBetween: vMark.TimeBetween,
			})
		}
		result.Marks = append(result.Marks, vResult)
	}
	return
}
