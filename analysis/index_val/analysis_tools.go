package AnalysisIndexVal

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
}

// GetAnalysisIndexValTotal 计算某个指标存在维度值的总体值分布
// 只查询该指标下，维度存在数据，但
func GetAnalysisIndexValTotal(args *ArgsGetAnalysisIndexValTotal) (data DataGetAnalysisIndexValTotal) {
	_ = indexValDB.GetClient().DB.GetPostgresql().Get(&data, "SELECT count(id) as val_count, sum(val_raw) as sum_val, avg(val_raw) as avg_val, max(val_raw) as max_val, min(val_raw) as min_val FROM analysis_index_vals WHERE code = $1 AND year_md = $2 AND (extend1 != '' OR extend2 != '' OR extend3 != '' OR extend4 != '' OR extend5 != '') LIMIT 1", args.Code, args.YearMD)
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
}

// RefAnalysisIndexValTotal 指标总体值及归一化计算
// 注意：此方法不能用于预测值，预测值请使用CreateVal方法
func RefAnalysisIndexValTotal(args *ArgsRefAnalysisIndexValTotal) (err error) {
	topData := GetAnalysisIndexValTotal(&ArgsGetAnalysisIndexValTotal{
		Code:   args.Code,
		YearMD: args.YearMD,
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
	return
}
