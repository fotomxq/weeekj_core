package AnalysisIndexRFM

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreMathRFM "github.com/fotomxq/weeekj_core/v5/core/math/rfm"
	"time"
)

// ArgsCalcRFM 计算RFM值并记录
type ArgsCalcRFM struct {
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
}

// CalcRFM 计算RFM值并记录
// 注意RFM计算中，R指标需根据需求提前根据需求反转（也可以不反转），否则计算结果将出现异常
func CalcRFM(args *ArgsCalcRFM) (err error) {
	//计算RFM值
	var rfmCore CoreMathRFM.Core
	//约定权重
	rfmCore.SetWeight([]CoreMathRFM.Weight{
		//参数给定值
		{
			Number: 0,
			R:      args.RWeight,
			F:      args.FWeight,
			M:      args.MWeight,
		},
		//默认预设值
		{
			Number: 1,
			R:      0.3,
			F:      0.3,
			M:      0.4,
		},
	})
	//是否使用预设权重
	useDefaultWeight := false
	//获取权重
	weightData := rfmCore.GetWeight(0)
	if weightData.F < 1 {
		//设置默认值
		weightData = rfmCore.GetWeight(1)
		useDefaultWeight = true
	}
	//设置范围
	rfmCore.SetDataRange(args.RMin, args.RMax, args.FMin, args.FMax, args.MMin, args.MMax)
	//计算结果
	var rfmResult float64
	if useDefaultWeight {
		rfmResult = rfmCore.GetScoreByWeight(args.RVal, args.FVal, args.MVal, 1)
	} else {
		rfmResult = rfmCore.GetScoreByWeight(args.RVal, args.FVal, args.MVal, 0)
	}
	//修正结果
	rfmResult = CoreFilter.RoundTo4DecimalPlaces(rfmResult)
	//获取数据
	rawData := getRFMByCodeAndYMAndExtendRaw(&ArgsGetRFMByCodeAndYMAndExtend{
		Code:    args.Code,
		YearM:   args.YearM,
		Extend1: args.Extend1,
		Extend2: args.Extend2,
		Extend3: args.Extend3,
		Extend4: args.Extend4,
		Extend5: args.Extend5,
	})
	//更新数据
	if rawData.ID > 0 {
		type updateType struct {
			// ID
			ID int64 `db:"id" json:"id" check:"id" unique:"true"`
			//更新时间
			UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
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
		updateData := updateType{
			ID:        rawData.ID,
			UpdateAt:  CoreFilter.GetNowTime(),
			RVal:      args.RVal,
			RMin:      args.RMin,
			RMax:      args.RMax,
			RWeight:   weightData.R,
			FVal:      args.FVal,
			FMin:      args.FMin,
			FMax:      args.FMax,
			FWeight:   weightData.F,
			MVal:      args.MVal,
			MMin:      args.MMin,
			MMax:      args.MMax,
			MWeight:   weightData.M,
			RFMResult: CoreFilter.RoundToTwoDecimalPlaces(rfmResult),
		}
		err = rfmDB.GetUpdate().UpdateByID(&updateData)
		if err != nil {
			return
		}
	} else {
		type insertType struct {
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
		insertData := insertType{
			Code:      args.Code,
			YearM:     args.YearM,
			Extend1:   args.Extend1,
			Extend2:   args.Extend2,
			Extend3:   args.Extend3,
			Extend4:   args.Extend4,
			Extend5:   args.Extend5,
			RVal:      args.RVal,
			RMin:      args.RMin,
			RMax:      args.RMax,
			RWeight:   weightData.R,
			FVal:      args.FVal,
			FMin:      args.FMin,
			FMax:      args.FMax,
			FWeight:   weightData.F,
			MVal:      args.MVal,
			MMin:      args.MMin,
			MMax:      args.MMax,
			MWeight:   weightData.M,
			RFMResult: CoreFilter.RoundToTwoDecimalPlaces(rfmResult),
		}
		_, err = rfmDB.GetInsert().InsertRow(&insertData)
		if err != nil {
			return
		}
	}
	//反馈
	return
}
