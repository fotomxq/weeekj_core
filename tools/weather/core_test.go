package ToolsWeather

import (
	CoreHttp "gitee.com/weeekj/weeekj_core/v5/core/http"
	MapArea "gitee.com/weeekj/weeekj_core/v5/map/area"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	isInit = false
)

func TestInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	ToolsTest.Init(t)
	if err := MapArea.Init(); err != nil {
		t.Error(err)
	}
}

func TestInitURL(t *testing.T) {
	t.Skip()
	// https://geoapi.qweather.com/v2/city/lookup?key=1cee237f34f8454299ac855b179e3009&location=太原市
	// %E5%A4%AA%E5%8E%9F%E5%B8%82
	getURL := "https://geoapi.qweather.com/v2/city/lookup?key=1cee237f34f8454299ac855b179e3009&location=" + CoreHttp.GetURLEncode("太原市")
	dataByte, err := CoreHttp.GetData(getURL, nil, "", false)
	if err != nil {
		t.Error("tools weather coll run, city, ", err, ", url: ", getURL)
		return
	} else {
		t.Log("get url: ", getURL, ", data byte: ", string(dataByte))
	}
}

func TestRun(t *testing.T) {
	data, err := GetWeather(&ArgsGetWeather{
		Country:  86,
		CityCode: 140100,
		DayCount: 1,
	})
	ToolsTest.ReportData(t, err, data)
	runColl()
}

func TestGetWeather(t *testing.T) {
	data, err := GetWeather(&ArgsGetWeather{
		Country:  86,
		CityCode: 110000,
		DayCount: 1,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetCity(t *testing.T) {
	data, err := GetCity(&ArgsGetCity{
		Country:  86,
		CityCode: 110000,
	})
	ToolsTest.ReportData(t, err, data)
}
