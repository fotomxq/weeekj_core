package BaseConfig

import (
	"encoding/json"
	"fmt"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
)

// DataGroup 分组配置
type DataGroup struct {
	Group []DataGroupTop `json:"groups"`
}

type DataGroupTop struct {
	Mark  string           `json:"mark"`
	Name  string           `json:"name"`
	Des   string           `json:"des"`
	Child []DataGroupChild `json:"child"`
}

type DataGroupChild struct {
	Mark string `json:"mark"`
	Name string `json:"name"`
	Des  string `json:"des"`
}

// loadGroupData 加载配置文件
func loadGroupData() (err error) {
	var dataByte []byte
	dataByte, err = CoreFile.LoadFile(fmt.Sprint(RootDir, CoreFile.Sep, "data", CoreFile.Sep, "group.json"))
	if err != nil {
		return
	}
	err = json.Unmarshal(dataByte, &groupData)
	if err != nil {
		return
	}
	return
}

// GetGroupData 获取分组数据
func GetGroupData() (data DataGroup) {
	data = groupData
	return
}
