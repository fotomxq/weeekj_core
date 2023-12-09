package AnalysisAny2

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

//排名模块设计
/**
1. 采用param1/param2等参数作为参照物，将数据打散到同一个时间点的多个记录内
2. 获取排名直接使用本方法即可，注意不能给与其他特征参数
3. 可以给与其他参数，但注意不想给与的参数标记-1，系统获取所有数据并做排序反馈
*/

type ArgsGetSort struct {
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
	//是否倒叙
	Desc bool `json:"desc"`
	//均化数量限制
	// 将数据以平均值形式散开投放到集合内
	// 此限制为切割长度
	Limit int `json:"limit"`
}

type DataGetSort struct {
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
	//数据
	Data int64 `db:"data" json:"data"`
}

func GetSort(args *ArgsGetSort) (dataList []DataGetSort) {
	var err error
	//限制参数
	if args.Limit < 1 {
		args.Limit = 1
	}
	if args.Limit > 1000 {
		args.Limit = 1000
	}
	//获取配置
	var configData FieldsConfig
	configData, err = getConfigByMark(args.Mark, true)
	if err != nil {
		return
	}
	//获取缓冲
	cacheMark := fmt.Sprint(getAnyCacheMark(configData.ID), ":sort:min.", args.MinAt, ".max.", args.MaxAt, ".", args.OrgID, ".", args.UserID, ".", args.BindID, ".", args.Params1, ".", args.Params2, ".", args.Desc, ".", args.Limit)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	//获取数据集
	var rawList []FieldsAny
	if args.Desc {
		_ = Router2SystemConfig.MainDB.Select(&rawList, "SELECT org_id, user_id, bind_id, params1, params2, data FROM (SELECT org_id, user_id, bind_id, params1, params2, SUM(data) as data FROM analysis_any2 WHERE config_id = $1 AND ($2 < 0 OR org_id = $2) AND ($3 < 0 OR user_id = $3) AND ($4 < 0 OR bind_id = $4) AND ($5 < 0 OR params1 = $5) AND ($6 < 0 OR params2 = $6) AND create_at >= $7 AND create_at <= $8 GROUP BY org_id, user_id, bind_id, params1, params2, data) as analysis_any2 ORDER BY data DESC LIMIT $9", configData.ID, args.OrgID, args.UserID, args.BindID, args.Params1, args.Params2, args.MinAt, args.MaxAt, args.Limit)
	} else {
		_ = Router2SystemConfig.MainDB.Select(&rawList, "SELECT org_id, user_id, bind_id, params1, params2, data FROM (SELECT org_id, user_id, bind_id, params1, params2, SUM(data) as data FROM analysis_any2 WHERE config_id = $1 AND ($2 < 0 OR org_id = $2) AND ($3 < 0 OR user_id = $3) AND ($4 < 0 OR bind_id = $4) AND ($5 < 0 OR params1 = $5) AND ($6 < 0 OR params2 = $6) AND create_at >= $7 AND create_at <= $8 GROUP BY org_id, user_id, bind_id, params1, params2, data) as analysis_any2 ORDER BY data LIMIT $9", configData.ID, args.OrgID, args.UserID, args.BindID, args.Params1, args.Params2, args.MinAt, args.MaxAt, args.Limit)
	}
	//重构数据集
	for kRaw := 0; kRaw < len(rawList); kRaw++ {
		vRaw := rawList[kRaw]
		dataList = append(dataList, DataGetSort{
			OrgID:   vRaw.OrgID,
			UserID:  vRaw.UserID,
			BindID:  vRaw.BindID,
			Params1: vRaw.Param1,
			Params2: vRaw.Param2,
			Data:    vRaw.Data,
		})
	}
	//保存缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, cacheExpire)
	//反馈
	return
}

// DataGetSortP 获取排名数据的占比排列数据
type DataGetSortP struct {
	//总额度
	AllDataCount int64 `json:"allDataCount"`
	//数据集合
	DataList []DataGetSortPData
}

type DataGetSortPData struct {
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
	//数据
	Data int64 `db:"data" json:"data"`
	//占比
	// 35.12% = 3512
	DataP int64 `json:"dataP"`
}

// GetSortP 获取排名数据的占比排列
// 将排名全部相加后，计算每个排名对象的占比
func GetSortP(args *ArgsGetSort) (data DataGetSortP) {
	//获取排名
	sortList := GetSort(args)
	if len(sortList) < 1 {
		data.DataList = []DataGetSortPData{}
		return
	}
	//计算排名总数量
	data.AllDataCount = SUMBetween(&ArgsSUMBetween{
		Mark:    args.Mark,
		OrgID:   args.OrgID,
		UserID:  args.UserID,
		BindID:  -1,
		Params1: -1,
		Params2: -1,
		MinAt:   args.MinAt,
		MaxAt:   args.MaxAt,
	})
	//二次遍历计算占比
	// 最后一个占比减去计算
	for _, v := range sortList {
		var dataP float64 = 0
		if data.AllDataCount > 0 {
			dataP = float64(v.Data) / float64(data.AllDataCount)
		}
		data.DataList = append(data.DataList, DataGetSortPData{
			OrgID:   v.OrgID,
			UserID:  v.UserID,
			BindID:  v.BindID,
			Params1: v.Params1,
			Params2: v.Params2,
			Data:    v.Data,
			DataP:   int64(dataP * 10000),
		})
	}
	//反馈
	return
}

// ArgsGetSortNum 获取指定参数在排名的第几位参数
type ArgsGetSortNum struct {
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

// GetSortNum 获取指定参数在排名的第几位
func GetSortNum(args *ArgsGetSortNum) (id int64) {
	//获取配置
	configData, err := getConfigByMark(args.Mark, true)
	if err != nil {
		return
	}
	//获取该范围的最后数据量
	lastData := SUMBetween(&ArgsSUMBetween{
		Mark:    args.Mark,
		OrgID:   args.OrgID,
		UserID:  args.UserID,
		BindID:  args.BindID,
		Params1: args.Params1,
		Params2: args.Params1,
		MinAt:   args.MinAt,
		MaxAt:   args.MaxAt,
	})
	//根据数据量检查之前的所有数据综合
	_ = Router2SystemConfig.MainDB.Get(&id, "SELECT COUNT(id) FROM analysis_any2 WHERE config_id = $1 AND ($2 < 0 OR org_id = $2) AND ($3 < 0 OR user_id = $3) AND ($4 < 0 OR bind_id = $4) AND ($5 < 0 OR params1 = $5) AND ($6 < 0 OR params2 = $6) AND create_at >= $7 AND create_at <= $8 AND data > $9 GROUP BY org_id, user_id, bind_id, params1, params2", configData.ID, args.OrgID, args.UserID, args.BindID, args.Params1, args.Params2, args.MinAt, args.MaxAt, lastData)
	//修正数据
	id += 1
	//反馈
	return
}
