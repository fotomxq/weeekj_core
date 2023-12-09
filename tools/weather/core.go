package ToolsWeather

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/robfig/cron"
	"sync"
)

//天气数据采集汇总模块

// 等待拉取天气列
type dataWaitWeatherCity struct {
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//城市编码
	CityCode int `db:"city_code" json:"cityCode"`
}

var (
	//预报拉取的城市气象数据
	waitWeatherCity []dataWaitWeatherCity
	waitWeatherLock sync.Mutex
	//定时器
	runTimer *cron.Cron
	runLock  = false
)

// ArgsGetWeather 获取指定城市的天气预报参数
type ArgsGetWeather struct {
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//城市编码
	CityCode int `db:"city_code" json:"cityCode" check:"intThan0" empty:"true"`
	//查询天数
	// 1 1天 / 3 3天 / 7 7天
	DayCount int `db:"day_count" json:"dayCount" check:"intThan0" empty:"true"`
}

// DataGetWeather 获取指定城市的天气预报数据
type DataGetWeather struct {
	//天气预报结构
	Weathers []FieldsWeatherData `json:"weathers"`
}

// GetWeather 获取指定城市的天气预报
func GetWeather(args *ArgsGetWeather) (data DataGetWeather, err error) {
	type dataType struct {
		//所属国家 国家代码
		// eg: china => 86
		Country int `db:"country" json:"country"`
		//城市编码
		CityCode int `db:"city_code" json:"cityCode"`
		//天气数据集合
		Weather FieldsWeatherData `db:"weather" json:"weather"`
	}
	var dataList []dataType
	if args.DayCount < 0 {
		args.DayCount = 1
	}
	if args.DayCount > 7 {
		args.DayCount = 7
	}
	beforeAt := CoreFilter.GetNowTimeCarbon().StartOfDay()
	afterAt := CoreFilter.GetNowTimeCarbon().EndOfDay().AddDays(args.DayCount)
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT c.country as country, c.city_code as city_code, w.weather as weather FROM tools_weather as w INNER JOIN tools_weather_city as c ON c.id = w.city_id WHERE c.country = $1 AND c.city_code = $2 AND w.day_time >= $3 AND w.day_time <= $4 ORDER BY w.id LIMIT 7", args.Country, args.CityCode, beforeAt.Time, afterAt.Time)
	if err != nil {
		appendWaitWeatherCity(args.Country, args.CityCode)
		return
	}
	for _, v := range dataList {
		data.Weathers = append(data.Weathers, v.Weather)
	}
	if len(data.Weathers) < 1 {
		err = errors.New("data is empty")
		appendWaitWeatherCity(args.Country, args.CityCode)
		return
	}
	return
}

// ArgsGetCityList 获取城市列表参数
type ArgsGetCityList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetCityList 获取城市列表
func GetCityList(args *ArgsGetCityList) (dataList []FieldsCity, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.Search != "" {
		where = where + " AND (city_data -> 'adm1' ? :search OR city_data -> 'adm2' ? :search OR city_data -> 'name' ? :search)"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"tools_weather_city",
		"id",
		"SELECT id, country, city_code, city_data FROM tools_weather_city WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "city_code"},
	)
	return
}

// ArgsGetCity 获取指定城市的数据包参数
type ArgsGetCity struct {
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//城市编码
	CityCode int `db:"city_code" json:"cityCode" check:"intThan0" empty:"true"`
}

// GetCity 获取指定城市的数据包参数
func GetCity(args *ArgsGetCity) (data FieldsCity, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, country, city_code, city_data FROM tools_weather_city WHERE country = $1 AND city_code = $2", args.Country, args.CityCode)
	return
}

// 添加等待拉取天气的城市列
func appendWaitWeatherCity(country int, cityCode int) {
	waitWeatherLock.Lock()
	isFind := false
	for _, v := range waitWeatherCity {
		if v.Country == country && v.CityCode == cityCode {
			isFind = true
			break
		}
	}
	if !isFind {
		waitWeatherCity = append(waitWeatherCity, dataWaitWeatherCity{
			Country:  country,
			CityCode: cityCode,
		})
	}
	waitWeatherLock.Unlock()
}
