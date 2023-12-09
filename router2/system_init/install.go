package Router2SystemInit

import (
	"errors"
	"fmt"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ToolsInstall "gitee.com/weeekj/weeekj_core/v5/tools/install"
)

func Install() (err error) {
	if !Router2SystemConfig.RunInstall {
		return
	}
	//初始化安装程序
	ToolsInstall.Init(fmt.Sprint(CoreFile.BaseDir()+CoreFile.Sep, "install", CoreFile.Sep, "data", CoreFile.Sep))
	//安装sql
	Router2SystemConfig.MainDB.InstallDir = fmt.Sprint(CoreFile.BaseDir()+CoreFile.Sep, "install", CoreFile.Sep, "sql")
	if err = Router2SystemConfig.MainDB.Install(); err != nil {
		err = errors.New("install sql, " + err.Error())
		return
	}
	//系统配置
	if err = ToolsInstall.InstallConfig(); err != nil {
		err = errors.New("install config, " + err.Error())
		return
	}
	//短信
	if err = ToolsInstall.InstallSMS(); err != nil {
		err = errors.New("install sms, " + err.Error())
		return
	}
	//邮箱
	if err = ToolsInstall.InstallEmail(); err != nil {
		return errors.New("install email, " + err.Error())
	}
	//用户
	if err = ToolsInstall.InstallUser(); err != nil {
		return errors.New("install user, " + err.Error())
	}
	//预警服务
	if err = ToolsInstall.InstallEarlyWarning(); err != nil {
		return errors.New("install early warning, " + err.Error())
	}
	//财务
	if err = ToolsInstall.InstallFinance(); err != nil {
		return errors.New("install finance, " + err.Error())
	}
	//组织
	if err = ToolsInstall.InstallOrg(); err != nil {
		err = errors.New("install org, " + err.Error())
		return
	}
	//提示信息
	fmt.Println("main router system install success.")
	//反馈
	return
}
