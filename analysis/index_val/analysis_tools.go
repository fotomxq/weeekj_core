package AnalysisIndexVal

import (
	"errors"
	"github.com/lib/pq"
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
func GetAnalysisIndexValTotal(args *ArgsGetAnalysisIndexValTotal) (data DataGetAnalysisIndexValTotal) {
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
	_ = indexValDB.GetClient().DB.GetPostgresql().Get(&data, "SELECT count(id) as val_count, sum(val_raw) as sum_val, avg(val_raw) as avg_val, max(val_raw) as max_val, min(val_raw) as min_val FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND code = $1 AND year_md = $2 AND is_forecast = $3 "+sqlAppend+" LIMIT 1", args.Code, args.YearMD, args.IsForecast)
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
	topData := GetAnalysisIndexValTotal(&ArgsGetAnalysisIndexValTotal{
		Code:       args.Code,
		YearMD:     args.YearMD,
		ExcludeCol: args.ExcludeCol,
		IsForecast: false,
	})
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
		norVal := (topData.AvgVal - topData.MinVal) / (topData.MaxVal - topData.MinVal)
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
