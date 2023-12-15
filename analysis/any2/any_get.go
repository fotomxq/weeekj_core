package AnalysisAny2

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetLast 获取指定最新的数据参数
type ArgsGetLast struct {
	//标识码
	Mark string `json:"mark" check:"mark"`
	//组织ID
	// 可留空
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可留空
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//绑定ID
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//扩展参数1
	// 可选扩展参数，默认给0无视
	Params1 int64 `json:"params1" check:"id" empty:"true"`
	//扩展参数2
	// 可选扩展参数，默认给0无视
	Params2 int64 `json:"params2" check:"id" empty:"true"`
}

// GetLast 获取指定时间最新的数据
func GetLast(args *ArgsGetLast) (data int64) {
	var err error
	//获取配置
	var configData FieldsConfig
	configData, err = getConfigByMark(args.Mark, true)
	if err != nil {
		return
	}
	//获取缓冲
	cacheMark := fmt.Sprint(getAnyCacheMark(configData.ID), ":last:", args.OrgID, ".", args.UserID, ".", args.BindID, ".", args.Params1, ".", args.Params2)
	data, err = Router2SystemConfig.MainCache.GetInt64(cacheMark)
	if err == nil && data > 0 {
		//return
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT data FROM analysis_any2 WHERE config_id = $1 AND org_id = $2 AND user_id = $3 AND bind_id = $4 AND params1 = $5 AND params2 = $6 ORDER BY id DESC LIMIT 1", configData.ID, args.OrgID, args.UserID, args.BindID, args.Params1, args.Params2)
	if err != nil || data < 1 {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT data FROM analysis_any2_file WHERE config_id = $1 AND org_id = $2 AND user_id = $3 AND bind_id = $4 AND params1 = $5 AND params2 = $6 ORDER BY id DESC LIMIT 1", configData.ID, args.OrgID, args.UserID, args.BindID, args.Params1, args.Params2)
		if err != nil {
			return
		}
	}
	//写入缓冲
	Router2SystemConfig.MainCache.SetInt64(cacheMark, data, cacheExpire)
	//反馈
	return
}

// ArgsGetBetween 获取指定时间范围的数据参数
type ArgsGetBetween struct {
	//标识码
	Mark string `json:"mark" check:"mark"`
	//组织ID
	// 可留空
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可留空
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//绑定ID
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//扩展参数1
	// 可选扩展参数，默认给0无视
	Params1 int64 `json:"params1" check:"id" empty:"true"`
	//扩展参数2
	// 可选扩展参数，默认给0无视
	Params2 int64 `json:"params2" check:"id" empty:"true"`
	//时间范围
	MinAt time.Time `json:"minAt"`
	MaxAt time.Time `json:"maxAt"`
}

// GetBetween 获取指定时间范围的数据
func GetBetween(args *ArgsGetBetween) (data int64) {
	//获取配置
	var configData FieldsConfig
	configData, _ = getConfigByMark(args.Mark, true)
	if configData.ID < 1 {
		return
	}
	//获取缓冲
	cacheMark := fmt.Sprint(getAnyCacheMark(configData.ID), ":between:min.", args.MinAt, ".max.", args.MaxAt, ".", args.OrgID, ".", args.UserID, ".", args.BindID, ".", args.Params1, ".", args.Params2)
	data, _ = Router2SystemConfig.MainCache.GetInt64(cacheMark)
	if data > 0 {
		return
	}
	//获取数据
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT data FROM analysis_any2 WHERE config_id = $1 AND org_id = $2 AND user_id = $3 AND bind_id = $4 AND params1 = $5 AND params2 = $6 AND create_at >= $7 AND create_at <= $8 ORDER BY id DESC LIMIT 1", configData.ID, args.OrgID, args.UserID, args.BindID, args.Params1, args.Params2, args.MinAt, args.MaxAt)
	if data < 1 {
		if args.MinAt.Unix() <= CoreFilter.GetNowTimeCarbon().SubDays(configData.FileDay).Time.Unix() {
			_ = Router2SystemConfig.MainDB.Get(&data, "SELECT data FROM analysis_any2_file WHERE config_id = $1 AND org_id = $2 AND user_id = $3 AND bind_id = $4 AND params1 = $5 AND params2 = $6 AND create_at >= $7 AND create_at <= $8 ORDER BY id DESC LIMIT 1", configData.ID, args.OrgID, args.UserID, args.BindID, args.Params1, args.Params2, args.MinAt, args.MaxAt)
			if data < 1 {
				return
			}
		} else {
			return
		}
	}
	//写入缓冲
	Router2SystemConfig.MainCache.SetInt64(cacheMark, data, cacheExpire)
	//反馈
	return
}

// ArgsGetBetweenP 获取两组时间数据并对比参数
type ArgsGetBetweenP struct {
	//标识码
	Mark string `json:"mark" check:"mark"`
	//组织ID
	// 可留空
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可留空
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//绑定ID
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//扩展参数1
	// 可选扩展参数，默认给0无视
	Params1 int64 `json:"params1" check:"id" empty:"true"`
	//扩展参数2
	// 可选扩展参数，默认给0无视
	Params2 int64 `json:"params2" check:"id" empty:"true"`
	//时间范围
	MinAtPrev time.Time `json:"minAtPrev"`
	MaxAtPrev time.Time `json:"maxAtPrev"`
	MinAtLast time.Time `json:"minAtLast"`
	MaxAtLast time.Time `json:"maxAtLast"`
}

// GetBetweenP 获取两组时间数据并对比
func GetBetweenP(args *ArgsGetBetweenP) (prevData int64, lastData int64, addCount int64, p int64) {
	prevData = GetBetween(&ArgsGetBetween{
		Mark:    args.Mark,
		OrgID:   args.OrgID,
		UserID:  args.UserID,
		BindID:  args.BindID,
		Params1: args.Params1,
		Params2: args.Params2,
		MinAt:   args.MinAtPrev,
		MaxAt:   args.MaxAtPrev,
	})
	lastData = GetBetween(&ArgsGetBetween{
		Mark:    args.Mark,
		OrgID:   args.OrgID,
		UserID:  args.UserID,
		BindID:  args.BindID,
		Params1: args.Params1,
		Params2: args.Params2,
		MinAt:   args.MinAtLast,
		MaxAt:   args.MaxAtLast,
	})
	var p2 float64
	addCount, p2 = CoreFilter.MathLastProportion(prevData, lastData)
	if p2 != 0 {
		p = CoreFilter.GetInt64ByFloat64(p2 * 10000)
	}
	return
}

// ArgsGetBetweenArea 获取指定时间范围的数据集参数
type ArgsGetBetweenArea struct {
	//标识码
	Mark string `json:"mark" check:"mark"`
	//组织ID
	// 可留空
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可留空
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//绑定ID
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//扩展参数1
	// 可选扩展参数，默认给0无视
	Params1 int64 `json:"params1" check:"id" empty:"true"`
	//扩展参数2
	// 可选扩展参数，默认给0无视
	Params2 int64 `json:"params2" check:"id" empty:"true"`
	//时间范围
	MinAt time.Time `json:"minAt"`
	MaxAt time.Time `json:"maxAt"`
	//均化数量限制
	// 将数据以平均值形式散开投放到集合内
	// 此限制为切割长度
	Limit int64 `json:"limit"`
}

type DataGetBetweenArea struct {
	//数据时间
	CreateAt time.Time `json:"createAt"`
	//数据
	Data int64 `json:"data"`
}

// GetBetweenArea 获取指定时间范围的数据集
// 按照指定颗粒数量分拆均化后展示
func GetBetweenArea(args *ArgsGetBetweenArea) (dataList []DataGetBetweenArea) {
	//修正限制
	if args.Limit < 1 {
		args.Limit = 1
	}
	if args.Limit > 3000 {
		args.Limit = 3000
	}
	//切割时间范围
	// 按照秒来切割
	areaSec := args.MaxAt.Unix() - args.MinAt.Unix()
	if areaSec < 1 {
		areaSec = 1
	}
	areaSrcPer := CoreFilter.GetInt64ByFloat64(float64(areaSec) / float64(args.Limit))
	// 如果间隔少于1，则按照1计算
	if areaSrcPer < 1 {
		areaSrcPer = 1
	}
	//分开获取数据包
	var step int64 = 1
	stepMin := args.MinAt.Unix()
	stepMax := stepMin + areaSrcPer
	if stepMax < stepMin {
		return
	}
	for {
		//超出范围后退出
		if stepMin > args.MaxAt.Unix() {
			break
		}
		if step > args.Limit {
			break
		}
		//获取指定时间范围数据
		vMinAt := CoreFilter.GetNowTimeCarbon().CreateFromTimestamp(stepMin).Time
		vMaxAt := CoreFilter.GetNowTimeCarbon().CreateFromTimestamp(stepMax).Time
		vData := GetBetween(&ArgsGetBetween{
			Mark:    args.Mark,
			OrgID:   args.OrgID,
			UserID:  args.UserID,
			BindID:  args.BindID,
			Params1: args.Params1,
			Params2: args.Params2,
			MinAt:   vMinAt,
			MaxAt:   vMaxAt,
		})
		dataList = append(dataList, DataGetBetweenArea{
			CreateAt: vMaxAt,
			Data:     vData,
		})
		//进一步
		stepMin += areaSrcPer
		stepMax = stepMin + areaSrcPer
		step += 1
	}
	//反馈
	return
}

// ArgsGetBetweenAreaMonth 获取指定范围的月份数据参数
type ArgsGetBetweenAreaMonth struct {
	//标识码
	Mark string `json:"mark" check:"mark"`
	//组织ID
	// 可留空
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可留空
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//绑定ID
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//扩展参数1
	// 可选扩展参数，默认给0无视
	Params1 int64 `json:"params1" check:"id" empty:"true"`
	//扩展参数2
	// 可选扩展参数，默认给0无视
	Params2 int64 `json:"params2" check:"id" empty:"true"`
	//时间范围
	MinAt time.Time `json:"minAt"`
	MaxAt time.Time `json:"maxAt"`
}

// DataGetBetweenAreaMonth 获取指定范围的月份数据集合
type DataGetBetweenAreaMonth struct {
	//年月
	Date string `json:"date"`
	//数据
	Data int64 `json:"data"`
}

// GetBetweenAreaMonth 获取指定范围的月份数据
func GetBetweenAreaMonth(args *ArgsGetBetweenAreaMonth) (dataList []DataGetBetweenAreaMonth) {
	//不能超出10年
	if args.MaxAt.Unix()-args.MinAt.Unix() > 315360000 {
		return
	}
	//以时间范围遍历数据
	nowAt := CoreFilter.GetNowTimeCarbon().CreateFromTimestamp(args.MinAt.Unix())
	for {
		//超出范围后退出
		if nowAt.Time.Unix() >= args.MaxAt.Unix() {
			break
		}
		//获取指定时间范围数据
		vMinAt := nowAt.StartOfMonth().Time
		vMaxAt := nowAt.EndOfMonth().Time
		vData := GetBetween(&ArgsGetBetween{
			Mark:    args.Mark,
			OrgID:   args.OrgID,
			UserID:  args.UserID,
			BindID:  args.BindID,
			Params1: args.Params1,
			Params2: args.Params2,
			MinAt:   vMinAt,
			MaxAt:   vMaxAt,
		})
		dataList = append(dataList, DataGetBetweenAreaMonth{
			Date: nowAt.Format("2006-01"),
			Data: vData,
		})
		//下一步
		nowAt = nowAt.AddMonth()
	}
	//反馈
	return
}
