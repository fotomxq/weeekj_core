package MapCityData

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 从文件加载数据到内存
func loadCityData() (err error) {
	var dataByte []byte
	dataByte, err = CoreFile.LoadFile(fmt.Sprint(Router2SystemConfig.RootDir, CoreFile.Sep, "data", CoreFile.Sep, "city_data.json"))
	if err != nil {
		err = errors.New("load city_data json file, " + err.Error())
		return
	}
	err = json.Unmarshal(dataByte, &globCityAreaData)
	if err != nil {
		return
	}
	return
}
