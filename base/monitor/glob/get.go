package BaseMonitorGlob

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// GetAll 获取所有数据信息
func GetAll() (dataList []DataGlob) {
	//搜索数据集合
	keys, err := Router2SystemConfig.MainCache.FindKeys(fmt.Sprint(cacheDataKey, "*"))
	if err != nil {
		return
	}
	for _, v := range keys {
		var vData DataGlob
		err = Router2SystemConfig.MainCache.GetStruct(v, &vData)
		if err != nil {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}
