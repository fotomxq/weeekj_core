package AnalysisIndexVal

import (
	"errors"
	"fmt"
	AnalysisIndex "github.com/fotomxq/weeekj_core/v5/analysis/index"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/lib/pq"
	"time"
)

// DataGetAnalysisIndexValTotal 计算某个指标存在维度值的总体值分布数据
type DataGetAnalysisIndexValTotal struct {
	//次数
	ValCount int64 `db:"val_count" json:"valCount"`
	//累计值
	SumVal float64 `db:"sum_val" json:"sumVal"`
	//平均值
	AvgVal float64 `db:"avg_val" json:"avgVal"`
	//最大值
	MaxVal float64 `db:"max_val" json:"maxVal"`
	//最小值
	MinVal float64 `db:"min_val" json:"minVal"`
}

type ArgsGetAnalysisIndexValTotal struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//年月日
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	YearMD string `db:"year_md" json:"yearMD" index:"true"`
	//排除列
	ExcludeCol pq.StringArray `db:"exclude_col" json:"excludeCol"`
	//是否为预测值
	IsForecast bool `db:"is_forecast" json:"isForecast"`
}

// GetAnalysisIndexValTotal 计算某个指标存在维度值的总体值分布
// 只查询该指标下，维度存在数据，但
func GetAnalysisIndexValTotal(args *ArgsGetAnalysisIndexValTotal) (data DataGetAnalysisIndexValTotal, err error) {
	//检查code是否存在
	var codeData AnalysisIndex.FieldsIndex
	codeData, err = AnalysisIndex.GetIndexByCode(args.Code)
	if err != nil || codeData.ID < 1 {
		err = errors.New(fmt.Sprint("code not exist, ", err))
		return
	}
	//组装筛选条件
	sqlAppend := ""
	if len(args.ExcludeCol) > 0 {
		sqlAppend = " AND ("
		for i, v := range args.ExcludeCol {
			if i > 0 {
				sqlAppend += " OR "
			}
			sqlAppend += v + " != ''"
		}
		sqlAppend += ")"
	}
	err = indexValDB.GetClient().DB.GetPostgresql().Get(&data, "SELECT count(id) as val_count, sum(val_raw) as sum_val, avg(val_raw) as avg_val, max(val_raw) as max_val, min(val_raw) as min_val FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND code = $1 AND year_md = $2 AND is_forecast = $3 "+sqlAppend+" LIMIT 1", args.Code, args.YearMD, args.IsForecast)
	if err != nil {
		err = errors.New(fmt.Sprint("get data failed, ", err))
		return
	}
	return
}

// GetAnalysisIndexValTotalCount 获取指标计算的数量
func GetAnalysisIndexValTotalCount(code string, afterAt time.Time) (count int64) {
	//检查code是否存在
	codeData, _ := AnalysisIndex.GetIndexByCode(code)
	if codeData.ID < 1 {
		return
	}
	//获取数据
	_ = indexValDB.GetClient().DB.GetPostgresql().Get(&count, "SELECT count(id) FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND code = $1 AND year_md >= $2 LIMIT 1", code, afterAt.Format("2006-01-02"))
	//反馈
	return
}

// GetAnalysisIndexValTotalCountByBetweenAt 获取指标在指定时间断内的数据量
func GetAnalysisIndexValTotalCountByBetweenAt(code string, startTime time.Time, endTime time.Time) (count int64) {
	//检查code是否存在
	codeData, _ := AnalysisIndex.GetIndexByCode(code)
	if codeData.ID < 1 {
		return
	}
	//获取数据
	_ = indexValDB.GetClient().DB.GetPostgresql().Get(&count, "SELECT count(id) FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND code = $1 AND year_md >= $2 AND year_md <= $3 LIMIT 1", code, startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
	//反馈
	return
}

type ArgsGetAnalysisIndexValTotalAll struct {
	//指定编码
	CodeList pq.StringArray `db:"code_list" json:"codeList"`
	//时间范围
	StartAt string `db:"start_at" json:"startAt"`
	EndAt   string `db:"end_at" json:"endAt"`
}

// DataGetAnalysisIndexValTotalAll 获取所有指标的总体值参数
type DataGetAnalysisIndexValTotalAll struct {
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

// GetAnalysisIndexValTotalAll 获取所有指标的总体值
// 不含维度筛选
func GetAnalysisIndexValTotalAll(args *ArgsGetAnalysisIndexValTotalAll) (dataList []DataGetAnalysisIndexValTotalAll, err error) {
	if len(args.CodeList) < 1 {
		_ = indexValDB.GetClient().DB.GetPostgresql().Select(&dataList, "SELECT code, count(id) as data_count, min(year_md) as min_time, max(year_md) as max_time, min(val_raw) as min_val, max(val_raw) as max_val FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND ($1 = '' OR year_md >= $1) AND ($2 = '' OR year_md <= $2) AND extend1 = '' AND extend2 = '' AND extend3 = '' AND extend4 = '' AND extend5 = '' GROUP BY code", args.StartAt, args.EndAt)
	} else {
		_ = indexValDB.GetClient().DB.GetPostgresql().Select(&dataList, "SELECT code, count(id) as data_count, min(year_md) as min_time, max(year_md) as max_time, min(val_raw) as min_val, max(val_raw) as max_val FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND ($1 = '' OR year_md >= $1) AND ($2 = '' OR year_md <= $2) AND extend1 = '' AND extend2 = '' AND extend3 = '' AND extend4 = '' AND extend5 = '' AND code = ANY($3) GROUP BY code", args.StartAt, args.EndAt, args.CodeList)
	}
	return
}

// GetEventExtendDistinctList 获取指定维度的所有可选值
func GetEventExtendDistinctList(extendNum int) (dataList []string, err error) {
	//获取数据
	dataList, err = indexValDB.GetList().GetDistinctList(fmt.Sprintf("extend%d", extendNum))
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsRefAnalysisIndexValTotal 指标总体值及归一化计算参数
type ArgsRefAnalysisIndexValTotal struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//年月日
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	YearMD string `db:"year_md" json:"yearMD" index:"true"`
	//计算方式
	// count/sum/avg/max/min
	CalcType string `db:"calc_type" json:"calcType"`
	//排除列
	ExcludeCol pq.StringArray `db:"exclude_col" json:"excludeCol"`
}

// RefAnalysisIndexValTotal 指标总体值及归一化计算
// 注意：此方法不能用于预测值，预测值请使用CreateVal方法; 归一化处理请使用reviseNormVal方法
func RefAnalysisIndexValTotal(args *ArgsRefAnalysisIndexValTotal) (err error) {
	//检查code是否存在
	codeData, _ := AnalysisIndex.GetIndexByCode(args.Code)
	if codeData.ID < 1 {
		return
	}
	//获取总体值
	var topData DataGetAnalysisIndexValTotal
	topData, err = GetAnalysisIndexValTotal(&ArgsGetAnalysisIndexValTotal{
		Code:       args.Code,
		YearMD:     args.YearMD,
		ExcludeCol: args.ExcludeCol,
		IsForecast: false,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("get top data failed, ", err))
		return
	}
	//计算值
	var valRaw float64
	switch args.CalcType {
	case "count":
		valRaw = float64(topData.ValCount)
	case "sum":
		valRaw = topData.SumVal
	case "avg":
		valRaw = topData.AvgVal
	case "max":
		valRaw = topData.MaxVal
	case "min":
		valRaw = topData.MinVal
	default:
		err = errors.New("calc_type error")
		return
	}
	//插入新的数据
	err = CreateVal(&ArgsCreateVal{
		Code:       args.Code,
		YearMD:     args.YearMD,
		Extend1:    "",
		Extend2:    "",
		Extend3:    "",
		Extend4:    "",
		Extend5:    "",
		ValRaw:     valRaw,
		IsForecast: false,
	})
	if err != nil {
		return
	}
	//如果args.CalcType为均值，自动归一化处理
	if args.CalcType == "avg" {
		var norVal float64
		if (topData.MaxVal - topData.MinVal) != 0 {
			norVal = CoreFilter.RoundToTwoDecimalPlaces(((topData.AvgVal - topData.MinVal) / (topData.MaxVal - topData.MinVal)) * 100)
		}
		err = reviseNormVal(&argsReviseNormVal{
			Code:       args.Code,
			YearMD:     args.YearMD,
			ValNorm:    norVal,
			IsForecast: false,
		})
		if err != nil {
			return
		}
	}
	//反馈
	return
}
