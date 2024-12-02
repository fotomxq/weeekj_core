package AnalysisIndexValCustom

import "github.com/lib/pq"

type ArgsGetAnalysisIndexValCustomTotalAll struct {
	//指定编码
	CodeList pq.StringArray `db:"code_list" json:"codeList"`
	//时间范围
	StartAt string `db:"start_at" json:"startAt"`
	EndAt   string `db:"end_at" json:"endAt"`
}

// DataGetAnalysisIndexValCustomTotalAll 获取所有指标的总体值参数
type DataGetAnalysisIndexValCustomTotalAll struct {
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
}

// GetAnalysisIndexValCustomTotalAll 获取所有指标的总体值
// 不含维度筛选
func GetAnalysisIndexValCustomTotalAll(args *ArgsGetAnalysisIndexValCustomTotalAll) (dataList []DataGetAnalysisIndexValCustomTotalAll, err error) {
	if len(args.CodeList) < 1 {
		err = indexValCustomDB.GetClient().DB.GetPostgresql().Select(&dataList, "SELECT code, count(id) as data_count, min(year_md) as min_time, max(year_md) as max_time, ROUND(min(val_raw)::numeric, 4) as min_val, ROUND(max(val_raw)::numeric, 4) as max_val FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND ($1 = '' OR year_md >= $1) AND ($2 = '' OR year_md <= $2) AND extend1 = '' AND extend2 = '' AND extend3 = '' AND extend4 = '' AND extend5 = '' AND is_forecast = false GROUP BY code", args.StartAt, args.EndAt)
	} else {
		err = indexValCustomDB.GetClient().DB.GetPostgresql().Select(&dataList, "SELECT code, count(id) as data_count, min(year_md) as min_time, max(year_md) as max_time, ROUND(min(val_raw)::numeric, 4) as min_val, ROUND(max(val_raw)::numeric, 4) as max_val FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND ($1 = '' OR year_md >= $1) AND ($2 = '' OR year_md <= $2) AND extend1 = '' AND extend2 = '' AND extend3 = '' AND extend4 = '' AND extend5 = '' AND code = ANY($3) AND is_forecast = false GROUP BY code", args.StartAt, args.EndAt, args.CodeList)
	}
	return
}
