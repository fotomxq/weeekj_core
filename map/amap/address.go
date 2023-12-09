package MapAMap

import (
	"fmt"
)

//ArgsGetGPSByAddress 逆地址解析处理参数
type ArgsGetGPSByAddress struct {
	//地址
	Address string `json:"address"`
	//城市
	// 行政编码
	City int `json:"city"`
}

//DataGetGPSByAddressGeocodes 子结构
type DataGetGPSByAddressGeocodes struct {
	//结构化地址信息
	FormattedAddress interface{} `json:"formatted_address"`
	//国家
	Country string `json:"country"`
	//地址所在的省份名
	Province string `json:"province"`
	//城市编码
	City string `json:"city"`
	//地址所在的城市名
	CityCode string `json:"citycode"`
	//地址所在的区
	District string `json:"district"`
	//街道
	// 空数组或string结构
	Street interface{} `json:"street"`
	//门牌
	// 空数组或string结构
	Number interface{} `json:"number"`
	//区域编码
	AdCode string `json:"adcode"`
	//坐标点
	Location string `json:"location"`
	//匹配级别
	Level string `json:"level"`
}

//DataGetGPSByAddress 逆地址解析处理反馈数据
type DataGetGPSByAddress struct {
	Status string `json:"status"`
	Info string `json:"info"`
	InfoCore string `json:"infocode"`
	Count string `json:"count"`
	Geocodes []DataGetGPSByAddressGeocodes `json:"geocodes"`
}

//GetGPSByAddress 逆地址解析处理反馈
func GetGPSByAddress(args *ArgsGetGPSByAddress) (data DataGetGPSByAddress, err error) {
	//请求加载数据
	params := map[string]string{
		"address": args.Address,
	}
	if args.City > 0 {
		params["city"] = fmt.Sprint(args.City)
	}
	err = httpGet("geocode/geo?", params, &data)
	if err != nil{
		return
	}
	//反馈数据
	return
}

//ArgsGetAddressByGPS 通过经纬度获取数据信息参数
type ArgsGetAddressByGPS struct {
	//坐标位置
	// 必须采用高德地图系坐标
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

type DataGetAddressByGPS struct {
	Status string `json:"status"`
	Info string `json:"info"`
	InfoCore string `json:"infocode"`
	Regeocode struct {
		FormattedAddress interface{} `json:"formatted_address"`
		AddressComponent struct {
			Country      string        `json:"country"`
			Province     string        `json:"province"`
			City         interface{}   `json:"city"`
			Citycode     string        `json:"citycode"`
			District     string        `json:"district"`
			Adcode       string        `json:"adcode"`
			Township     string        `json:"township"`
			Towncode     string        `json:"towncode"`
			Neighborhood struct {
				Name interface{} `json:"name"`
				Type interface{} `json:"type"`
			} `json:"neighborhood"`
			Building struct {
				Name interface{} `json:"name"`
				Type interface{} `json:"type"`
			} `json:"building"`
			StreetNumber struct {
				Street    interface{} `json:"street"`
				Number    interface{} `json:"number"`
				Location  interface{} `json:"location"`
				Direction interface{} `json:"direction"`
				Distance  interface{} `json:"distance"`
			} `json:"streetNumber"`
			BusinessAreas interface{} `json:"businessAreas"`
		} `json:"addressComponent"`
		Pois []struct {
			Id           string `json:"id"`
			Name         string `json:"name"`
			Type         string `json:"type"`
			Tel          string `json:"tel"`
			Direction    string `json:"direction"`
			Distance     string `json:"distance"`
			Location     string `json:"location"`
			Address      string `json:"address"`
			Poiweight    string `json:"poiweight"`
			Businessarea string `json:"businessarea"`
		} `json:"pois"`
		Roads []struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			Direction string `json:"direction"`
			Distance  string `json:"distance"`
			Location  string `json:"location"`
		} `json:"roads"`
		Roadinters []struct {
			Direction  string `json:"direction"`
			Distance   string `json:"distance"`
			Location   string `json:"location"`
			FirstId    string `json:"first_id"`
			FirstName  string `json:"first_name"`
			SecondId   string `json:"second_id"`
			SecondName string `json:"second_name"`
		} `json:"roadinters"`
		Aois []struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			Adcode   string `json:"adcode"`
			Location string `json:"location"`
			Area     string `json:"area"`
			Distance string `json:"distance"`
			Type     string `json:"type"`
		} `json:"aois"`
	} `json:"regeocode"`
}

//GetAddressByGPS 通过经纬度获取数据信息
func GetAddressByGPS(args *ArgsGetAddressByGPS) (data DataGetAddressByGPS, err error) {
	//请求加载数据
	params := map[string]string{
		"location": fmt.Sprint(args.Longitude, ",", args.Latitude),
	}
	err = httpGet("geocode/regeo?", params, &data)
	if err != nil{
		return
	}
	//反馈数据
	return
}