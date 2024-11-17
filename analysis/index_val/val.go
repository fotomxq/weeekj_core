package AnalysisIndexVal

// ArgsGetValsByBetweenAt 获取指定时间范围内的指标值参数
type ArgsGetValsByBetweenAt struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//时间范围
	StartAt string `json:"startAt"`
	EndAt   string `json:"endAt"`
	//扩展维度1
	// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
	Extend1 string `db:"extend1" json:"extend1" index:"true"`
	//扩展维度2
	Extend2 string `db:"extend2" json:"extend2" index:"true"`
	//扩展维度3
	Extend3 string `db:"extend3" json:"extend3" index:"true"`
	//扩展维度4
	Extend4 string `db:"extend4" json:"extend4" index:"true"`
	//扩展维度5
	Extend5 string `db:"extend5" json:"extend5" index:"true"`
	//原始值范围
	ValRawMin float64 `db:"val_raw_min" json:"valRawMin"`
	ValRawMax float64 `db:"val_raw_max" json:"valRawMax"`
	//归一化值范围
	ValNormMin float64 `db:"val_norm_min" json:"valNormMin"`
	ValNormMax float64 `db:"val_norm_max" json:"valNormMax"`
}

// DataGetValsByBetweenAt 获取指定时间范围内的指标值数据
type DataGetValsByBetweenAt struct {
	//年月日
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	YearMD string `db:"year_md" json:"yearMD" index:"true"`
	//扩展维度1
	// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
	Extend1 string `db:"extend1" json:"extend1" index:"true"`
	//扩展维度2
	Extend2 string `db:"extend2" json:"extend2" index:"true"`
	//扩展维度3
	Extend3 string `db:"extend3" json:"extend3" index:"true"`
	//扩展维度4
	Extend4 string `db:"extend4" json:"extend4" index:"true"`
	//扩展维度5
	Extend5 string `db:"extend5" json:"extend5" index:"true"`
	//原始值
	ValRaw float64 `db:"val_raw" json:"valRaw" index:"true"`
	//归一化值
	ValNorm float64 `db:"val_norm" json:"valNorm" index:"true"`
	//是否为预测值
	IsForecast bool `db:"is_forecast" json:"isForecast"`
}

// GetValsByBetweenAt 获取指定时间范围内的指标值
func GetValsByBetweenAt(args *ArgsGetValsByBetweenAt) (dataList []DataGetValsByBetweenAt, err error) {
	err = indexValDB.GetClient().DB.GetPostgresql().Select(&dataList, "SELECT year_md, extend1, extend2, extend3, extend4, extend5, val_raw, val_norm, is_forecast FROM analysis_index_vals WHERE code = $1 AND ($2 != '' OR year_md >= $2) AND ($3 != '' OR year_md <= $3) AND extend1 = $4 AND extend2 = $5 AND extend3 = $6 AND extend4 = $7 AND extend5 = $8 AND ($9 < 0 OR (val_raw >= $9 AND val_raw <= $10)) AND ($11 < 0 OR (val_norm >= $11 AND val_norm <= $12)) ORDER BY year_md", args.Code, args.StartAt, args.EndAt, args.Extend1, args.Extend2, args.Extend3, args.Extend4, args.Extend5, args.ValRawMin, args.ValRawMax, args.ValNormMin, args.ValNormMax)
	if err != nil {
		return
	}
	return
}

// ArgsGetValsByFilter 获取指定的数据
type ArgsGetValsByFilter struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//年月日
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	YearMD string `db:"year_md" json:"yearMD" index:"true"`
	//扩展维度1
	// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
	Extend1 string `db:"extend1" json:"extend1" index:"true"`
	//扩展维度2
	Extend2 string `db:"extend2" json:"extend2" index:"true"`
	//扩展维度3
	Extend3 string `db:"extend3" json:"extend3" index:"true"`
	//扩展维度4
	Extend4 string `db:"extend4" json:"extend4" index:"true"`
	//扩展维度5
	Extend5 string `db:"extend5" json:"extend5" index:"true"`
}

// GetValsByFilter 获取指定的数据
// 注意判断YearMD，如果为空则没有数据
func GetValsByFilter(args *ArgsGetValsByFilter) (data DataGetValsByBetweenAt) {
	var rawData FieldsVal
	_ = indexValDB.GetInfo().GetInfoByFields(map[string]any{
		"code":    args.Code,
		"year_md": args.YearMD,
		"extend1": args.Extend1,
		"extend2": args.Extend2,
		"extend3": args.Extend3,
		"extend4": args.Extend4,
		"extend5": args.Extend5,
	}, true, &rawData)
	data = DataGetValsByBetweenAt{
		YearMD:  rawData.YearMD,
		Extend1: rawData.Extend1,
		Extend2: rawData.Extend2,
		Extend3: rawData.Extend3,
		Extend4: rawData.Extend4,
		Extend5: rawData.Extend5,
	}
	return
}

// ArgsCreateVal 插入新的统计数据参数
type ArgsCreateVal struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//年月日
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	YearMD string `db:"year_md" json:"yearMD" index:"true"`
	//扩展维度1
	// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
	Extend1 string `db:"extend1" json:"extend1" index:"true"`
	//扩展维度2
	Extend2 string `db:"extend2" json:"extend2" index:"true"`
	//扩展维度3
	Extend3 string `db:"extend3" json:"extend3" index:"true"`
	//扩展维度4
	Extend4 string `db:"extend4" json:"extend4" index:"true"`
	//扩展维度5
	Extend5 string `db:"extend5" json:"extend5" index:"true"`
	//原始值
	ValRaw float64 `db:"val_raw" json:"valRaw" index:"true"`
	//是否为预测值
	IsForecast bool `db:"is_forecast" json:"isForecast"`
}

// CreateVal 插入新的统计数据
func CreateVal(args *ArgsCreateVal) (err error) {
	//获取数据
	var rawData FieldsVal
	_ = indexValDB.GetInfo().GetInfoByFields(map[string]any{
		"code":    args.Code,
		"year_md": args.YearMD,
		"extend1": args.Extend1,
		"extend2": args.Extend2,
		"extend3": args.Extend3,
		"extend4": args.Extend4,
		"extend5": args.Extend5,
	}, true, &rawData)
	if rawData.YearMD == "" {
		//插入数据
		type insertType struct {
			//指标编码
			Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
			//年月日
			// 可任意持续，如年，或仅年月
			// 不建议构建小时及以下级别的指标
			// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
			YearMD string `db:"year_md" json:"yearMD" index:"true"`
			//扩展维度1
			// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
			Extend1 string `db:"extend1" json:"extend1" index:"true"`
			//扩展维度2
			Extend2 string `db:"extend2" json:"extend2" index:"true"`
			//扩展维度3
			Extend3 string `db:"extend3" json:"extend3" index:"true"`
			//扩展维度4
			Extend4 string `db:"extend4" json:"extend4" index:"true"`
			//扩展维度5
			Extend5 string `db:"extend5" json:"extend5" index:"true"`
			//原始值
			ValRaw float64 `db:"val_raw" json:"valRaw" index:"true"`
			//归一化值
			ValNorm float64 `db:"val_norm" json:"valNorm" index:"true"`
			//是否为预测值
			IsForecast bool `db:"is_forecast" json:"isForecast"`
		}
		insertData := insertType{
			Code:       args.Code,
			YearMD:     args.YearMD,
			Extend1:    args.Extend1,
			Extend2:    args.Extend2,
			Extend3:    args.Extend3,
			Extend4:    args.Extend4,
			Extend5:    args.Extend5,
			ValRaw:     args.ValRaw,
			ValNorm:    0,
			IsForecast: args.IsForecast,
		}
		rawData.ID, err = indexValDB.GetInsert().InsertRow(&insertData)
		if err != nil {
			return
		}
	} else {
		//更新数据
		type updateType struct {
			// ID
			ID int64 `db:"id" json:"id" check:"id" unique:"true"`
			//原始值
			ValRaw float64 `db:"val_raw" json:"valRaw" index:"true"`
			//是否为预测值
			IsForecast bool `db:"is_forecast" json:"isForecast"`
		}
		updateData := updateType{
			ID:         rawData.ID,
			ValRaw:     args.ValRaw,
			IsForecast: args.IsForecast,
		}
		err = indexValDB.GetUpdate().UpdateByID(&updateData)
		if err != nil {
			return
		}
	}
	//修正归一化值
	_, _ = indexValDB.GetClient().DB.GetPostgresql().Exec("UPDATE analysis_index_vals SET val_norm = (val_raw - (SELECT MIN(val_raw) FROM analysis_index_vals WHERE code = $1 AND year_md = $2 AND extend1 = $3 AND extend2 = $4 AND extend3 = $5 AND extend4 = $6 AND extend5 = $7)) / ((SELECT MAX(val_raw) FROM analysis_index_vals WHERE code = $1 AND year_md = $2 AND extend1 = $3 AND extend2 = $4 AND extend3 = $5 AND extend4 = $6 AND extend5 = $7) - (SELECT MIN(val_raw) FROM analysis_index_vals WHERE code = $1 AND year_md = $2 AND extend1 = $3 AND extend2 = $4 AND extend3 = $5 AND extend4 = $6 AND extend5 = $7)) WHERE code = $1 AND year_md = $2 AND extend1 = $3 AND extend2 = $4 AND extend3 = $5 AND extend4 = $6 AND extend5 = $7", args.Code, args.YearMD, args.Extend1, args.Extend2, args.Extend3, args.Extend4, args.Extend5)
	//返回
	return
}
