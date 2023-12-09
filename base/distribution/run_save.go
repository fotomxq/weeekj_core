package BaseDistribution

import (
	"encoding/json"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
)

// 将数据定期存储进入数据库中，方便突然中断服务后加载数据
func runSave() error {
	//生成数据集合
	configData := DataType{
		Service:  cacheData,
		Child:    cacheChildData,
		ChildRun: cacheChildRunData,
	}
	//解析数据
	configDataByte, err := json.Marshal(configData)
	if err != nil {
		return err
	}
	//写入数据
	if err := CoreFile.WriteFile(localDataSrc, configDataByte); err != nil {
		return err
	}
	return nil
}
