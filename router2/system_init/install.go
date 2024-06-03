package Router2SystemInit

import (
	"errors"
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ToolsInstall "github.com/fotomxq/weeekj_core/v5/tools/install"
)

var (
	//OpenInstallLocalSQL 安装本地化SQL
	OpenInstallLocalSQL = true
	// OpenInstallConfig 启动开关
	OpenInstallConfig = true
	// OpenInstallSMS 启动SMS
	OpenInstallSMS = true
	// OpenInstallEmail 启动邮箱
	OpenInstallEmail = true
	// OpenInstallUser 启动用户
	OpenInstallUser = true
	// OpenInstallEarlyWarning 启动预警模块
	OpenInstallEarlyWarning = true
	// OpenInstallFinance 启动财务
	OpenInstallFinance = true
	// OpenInstallOrg 启动组织
	OpenInstallOrg = true
)

func Install() (err error) {
	if !Router2SystemConfig.RunInstall {
		return
	}
	//记录开始时间
	startRun := RunStartPrint{
		Pre: "main router system install success",
		Suf: "",
	}
	startRun.Start()
	//初始化安装程序
	ToolsInstall.Init(fmt.Sprint(CoreFile.BaseDir()+CoreFile.Sep, "install", CoreFile.Sep, "data", CoreFile.Sep))
	//安装sql
	Router2SystemConfig.MainDB.InstallDir = fmt.Sprint(CoreFile.BaseDir()+CoreFile.Sep, "install", CoreFile.Sep, "sql")
	if OpenInstallLocalSQL {
		if err = Router2SystemConfig.MainDB.Install(); err != nil {
			err = errors.New("install sql, " + err.Error())
			return
		}
	}
	//系统配置
	if OpenInstallConfig {
		if err = ToolsInstall.InstallConfig(); err != nil {
			err = errors.New("install config, " + err.Error())
			return
		}
	}
	//短信
	if OpenInstallSMS {
		if err = ToolsInstall.InstallSMS(); err != nil {
			err = errors.New("install sms, " + err.Error())
			return
		}
	}
	//邮箱
	if OpenInstallEmail {
		if err = ToolsInstall.InstallEmail(); err != nil {
			return errors.New("install email, " + err.Error())
		}
	}
	//用户
	if OpenInstallUser {
		if err = ToolsInstall.InstallUser(); err != nil {
			return errors.New("install user, " + err.Error())
		}
	}
	//预警服务
	if OpenInstallEarlyWarning {
		if err = ToolsInstall.InstallEarlyWarning(); err != nil {
			return errors.New("install early warning, " + err.Error())
		}
	}
	//财务
	if OpenInstallFinance {
		if err = ToolsInstall.InstallFinance(); err != nil {
			return errors.New("install finance, " + err.Error())
		}
	}
	//组织
	if OpenInstallOrg {
		if err = ToolsInstall.InstallOrg(); err != nil {
			err = errors.New("install org, " + err.Error())
			return
		}
	}
	//启动完成提示
	startRun.EndPrint()
	//反馈
	return
}
