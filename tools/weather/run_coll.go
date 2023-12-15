package ToolsWeather

import (
	"encoding/json"
	"fmt"
	MapCityData "github.com/fotomxq/weeekj_core/v5/map/city_data"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"

	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreHttp "github.com/fotomxq/weeekj_core/v5/core/http"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	"github.com/golang-module/carbon"
	"strings"
	"time"
)

var (
	hefengIsBusiness = false
	hefengWebID      = ""
	hefengWebKey     = ""
)

func runColl() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("tools weather coll run, ", r)
		}
	}()
	//获取列不同的数据
	var diffList []FieldsCity
	if err := Router2SystemConfig.MainDB.Select(&diffList, "SELECT country, city_code, city_data FROM tools_weather_city"); err != nil {
		//继续
	}
	if len(diffList) < 1 && len(waitWeatherCity) < 1 {
		//等待也为空，则退出
		return
	}
	var err error
	hefengIsBusiness, err = BaseConfig.GetDataBool("ToolsHefengWeatherBusiness")
	if err != nil {
		CoreLog.Error("tools weather coll run, get config by ToolsHefengWeatherBusiness, ", err)
		return
	}
	hefengWebID, err = BaseConfig.GetDataString("ToolsHefengWeatherWebID")
	if err != nil {
		CoreLog.Error("tools weather coll run, get config by ToolsHefengWeatherWebID, ", err)
		return
	}
	hefengWebKey, err = BaseConfig.GetDataString("ToolsHefengWeatherWebKey")
	if err != nil {
		CoreLog.Error("tools weather coll run, get config by ToolsHefengWeatherWebKey, ", err)
		return
	}
	var waitGetCity []dataWaitWeatherCity
	waitWeatherLock.Lock()
	for _, v := range waitWeatherCity {
		isFind := false
		for _, v2 := range diffList {
			if v2.Country == v.Country && v2.CityCode == v.CityCode {
				isFind = true
				break
			}
		}
		if !isFind {
			waitGetCity = append(waitGetCity, dataWaitWeatherCity{
				Country:  v.Country,
				CityCode: v.CityCode,
			})
		}
	}
	waitWeatherCity = []dataWaitWeatherCity{}
	waitWeatherLock.Unlock()
	//获取城市数据
	for _, v := range waitGetCity {
		runCollCity(v.Country, v.CityCode)
		time.Sleep(time.Second * 1)
	}
	if len(waitGetCity) > 0 {
		if err := Router2SystemConfig.MainDB.Select(&diffList, "SELECT country, city_code, city_data FROM tools_weather_city"); err != nil || len(diffList) < 1 {
			return
		}
	}
	//获取天气数据
	for _, v := range diffList {
		//获取城市数据集合
		var weatherCity FieldsCity
		if err := Router2SystemConfig.MainDB.Get(&weatherCity, "SELECT id, city_data FROM tools_weather_city WHERE country = $1 AND city_code = $2", v.Country, v.CityCode); err != nil || weatherCity.ID < 1 {
			continue
		}
		runCollWeather(weatherCity.ID, weatherCity.CityData.Id)
		time.Sleep(time.Second * 1)
	}
}

// 获取指定城市的数据包
// https://geoapi.qweather.com/v2/city/lookup?key=[key]&location=[城市名称]
func runCollCity(country int, cityCode int) {
	//从集合中查询到该城市的名称
	if country != 86 {
		return
	}
	//检查数据是否存在？
	var weatherCity FieldsCity
	if err := Router2SystemConfig.MainDB.Get(&weatherCity, "SELECT id FROM tools_weather_city WHERE country = $1 AND city_code = $2", country, cityCode); err == nil && weatherCity.ID > 0 {
		return
	}
	//获取城市集合，在里面查询城市名称的信息
	cityData := MapCityData.GetCityData()
	cityName := ""
	for _, v := range cityData {
		for _, v2 := range v.CityList {
			if v2.Code == fmt.Sprint(cityCode) {
				cityName = v2.Name
				break
			}
		}
		if cityName != "" {
			break
		}
	}
	if cityName == "" {
		CoreLog.Error("tools weather coll run, city not find, country: ", country, ", cityCode: ", cityCode)
		return
	} else {
		cityName = CoreHttp.GetURLEncode(cityName)
	}
	type referType struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	}
	type dataType struct {
		Code     string           `json:"code"`
		Location []FieldsCityData `json:"location"`
		Refer    referType        `json:"refer"`
	}
	var data dataType
	//请求数据
	getURL := fmt.Sprint("https://geoapi.qweather.com/v2/city/lookup?key=", hefengWebKey, "&location=", cityName)
	dataByte, err := CoreHttp.GetData(getURL, nil, "", false)
	if err != nil {
		CoreLog.Error("tools weather coll run, city, ", err, ", url: ", getURL)
		return
	}
	//解析数据
	if err = json.Unmarshal(dataByte, &data); err != nil {
		CoreLog.Error("tools weather coll run, city, json, ", err)
		return
	}
	//检查结果
	if data.Code != "200" {
		CoreLog.Warn("tools weather coll run, city, not 200, code: ", data.Code, ", url: ", getURL)
		return
	}
	if len(data.Location) < 1 {
		CoreLog.Error("tools weather coll run, city, location data is empty, ", err)
		return
	}
	//记录数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tools_weather_city (country, city_code, city_data) VALUES (:country,:city_code,:city_data)", map[string]interface{}{
		"country":   country,
		"city_code": cityCode,
		"city_data": data.Location[0],
	})
	if err != nil {
		CoreLog.Error("tools weather coll run, city, create data, ", err)
		return
	}
}

// 采集指定城市的最近7天天气预报数据
// https://devapi.qweather.com/v7/weather/3d?key=[key]&location=[城市坐标十进制组合]
func runCollWeather(cityID int64, location string) {
	//获取数据
	type referType struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	}
	type dataType struct {
		Code       string              `json:"code"`
		UpdateTime string              `json:"updateTime"`
		FxLink     string              `json:"fxLink"`
		Daily      []FieldsWeatherData `json:"daily"`
		Refer      referType           `json:"refer"`
	}
	var data dataType
	//请求数据
	getURL := fmt.Sprint("https://devapi.qweather.com/v7/weather/3d?key=", hefengWebKey, "&location=", location)
	dataByte, err := CoreHttp.GetData(getURL, nil, "", false)
	if err != nil {
		CoreLog.Error("tools weather coll run, weather, ", err, ", url: ", getURL)
		return
	}
	//解析数据
	if err = json.Unmarshal(dataByte, &data); err != nil {
		CoreLog.Error("tools weather coll run, weather, json, ", err)
		return
	}
	if data.Code != "200" {
		CoreLog.Error("tools weather coll run, weather, code not 200, code: ", data.Code, ", url: ", getURL)
		return
	}
	//检查是否是否已经存在？
	// 存在相同日期的，则更新该数据
	// 否则创建新的数据
	for _, v := range data.Daily {
		//获取该时间的数据
		var weatherData FieldsWeather
		if err := Router2SystemConfig.MainDB.Get(&weatherData, "SELECT id FROM tools_weather WHERE city_id = $1 AND weather @> '{\"fxDate\":\""+v.FxDate+"\"}'::jsonb", cityID); err == nil && weatherData.ID > 0 {
			//找到数据
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tools_weather SET weather = :weather WHERE id = :id", map[string]interface{}{
				"id":      weatherData.ID,
				"weather": v,
			})
			if err != nil {
				CoreLog.Error("tools weather coll run, weather, update data, ", err)
			}
			continue
		}
		//构建日期
		times := strings.Split(v.FxDate, "-")
		if len(times) != 3 {
			CoreLog.Error("tools weather coll run, weather, fx date: ", v.FxDate, ", len not 3")
			continue
		}
		timeY, err := CoreFilter.GetIntByString(times[0])
		if err != nil {
			CoreLog.Error("tools weather coll run, weather, fx date: ", v.FxDate, ", y not int, ", err)
			continue
		}
		timeM, err := CoreFilter.GetIntByString(times[1])
		if err != nil {
			CoreLog.Error("tools weather coll run, weather, fx date: ", v.FxDate, ", m not int, ", err)
			continue
		}
		timeD, err := CoreFilter.GetIntByString(times[2])
		if err != nil {
			CoreLog.Error("tools weather coll run, weather, fx date: ", v.FxDate, ", d not int, ", err)
			continue
		}
		vTimeAt := carbon.CreateFromDate(timeY, timeM, timeD).StartOfDay()
		//写入数据
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tools_weather (city_id, day_time, weather) VALUES (:city_id,:day_time,:weather)", map[string]interface{}{
			"city_id":  cityID,
			"day_time": vTimeAt.Time,
			"weather":  v,
		})
		if err != nil {
			CoreLog.Error("tools weather coll run, weather, create data, ", err)
		}
	}
}
