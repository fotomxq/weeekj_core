package CoreLog

import (
	"fmt"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	"github.com/robfig/cron"
)

var (
	//定时器
	runTimer    *cron.Cron
	runMakeLock = false
	runZipLock  = false

	//是否启动转存数据库模式
	openToDB bool
	//日志存储目录
	logDir string
	//debug模式
	debugOn bool
	//gin日志对象
	ginLog logConfig
	//全局普通日志
	globLog logConfig
	//错误日志
	errLog logConfig
	//警告日志
	warnLog logConfig
	//mqtt
	mqttLog logConfig
	//客户端采集日志
	appLog logConfig
	//是否启动gin日志
	allowGin = true
	//是否启动普通日志
	allowDefault = true

	//归档目录
	logFileDir string
)

// Init 初始化日志结构
func Init(tDebugOn bool, prefix string, tOpenToDB bool) {
	//设定debug
	debugOn = tDebugOn
	//设置转存模式
	openToDB = tOpenToDB
	//设定logger
	globLog.Init()
	globLog.Prefix = prefix + "_glob."
	errLog.Init()
	errLog.Prefix = prefix + "_err."
	warnLog.Init()
	warnLog.Prefix = prefix + "_warn."
	appLog.Init()
	appLog.Prefix = prefix + "_app."
	mqttLog.Init()
	mqttLog.Prefix = prefix + "_mqtt."
	//设定gin
	ginLog.Init()
	ginLog.Prefix = prefix + "_gin."
	//如果不存在则构建log目录
	if !debugOn {
		//如果没启动log则重建目录
		logDir = CoreFile.BaseSrc + CoreFile.Sep + "log"
		//确保log目录存在
		if err := CoreFile.CreateFolder(logDir); err != nil {
			fmt.Println("log set config, create folder is error, " + err.Error())
		}
		//归档目录
		logFileDir = CoreFile.BaseSrc + CoreFile.Sep + "log_file"
		if err := CoreFile.CreateFolder(logFileDir); err != nil {
			fmt.Println("log set config, create folder is error, " + err.Error())
		}
	}
	//预先加载一次
	runMake()
}

// SetAllowGin 设置日志模式
func SetAllowGin(b bool) {
	allowGin = b
}

func SetAllowDefault(b bool) {
	allowDefault = b
}
