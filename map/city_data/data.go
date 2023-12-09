package MapCityData

// 城市数据结构
type dataAreaData struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
type dataCityData struct {
	Code     string         `json:"code"`
	Name     string         `json:"name"`
	AreaList []dataAreaData `json:"areaList"`
}
type DataProvinceData struct {
	Code     string         `json:"code"`
	Name     string         `json:"name"`
	CityList []dataCityData `json:"cityList"`
}
