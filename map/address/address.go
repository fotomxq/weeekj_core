package MapAddress

import (
	MapCityData "gitee.com/weeekj/weeekj_core/v5/map/city_data"
	"github.com/pupuk/addr"
)

// DataGetAddressByStr 分析地址字符串为一个组合数据集合
type DataGetAddressByStr struct {
	//所属国家 国家代码
	Country int `json:"country"`
	//省份
	Province string `json:"province"`
	//省份编码
	ProvinceCode string `json:"provinceCode"`
	//所属城市
	City string `json:"city"`
	//城市编码
	CityCode string `json:"cityCode"`
	//地区
	Region string `json:"region"`
	//地区编码
	RegionCode string `json:"regionCode"`
	//街道
	Street string `json:"street"`
	//街道详细信息
	Address string `json:"address"`
	//邮编
	PostCode string `json:"postCode"`
	//姓名
	Name string `json:"name"`
	//联系电话
	Phone string `json:"phone"`
	//身份证
	IDCard string `json:"idCard"`
}

// GetAddressByStr 分析地址字符串为一个组合
func GetAddressByStr(str string) (data DataGetAddressByStr) {
	//如果为空则反馈
	if str == "" {
		return
	}
	//解析数据
	parse := addr.Smart(str)
	//组合数据
	data = DataGetAddressByStr{
		Country:  86,
		Province: parse.Province,
		City:     parse.City,
		Region:   parse.Region,
		Street:   parse.Street,
		Address:  parse.Address,
		PostCode: parse.PostCode,
		Name:     parse.Name,
		Phone:    parse.Mobile,
		IDCard:   parse.IdNumber,
	}
	//附加城市编码等信息
	data.ProvinceCode, data.CityCode, data.RegionCode = MapCityData.SearchPCR(data.Province, data.City, data.Region)
	//反馈数据
	return
}
