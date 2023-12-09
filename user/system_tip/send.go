package UserSystemTip

import (
	"fmt"
	UserMessageMod "gitee.com/weeekj/weeekj_core/v5/user/message/mod"
	"time"
)

// SendSuccess 发送通过审核通知
func SendSuccess(userID int64, systemName string, systemID int64, modName string) {
	if modName == "" {
		UserMessageMod.CreateSystemToUser(time.Time{}, userID, fmt.Sprint(systemName, "审核通过"), fmt.Sprint("您申请的", systemName, "已经通过审核"), nil, nil)
	} else {
		UserMessageMod.CreateSystemToUser(time.Time{}, userID, fmt.Sprint(systemName, "审核通过"), fmt.Sprint("您申请的", systemName, "[", modName, "]已经通过审核"), nil, nil)
	}
}
