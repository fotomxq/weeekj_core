package AnalysisIndexRFM

// ArgsGetRFMByCodeAndBetweenAt 获取指定时间范围的数据参数
type ArgsGetRFMByCodeAndBetweenAt struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
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
	//时间范围
	StartAt string `json:"startAt"`
	EndAt   string `json:"endAt"`
}

type DataGetRFMByCodeAndBetweenAt struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//年月
	YearM string `db:"year_m" json:"yearM" index:"true" field_list:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 维度和IndexVals模块一致
	///////////////////////////////////////////////////////////////////////////////////////////////////
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
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 计算过程值
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//R
	RVal float64 `db:"r_val" json:"rVal"`
	//RMin
	RMin float64 `db:"r_min" json:"rMin"`
	//RMax
	RMax float64 `db:"r_max" json:"rMax"`
	//R 权重
	RWeight float64 `db:"r_weight" json:"rWeight"`
	//F
	FVal float64 `db:"f_val" json:"fVal"`
	//FMin
	FMin float64 `db:"f_min" json:"fMin"`
	//FMax
	FMax float64 `db:"f_max" json:"fMax"`
	//F 权重
	FWeight float64 `db:"f_weight" json:"fWeight"`
	//M
	MVal float64 `db:"m_val" json:"mVal"`
	//MMin
	MMin float64 `db:"m_min" json:"mMin"`
	//MMax
	MMax float64 `db:"m_max" json:"mMax"`
	//M 权重
	MWeight float64 `db:"m_weight" json:"mWeight"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 计算结果
	///////////////////////////////////////////////////////////////////////////////////////////////////
	RFMResult float64 `db:"rfm_result" json:"rfmResult"`
}

// GetRFMByCodeAndBetweenAt 获取指定时间范围的数据
func GetRFMByCodeAndBetweenAt(args *ArgsGetRFMByCodeAndBetweenAt) (dataList []DataGetRFMByCodeAndBetweenAt) {
	_ = rfmDB.GetClient().DB.GetPostgresql().Select(&dataList, "SELECT code, year_m, extend1, extend2, extend3, extend4, extend5, r_val, r_min, r_max, r_weight, f_val, f_min, f_max, f_weight, m_val, m_min, m_max, m_weight, rfm_result FROM analysis_index_rfm WHERE code = $1 AND extend1 = $2 AND extend2 = $3 AND extend3 = $4 AND extend4 = $5 AND extend5 = $6 AND delete_at < to_timestamp(1000000) AND year_m >= $7 AND year_m <= $8", args.Code, args.Extend1, args.Extend2, args.Extend3, args.Extend4, args.Extend5, args.StartAt, args.EndAt)
	return
}

type DataGetRFMByCodeAndBetweenAtResult struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//年月
	YearM string `db:"year_m" json:"yearM" index:"true" field_list:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 维度和IndexVals模块一致
	///////////////////////////////////////////////////////////////////////////////////////////////////
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
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 计算结果
	///////////////////////////////////////////////////////////////////////////////////////////////////
	RFMResult float64 `db:"rfm_result" json:"rfmResult"`
}

// GetRFMByCodeAndBetweenAtResult 获取指定时间范围的数据结果
// 仅反馈结果
func GetRFMByCodeAndBetweenAtResult(args *ArgsGetRFMByCodeAndBetweenAt) (dataList []DataGetRFMByCodeAndBetweenAt) {
	_ = rfmDB.GetClient().DB.GetPostgresql().Select(&dataList, "SELECT code, year_m, extend1, extend2, extend3, extend4, extend5, rfm_result FROM analysis_index_rfm WHERE code = $1 AND extend1 = $2 AND extend2 = $3 AND extend3 = $4 AND extend4 = $5 AND extend5 = $6 AND delete_at < to_timestamp(1000000) AND year_m >= $7 AND year_m <= $8", args.Code, args.Extend1, args.Extend2, args.Extend3, args.Extend4, args.Extend5, args.StartAt, args.EndAt)
	return
}

// ArgsGetRFMByCodeAndYMAndExtend 获取指定编码和日期的RFM数据参数
type ArgsGetRFMByCodeAndYMAndExtend struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//年月
	YearM string `db:"year_m" json:"yearM" index:"true" field_list:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 维度和IndexVals模块一致
	///////////////////////////////////////////////////////////////////////////////////////////////////
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
}

// GetRFMByCodeAndYMAndExtend 获取指定编码和日期的RFM数据
func GetRFMByCodeAndYMAndExtend(args *ArgsGetRFMByCodeAndYMAndExtend) (result float64) {
	data := getRFMByCodeAndYMAndExtendRaw(args)
	if data.ID > 0 {
		result = data.RFMResult
	}
	return
}

// 获取RFM指定条件的原始数据
func getRFMByCodeAndYMAndExtendRaw(args *ArgsGetRFMByCodeAndYMAndExtend) (result FieldsRFM) {
	_ = rfmDB.GetInfo().GetInfoByFields(map[string]any{
		"code":    args.Code,
		"year_m":  args.YearM,
		"extend1": args.Extend1,
		"extend2": args.Extend2,
		"extend3": args.Extend3,
		"extend4": args.Extend4,
		"extend5": args.Extend5,
	}, true, &result)
	return
}