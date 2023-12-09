package AnalysisAny2

import (
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsSUMBetween 获取指定时间范围的合计数参数
type ArgsSUMBetween struct {
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

// SUMBetween 获取指定时间范围的合计数
func SUMBetween(args *ArgsSUMBetween) (data int64) {
	var err error
	//获取配置
	var configData FieldsConfig
	configData, err = getConfigByMark(args.Mark, true)
	if err != nil {
		return
	}
	//获取缓冲
	cacheMark := fmt.Sprint(getAnyCacheMark(configData.ID), ":sum:between:min.", args.MinAt, ".max.", args.MaxAt, ".", args.OrgID, ".", args.UserID, ".", args.BindID, ".", args.Params1, ".", args.Params2)
	data, err = Router2SystemConfig.MainCache.GetInt64(cacheMark)
	if err == nil && data > 0 {
		return
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT SUM(data) FROM analysis_any2 WHERE config_id = $1 AND ($2 < 0 OR org_id = $2) AND ($3 < 0 OR user_id = $3) AND ($4 < 0 OR bind_id = $4) AND ($5 < 0 OR params1 = $5) AND ($6 < 0 OR params2 = $6) AND create_at >= $7 AND create_at <= $8", configData.ID, args.OrgID, args.UserID, args.BindID, args.Params1, args.Params2, args.MinAt, args.MaxAt)
	if err != nil || data < 1 {
		if args.MinAt.Unix() <= CoreFilter.GetNowTimeCarbon().SubDays(configData.FileDay).Time.Unix() {
			err = Router2SystemConfig.MainDB.Get(&data, "SELECT SUM(data) FROM analysis_any2_file WHERE config_id = $1 AND ($2 < 0 OR org_id = $2) AND ($3 < 0 OR user_id = $3) AND ($4 < 0 OR bind_id = $4) AND ($5 < 0 OR params1 = $5) AND ($6 < 0 OR params2 = $6) AND create_at >= $7 AND create_at <= $8", configData.ID, args.OrgID, args.UserID, args.BindID, args.Params1, args.Params2, args.MinAt, args.MaxAt)
			if err != nil {
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

// ArgsSUMBetweenP 获取两组时间合计数并对比参数
type ArgsSUMBetweenP struct {
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

// SUMBetweenP 获取两组时间合计数并对比
func SUMBetweenP(args *ArgsSUMBetweenP) (prevData int64, lastData int64, addCount int64, p int64) {
	prevData = SUMBetween(&ArgsSUMBetween{
		Mark:    args.Mark,
		OrgID:   args.OrgID,
		UserID:  args.UserID,
		BindID:  args.BindID,
		Params1: args.Params1,
		Params2: args.Params2,
		MinAt:   args.MinAtPrev,
		MaxAt:   args.MaxAtPrev,
	})
	lastData = SUMBetween(&ArgsSUMBetween{
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

// ArgsSUMBetweenArea 获取指定时间范围的数据集参数
type ArgsSUMBetweenArea struct {
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

type DataSUMBetweenArea struct {
	//数据时间
	CreateAt time.Time `json:"createAt"`
	//数据
	Data int64 `json:"data"`
}

// SUMBetweenArea 获取指定时间范围的数据集
// 按照指定颗粒数量分拆均化后展示
func SUMBetweenArea(args *ArgsSUMBetweenArea) (dataList []DataSUMBetweenArea) {
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
		vData := SUMBetween(&ArgsSUMBetween{
			Mark:    args.Mark,
			OrgID:   args.OrgID,
			UserID:  args.UserID,
			BindID:  args.BindID,
			Params1: args.Params1,
			Params2: args.Params2,
			MinAt:   vMinAt,
			MaxAt:   vMaxAt,
		})
		dataList = append(dataList, DataSUMBetweenArea{
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

// ArgsSUMBetweenAreaMonth 获取指定范围的月份数据参数
type ArgsSUMBetweenAreaMonth struct {
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

// DataSUMBetweenAreaMonth 获取指定范围的月份数据集合
type DataSUMBetweenAreaMonth struct {
	//年月
	Date string `json:"date"`
	//数据
	Data int64 `json:"data"`
}

// SUMBetweenAreaMonth 获取指定范围的月份数据
func SUMBetweenAreaMonth(args *ArgsSUMBetweenAreaMonth) (dataList []DataSUMBetweenAreaMonth) {
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
		vData := SUMBetween(&ArgsSUMBetween{
			Mark:    args.Mark,
			OrgID:   args.OrgID,
			UserID:  args.UserID,
			BindID:  args.BindID,
			Params1: args.Params1,
			Params2: args.Params2,
			MinAt:   vMinAt,
			MaxAt:   vMaxAt,
		})
		dataList = append(dataList, DataSUMBetweenAreaMonth{
			Date: nowAt.Format("2006-01"),
			Data: vData,
		})
		//下一步
		nowAt = nowAt.AddMonth()
	}
	//反馈
	return
}
