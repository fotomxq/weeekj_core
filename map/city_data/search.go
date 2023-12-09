package MapCityData

// SearchPCR 根据三要素查询符合条件的数据集
func SearchPCR(province, city, region string) (provinceCode string, cityCode string, regionCode string) {
	if region != "" {
		isFind := false
		for _, v := range globCityAreaData {
			for _, v2 := range v.CityList {
				for _, v3 := range v2.AreaList {
					if v3.Name == region {
						provinceCode = v.Code
						cityCode = v2.Code
						regionCode = v3.Code
						isFind = true
						break
					}
				}
				if isFind {
					break
				}
			}
			if isFind {
				break
			}
		}
	} else {
		if city != "" {
			isFind := false
			for _, v := range globCityAreaData {
				for _, v2 := range v.CityList {
					if v2.Name == city {
						provinceCode = v.Code
						cityCode = v2.Code
						isFind = true
						break
					}
				}
				if isFind {
					break
				}
			}
		} else {
			if province != "" {
				for _, v := range globCityAreaData {
					provinceCode = v.Code
					break
				}
			}
		}
	}
	return
}
