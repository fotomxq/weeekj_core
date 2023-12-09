package MapCityData

// GetCityData 获取全局城市数据包
func GetCityData() []DataProvinceData {
	return globCityAreaData
}

func GetOnlyCityData() []dataAreaData {
	var result []dataAreaData
	for _, v := range globCityAreaData {
		for _, v2 := range v.CityList {
			result = append(result, dataAreaData{
				Code: v2.Code,
				Name: v2.Name,
			})
		}
	}
	return result
}

// 正向查询所属分区
// 如果是城市区级别，则反馈城市分区
// 如果完全找不到，则原样反馈
func getCityData(code string) (cityCode string) {
	for _, vProvince := range globCityAreaData {
		if vProvince.Code == code {
			cityCode = vProvince.Code
			return
		}
		for _, vCity := range vProvince.CityList {
			if vCity.Code == code {
				cityCode = vCity.Code
				return
			}
			for _, vArea := range vCity.AreaList {
				if vArea.Code == code {
					cityCode = vCity.Code
					return
				}
			}
		}
	}
	cityCode = code
	return
}

// GetCodeByCityName 通过城市检索行政编码
func GetCodeByCityName(cityName string) (provinceCode string, cityCode string) {
	for _, vProvince := range globCityAreaData {
		for _, vCity := range vProvince.CityList {
			if vCity.Name == cityName {
				provinceCode = vProvince.Code
				cityCode = vCity.Code
				return
			}
		}
	}
	return
}

// GetNameByCityCode 通过城市编码获取名称
func GetNameByCityCode(code string) string {
	for _, vProvince := range globCityAreaData {
		for _, vCity := range vProvince.CityList {
			if vCity.Code == code {
				return vCity.Name
			}
		}
	}
	return ""
}
