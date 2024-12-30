package AnalysisIndexVal

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/golang-module/carbon"
	"math"
)

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
	err = indexValDB.GetClient().DB.GetPostgresql().Select(&dataList, "SELECT year_md, extend1, extend2, extend3, extend4, extend5, val_raw, val_norm, is_forecast FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND code = $1 AND ($2 = '' OR year_md >= $2) AND ($3 = '' OR year_md <= $3) AND extend1 = $4 AND extend2 = $5 AND extend3 = $6 AND extend4 = $7 AND extend5 = $8 AND ($9 < 0 OR (val_raw >= $9 AND val_raw <= $10)) AND ($11 < 0 OR (val_norm >= $11 AND val_norm <= $12)) ORDER BY year_md", args.Code, args.StartAt, args.EndAt, args.Extend1, args.Extend2, args.Extend3, args.Extend4, args.Extend5, args.ValRawMin, args.ValRawMax, args.ValNormMin, args.ValNormMax)
	if err != nil {
		return
	}
	return
}

// ArgsGetValsByBetweenAtAndAutoFull 获取指定时间范围的多指标数据集参数
type ArgsGetValsByBetweenAtAndAutoFull struct {
	//时间范围
	StartAt string `json:"startAt"`
	EndAt   string `json:"endAt"`
	//时间类型
	// 可选值：year, month, day
	DateType string `json:"dateType"`
	//指标参数集合
	// 允许给予不同指标、不同筛选条件
	IndexList []ArgsGetValsByBetweenAtAndAutoFullChild `json:"indexList"`
}

type ArgsGetValsByBetweenAtAndAutoFullChild struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
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

// DataGetValsByBetweenAtAndAutoFull 获取指定时间范围的多指标数据集数据
type DataGetValsByBetweenAtAndAutoFull struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
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
	//数据集合
	DataList []DataGetValsByBetweenAtAndAutoFullChild `db:"data_list" json:"dataList"`
}

type DataGetValsByBetweenAtAndAutoFullChild struct {
	//年月日
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	YearMD string `db:"year_md" json:"yearMD" index:"true"`
	//原始值
	ValRaw float64 `db:"val_raw" json:"valRaw" index:"true"`
	//归一化值
	ValNorm float64 `db:"val_norm" json:"valNorm" index:"true"`
	//是否为预测值
	IsForecast bool `db:"is_forecast" json:"isForecast"`
}

// GetValsByBetweenAtAndAutoFull 获取指定时间范围的多指标数据集
func GetValsByBetweenAtAndAutoFull(args *ArgsGetValsByBetweenAtAndAutoFull) (dataList []DataGetValsByBetweenAtAndAutoFull, err error) {
	//缓冲器设计
	//cacheMark := fmt.Sprint("analysis:index:vals:between.auto.full:args:", CoreFilter.GetMd5StrByStr(fmt.Sprint(args)))
	//if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
	//	return
	//}
	//dataList = []DataGetValsByBetweenAtAndAutoFull{}
	//参数检查
	switch args.DateType {
	case "year":
	case "month":
	case "day":
	default:
		err = errors.New("no support between time")
		return
	}
	//定义时间结构
	var startAt, endAt carbon.Carbon
	if args.StartAt != "" {
		switch args.DateType {
		case "year":
			startAt, err = CoreFilter.GetTimeCarbonByDefault(args.StartAt + "-01-01")
		case "month":
			startAt, err = CoreFilter.GetTimeCarbonByDefault(args.StartAt + "-01")
		case "day":
			startAt, err = CoreFilter.GetTimeCarbonByDefault(args.StartAt)
		}
		if err != nil {
			err = errors.New(fmt.Sprint("arg start at error, ", err, ", ", args.StartAt))
			return
		}
		startAt = startAt.StartOfMonth()
	}
	if args.EndAt != "" {
		switch args.DateType {
		case "year":
			endAt, err = CoreFilter.GetTimeCarbonByDefault(args.EndAt + "-01-01")
		case "month":
			endAt, err = CoreFilter.GetTimeCarbonByDefault(args.EndAt + "-01")
		case "day":
			endAt, err = CoreFilter.GetTimeCarbonByDefault(args.EndAt)
		}
		if err != nil {
			err = errors.New(fmt.Sprint("arg end at error, ", err, ", ", args.EndAt))
			return
		}
		endAt = endAt.EndOfMonth()
	}
	if args.StartAt != "" && args.EndAt != "" {
		if startAt.Gte(endAt) {
			//时间范围异常，不进行处理
			err = errors.New(fmt.Sprint("arg start and end at error, ", err))
			return
		}
	}
	if args.StartAt == "" && args.EndAt == "" {
		startAt = CoreFilter.GetNowTimeCarbon()
		endAt = startAt.SubYears(5)
		switch args.DateType {
		case "year":
			endAt = startAt.EndOfYear()
			startAt = endAt.SubYears(5).StartOfYear()
		case "month":
			endAt = startAt.EndOfYear()
			startAt = endAt.SubYears(5).StartOfYear()
		case "day":
			endAt = startAt.EndOfYear()
			startAt = endAt.StartOfYear()
		}
		args.StartAt = startAt.Format("2006-01-02")
		args.EndAt = endAt.Format("2006-01-02")
	}
	//获取数据集
	for _, vIndex := range args.IndexList {
		//是否写入过数据，用于后续的时间排序处理
		haveAppendData := false
		//如果没有时间范围，则按照数据集合的时间范围
		var rawList []DataGetValsByBetweenAt
		rawList, _ = GetValsByBetweenAt(&ArgsGetValsByBetweenAt{
			Code:       vIndex.Code,
			StartAt:    args.StartAt,
			EndAt:      args.EndAt,
			Extend1:    vIndex.Extend1,
			Extend2:    vIndex.Extend2,
			Extend3:    vIndex.Extend3,
			Extend4:    vIndex.Extend4,
			Extend5:    vIndex.Extend5,
			ValRawMin:  vIndex.ValRawMin,
			ValRawMax:  vIndex.ValRawMax,
			ValNormMin: vIndex.ValNormMin,
			ValNormMax: vIndex.ValNormMax,
		})
		//准备写入的数据
		vDataList := make([]DataGetValsByBetweenAtAndAutoFullChild, 0)
		//设定开始时间
		stepAt := startAt
		//如果是否存在数据，都会进行时间维度的补全处理
		for {
			var vYM string
			switch args.DateType {
			case "year":
				vYM = stepAt.Format("2006")
			case "month":
				vYM = stepAt.Format("2006-01")
			case "day":
				vYM = stepAt.Format("2006-01-02")
			}
			//补全数据
			isFind := false
			for _, vData := range rawList {
				if vYM != vData.YearMD {
					continue
				}
				vDataList = append(vDataList, DataGetValsByBetweenAtAndAutoFullChild{
					YearMD:     vData.YearMD,
					ValRaw:     vData.ValRaw,
					ValNorm:    vData.ValNorm,
					IsForecast: vData.IsForecast,
				})
				isFind = true
				break
			}
			if !isFind {
				haveAppendData = true
				vDataList = append(vDataList, DataGetValsByBetweenAtAndAutoFullChild{
					YearMD:     vYM,
					ValRaw:     0,
					ValNorm:    0,
					IsForecast: false,
				})
			}
			//下一个月
			stepAt = stepAt.AddMonth()
			if stepAt.Gt(endAt) {
				break
			}
		}
		//如果haveAppendData为true，说明有数据补全，需要重新排序
		if haveAppendData {
			//排序
			for i := 0; i < len(vDataList); i++ {
				for j := i + 1; j < len(vDataList); j++ {
					if vDataList[i].YearMD > vDataList[j].YearMD {
						vDataList[i], vDataList[j] = vDataList[j], vDataList[i]
					}
				}
			}
		}
		//写入数据
		dataList = append(dataList, DataGetValsByBetweenAtAndAutoFull{
			Code:     vIndex.Code,
			Extend1:  vIndex.Extend1,
			Extend2:  vIndex.Extend2,
			Extend3:  vIndex.Extend3,
			Extend4:  vIndex.Extend4,
			Extend5:  vIndex.Extend5,
			DataList: vDataList,
		})
	}
	//存储缓冲
	//Router2SystemConfig.MainCache.DeleteMark(cacheMark)
	//Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, CoreCache.CacheTime1Hour)
	//反馈
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
	//是否为预测值
	IsForecast bool `db:"is_forecast" json:"isForecast"`
}

// GetValsByFilter 获取指定的数据
// 注意判断YearMD，如果为空则没有数据
func GetValsByFilter(args *ArgsGetValsByFilter) (data FieldsVal) {
	_ = indexValDB.GetClient().DB.GetPostgresql().Get(&data, "SELECT * FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND code = $1 AND year_md = $2 AND extend1 = $3 AND extend2 = $4 AND extend3 = $5 AND extend4 = $6 AND extend5 = $7 AND is_forecast = $8", args.Code, args.YearMD, args.Extend1, args.Extend2, args.Extend3, args.Extend4, args.Extend5, args.IsForecast)
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
	_ = indexValDB.GetClient().DB.GetPostgresql().Get(&rawData, "SELECT * FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND code = $1 AND year_md = $2 AND extend1 = $3 AND extend2 = $4 AND extend3 = $5 AND extend4 = $6 AND extend5 = $7 AND is_forecast = $8", args.Code, args.YearMD, args.Extend1, args.Extend2, args.Extend3, args.Extend4, args.Extend5, args.IsForecast)
	//修正浮点数
	if math.IsNaN(rawData.ValRaw) || math.IsInf(rawData.ValRaw, 0) {
		rawData.ValRaw = 0
	}
	//修正结果
	rawData.ValRaw = CoreFilter.RoundTo4DecimalPlaces(rawData.ValRaw)
	if rawData.ID < 1 {
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
		}
		updateData := updateType{
			ID:     rawData.ID,
			ValRaw: args.ValRaw,
		}
		err = indexValDB.GetUpdate().UpdateByID(&updateData)
		if err != nil {
			return
		}
	}
	//返回
	return
}

type ArgsUpdateNormalValByCode struct {
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
	//归一化值
	ValNorm float64 `db:"val_norm" json:"valNorm" index:"true"`
	//是否为预测值
	IsForecast bool `db:"is_forecast" json:"isForecast"`
}

// UpdateNormalValByCode 修订指定数据的归一化值
func UpdateNormalValByCode(args *ArgsUpdateNormalValByCode) (err error) {
	//获取数据
	var rawData FieldsVal
	_ = indexValDB.GetClient().DB.GetPostgresql().Get(&rawData, "SELECT * FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND code = $1 AND year_md = $2 AND extend1 = $3 AND extend2 = $4 AND extend3 = $5 AND extend4 = $6 AND extend5 = $7 AND is_forecast = $8", args.Code, args.YearMD, args.Extend1, args.Extend2, args.Extend3, args.Extend4, args.Extend5, args.IsForecast)
	if rawData.ID < 1 {
		err = errors.New(fmt.Sprint("no data, ", err))
		return
	}
	//更新数据
	type updateType struct {
		// ID
		ID int64 `db:"id" json:"id" check:"id" unique:"true"`
		//归一化值
		ValNorm float64 `db:"val_norm" json:"valNorm" index:"true"`
	}
	updateData := updateType{
		ID:      rawData.ID,
		ValNorm: args.ValNorm,
	}
	//修正浮点数
	if math.IsNaN(updateData.ValNorm) || math.IsInf(updateData.ValNorm, 0) {
		updateData.ValNorm = 0
	}
	//修正结果
	updateData.ValNorm = CoreFilter.RoundTo4DecimalPlaces(updateData.ValNorm)
	//执行修改
	err = indexValDB.GetUpdate().UpdateByID(&updateData)
	if err != nil {
		return
	}
	//返回
	return
}

// argsReviseNormVal 修订归一化值参数
type argsReviseNormVal struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//年月日
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	YearMD string `db:"year_md" json:"yearMD" index:"true"`
	//归一化值
	ValNorm float64 `db:"val_norm" json:"valNorm" index:"true"`
	//是否为预测值
	IsForecast bool `db:"is_forecast" json:"isForecast"`
}

// reviseNormVal 修订归一化值
func reviseNormVal(args *argsReviseNormVal) (err error) {
	//获取数据
	var rawData FieldsVal
	_ = indexValDB.GetClient().DB.GetPostgresql().Get(&rawData, "SELECT * FROM analysis_index_vals WHERE delete_at < to_timestamp(1000000) AND code = $1 AND year_md = $2 AND extend1 = $3 AND extend2 = $4 AND extend3 = $5 AND extend4 = $6 AND extend5 = $7 AND is_forecast = $8", args.Code, args.YearMD, "", "", "", "", "", args.IsForecast)
	if rawData.ID < 1 {
		err = errors.New(fmt.Sprint("no data, ", err))
		return
	}
	//更新数据
	type updateType struct {
		// ID
		ID int64 `db:"id" json:"id" check:"id" unique:"true"`
		//归一化值
		ValNorm float64 `db:"val_norm" json:"valNorm" index:"true"`
	}
	updateData := updateType{
		ID:      rawData.ID,
		ValNorm: args.ValNorm,
	}
	//修正浮点数
	if math.IsNaN(updateData.ValNorm) || math.IsInf(updateData.ValNorm, 0) {
		updateData.ValNorm = 0
	}
	//修正结果
	updateData.ValNorm = CoreFilter.RoundTo4DecimalPlaces(updateData.ValNorm)
	//执行修改
	err = indexValDB.GetUpdate().UpdateByID(&updateData)
	if err != nil {
		return
	}
	//返回
	return
}

// GetAvgValByCode 获取指定指标平均值
func GetAvgValByCode(code string) (avg float64) {
	_ = indexValDB.GetClient().DB.GetPostgresql().Get(&avg, "SELECT AVG(val_raw) FROM analysis_index_vals WHERE code = $1", code)
	return

}

// GetAvgValNormByCode 获取指定指标平均值
func GetAvgValNormByCode(code string) (avg float64) {
	_ = indexValDB.GetClient().DB.GetPostgresql().Get(&avg, "SELECT AVG(val_norm) FROM analysis_index_vals WHERE code = $1", code)
	return

}
