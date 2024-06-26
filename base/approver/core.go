package BaseApprover

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

//审批流
/**
1. 用于任意模块审批流设计和实现
*/

var (
	//缓冲时间
	cacheConfigTime     = 1800
	cacheConfigItemTime = 1800
	cacheLogTime        = 1800
	cacheLogFlowTime    = 1800
	//数据表
	configDB     CoreSQL2.Client
	configItemDB CoreSQL2.Client
	logDB        CoreSQL2.Client
	logFlowDB    CoreSQL2.Client
	//OpenSub 消息
	OpenSub = false
)

// Init 初始化
func Init() {
	//初始化数据表
	configDB.Init(&Router2SystemConfig.MainSQL, "base_approver_config")
	configItemDB.Init(&Router2SystemConfig.MainSQL, "base_approver_config_item")
	logDB.Init(&Router2SystemConfig.MainSQL, "base_approver_log")
	logFlowDB.Init(&Router2SystemConfig.MainSQL, "base_approver_log_flow")
	//启动消息
	if OpenSub {
		subNats()
	}
}

// getApproverName 通过组织成员ID或用户ID获取姓名
func getApproverName(orgBindID int64, userID int64) (name string) {
	var approverName string
	if orgBindID > 0 {
		approverName = OrgCore.GetBindName(orgBindID)
	}
	if approverName == "" && userID > 0 {
		approverName = UserCore.GetUserNameByID(userID)
	}
	return
}
