package Router2SystemInit

import (
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreReg "gitee.com/weeekj/weeekj_core/v5/core/reg"
	RouterGinSet "gitee.com/weeekj/weeekj_core/v5/router/gin_set"
	RouterSystem "gitee.com/weeekj/weeekj_core/v5/router/system"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// Main 初始化程序设计
func Main() {
	//初始化配置
	if err := Router2SystemConfig.Init(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("main router core config init success.")
	//初始化日志模块
	logSaveDB, err := Router2SystemConfig.Cfg.Section("log").Key("log_save_db").Bool()
	if err != nil {
		logSaveDB = true
	}
	CoreLog.Init(Router2SystemConfig.Debug, AppName, logSaveDB)
	fmt.Println("main router core log init success.")
	//装载gin
	RouterGinSet.Init()
	//启动注册机
	if OpenSystemReg {
		//计算注册地址
		CoreReg.Init(AppName + AppVersion)
		//fmt.Println("reg: ", CoreReg.GetKey("245d5c2d6bedd03abe1d", "202301", "202512"))
		//注册机启动
		if b := RouterSystem.Reg(AppName + AppVersion); !b {
			return
		}
		fmt.Println("main router system reg success.")
	}
}
