package ToolsWeather

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

//FieldsCity 城市和关联数据集合
type FieldsCity struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//城市编码
	CityCode int `db:"city_code" json:"cityCode"`
	//扩展数据
	CityData FieldsCityData `db:"city_data" json:"cityData"`
}

//FieldsCityData 日风天气数据集合
type FieldsCityData struct {
	Name      string `json:"name"`
	Id        string `json:"id"`
	Lat       string `json:"lat"`
	Lon       string `json:"lon"`
	Adm2      string `json:"adm2"`
	Adm1      string `json:"adm1"`
	Country   string `json:"country"`
	Tz        string `json:"tz"`
	UtcOffset string `json:"utcOffset"`
	IsDst     string `json:"isDst"`
	Type      string `json:"type"`
	Rank      string `json:"rank"`
	FxLink    string `json:"fxLink"`
}

//Value sql底层处理器
func (t FieldsCityData) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsCityData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}