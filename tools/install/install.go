package ToolsInstall

import (
	"encoding/json"
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
)

var (
	//配置路径
	configDir = ""
)

func Init(fixDir string) {
	//修正路径位置
	if fixDir != "" {
		configDir = fixDir
	} else {
		configDir = fmt.Sprint(CoreFile.BaseDir(), CoreFile.Sep, "conf", CoreFile.Sep, "install_data", CoreFile.Sep)
	}
}

// 加载配置文件并附加到变量上
func loadConfigFile(name string, data interface{}) error {
	dataByte, err := CoreFile.LoadFile(fmt.Sprint(configDir, name))
	if err != nil {
		return err
	}
	err = json.Unmarshal(dataByte, &data)
	if err != nil {
		return err
	}
	return nil
}

// 检查文件是否存在
func checkConfigFile(name string) bool {
	if CoreFile.IsFile(fmt.Sprint(configDir, name)) {
		return true
	}
	return false
}
