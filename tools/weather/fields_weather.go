package ToolsWeather

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

//FieldsWeather 天气数据
type FieldsWeather struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//天气时间
	// 20200101 没有小时等
	DayTime time.Time `db:"day_time" json:"dayTime"`
	//城市ID
	CityID int64 `db:"city_id" json:"cityID"`
	//天气数据集合
	Weather FieldsWeatherData `db:"weather" json:"weather"`
}

//FieldsWeatherData 日风天气气象数据集合
type FieldsWeatherData struct {
	FxDate         string `json:"fxDate"`
	Sunrise        string `json:"sunrise"`
	Sunset         string `json:"sunset"`
	Moonrise       string `json:"moonrise"`
	Moonset        string `json:"moonset"`
	MoonPhase      string `json:"moonPhase"`
	TempMax        string `json:"tempMax"`
	TempMin        string `json:"tempMin"`
	IconDay        string `json:"iconDay"`
	TextDay        string `json:"textDay"`
	IconNight      string `json:"iconNight"`
	TextNight      string `json:"textNight"`
	Wind360Day     string `json:"wind360Day"`
	WindDirDay     string `json:"windDirDay"`
	WindScaleDay   string `json:"windScaleDay"`
	WindSpeedDay   string `json:"windSpeedDay"`
	Wind360Night   string `json:"wind360Night"`
	WindDirNight   string `json:"windDirNight"`
	WindScaleNight string `json:"windScaleNight"`
	WindSpeedNight string `json:"windSpeedNight"`
	Humidity       string `json:"humidity"`
	Precip         string `json:"precip"`
	Pressure       string `json:"pressure"`
	Vis            string `json:"vis"`
	Cloud          string `json:"cloud"`
	UvIndex        string `json:"uvIndex"`
}

//Value sql底层处理器
func (t FieldsWeatherData) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsWeatherData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}