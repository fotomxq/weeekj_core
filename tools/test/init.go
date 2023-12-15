package ToolsTest

import (
	"fmt"
	"testing"

	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//API基础配置结构
// 用于frame、api层级的单元测试处理

var (
	//isInitConfig 是否进行了initConfig
	isInitConfig = false
	//ConfigDirAppend 路径修正
	ConfigDirAppend = "/../../builds/test"
)

// Init 测试初始化入口
func Init(t *testing.T) {
	var err error
	//检查是否已经初始化过
	if isInitConfig {
		return
	} else {
		isInitConfig = true
	}
	//初始化Router2SystemConfig.RootDir
	Router2SystemConfig.RootDir, _ = CoreFile.BaseWDDir()
	Router2SystemConfig.RootDir = fmt.Sprint(Router2SystemConfig.RootDir, ConfigDirAppend)
	fmt.Println("RootDir: ", Router2SystemConfig.RootDir)
	//初始化配置
	//Router2SystemConfig.GlobConfig.PostgresqlNeedConnect = true
	if err = Router2SystemConfig.Init(); err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	Router2SystemConfig.Debug = true
	//初始化日志模块
	CoreLog.Init(Router2SystemConfig.Debug, "weeekj", true)
	//postgres安装包加载
	//CorePostgres.InstallDir = fmt.Sprint(Router2SystemConfig.RootDir, CoreFile.Sep, "install")
}
