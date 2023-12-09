package CoreSQL2

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ArgsAddress struct {
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country" empty:"true"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province"`
	//所属城市
	City int `db:"city" json:"city" check:"city" empty:"true"`
	//街道详细信息
	Address string `db:"address" json:"address"`
	//地图制式
	// 0 WGS-84 / 1 GCJ-02 / 2 BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//联系人姓名
	Name string `db:"name" json:"name"`
	//联系人国家代码
	NationCode string `db:"nation_code" json:"nationCode"`
	//联系人手机号
	Phone string `db:"phone" json:"phone"`
}

func (t *ArgsAddress) GetField() FieldsAddress {
	return FieldsAddress{
		Country:    t.Country,
		Province:   t.Province,
		City:       t.City,
		Address:    t.Address,
		MapType:    t.MapType,
		Longitude:  t.Longitude,
		Latitude:   t.Latitude,
		Name:       t.Name,
		NationCode: t.NationCode,
		Phone:      t.Phone,
	}
}

// FieldsAddress 通用地址结构
type FieldsAddress struct {
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province"`
	//所属城市
	City int `db:"city" json:"city" check:"city"`
	//街道详细信息
	Address string `db:"address" json:"address"`
	//地图制式
	// 0 WGS-84 / 1 GCJ-02 / 2 BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//联系人姓名
	Name string `db:"name" json:"name"`
	//联系人国家代码
	NationCode string `db:"nation_code" json:"nationCode"`
	//联系人手机号
	Phone string `db:"phone" json:"phone"`
}

// Value sql底层处理器
func (t FieldsAddress) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsAddress) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

func (t *FieldsAddress) Check() (err error) {
	switch t.MapType {
	case 0:
		//"WGS-84"
	case 1:
		//"GCJ-02"
	case 2:
		//"BD-09"
	default:
		err = errors.New("gps map type error")
	}
	return
}

func (t *FieldsAddress) GetMapType() (data string) {
	switch t.MapType {
	case 0:
		data = "WGS-84"
	case 1:
		data = "GCJ-02"
	case 2:
		data = "BD-09"
	default:
		return
	}
	return
}
