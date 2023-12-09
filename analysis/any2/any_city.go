package AnalysisAny2

import (
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	MapCityData "gitee.com/weeekj/weeekj_core/v5/map/city_data"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

//获取城市数据集合
/**
1. 采用bindID/param1/param2作为城市的几个重要参数，mode=0时，其中bindID作为城市编码; 其他形式根据未来需求待定
2. mode可以作为参数选择投放模式，注意该设计和添加数据的参数一致
*/

type ArgsGetCity struct {
	//标识码
	Mark string `json:"mark" check:"mark"`
	//组织ID
	// 可留空
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//模式
	// 0 城市模式
	Mode int `json:"mode"`
	//时间范围
	MinAt time.Time `json:"minAt"`
	MaxAt time.Time `json:"maxAt"`
}

type DataGetCity struct {
	//城市名称
	CityName string `json:"cityName"`
	//数据
	Data int64 `json:"data"`
}

func GetCity(args *ArgsGetCity) (dataList []DataGetCity) {
	var err error
	//获取配置
	var configData FieldsConfig
	configData, err = getConfigByMark(args.Mark, true)
	if err != nil {
		return
	}
	//获取缓冲
	cacheMark := fmt.Sprint(getAnyCacheMark(configData.ID), ":city:min.", args.MinAt, ".max.", args.MaxAt, ".", args.OrgID, ".", args.Mode)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	//获取数据
	var rawList []FieldsAny
	_ = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id, org_id, user_id, bind_id, params1, params2, config_id, data FROM analysis_any2 WHERE config_id = $1 AND ($2 < 0 OR org_id = $2) AND create_at >= $7 AND create_at <= $8", configData.ID, args.OrgID, args.MinAt, args.MaxAt)
	//遍历数据
	for _, v := range rawList {
		dataList = append(dataList, DataGetCity{
			CityName: getCityNameBySN(v.BindID),
			Data:     v.Data,
		})
	}
	//写入缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, cacheExpire)
	//反馈
	return
}

// 通过城市编码获取城市名称
func getCityNameBySN(sn int64) string {
	return MapCityData.GetNameByCityCode(fmt.Sprint(sn))
}

// GetCitySNByName 通过城市名称获取编码
// 可用于添加数据时，bindID的指定
func GetCitySNByName(name string) int64 {
	_, sn := MapCityData.GetCodeByCityName(name)
	snInt64, err := CoreFilter.GetInt64ByString(sn)
	if err != nil {
		return 0
	}
	return snInt64
}
