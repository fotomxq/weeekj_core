package IOTMQTT

import (
	"sync"
)

type waitPushUpdateDateListType struct {
	OrgID    int64
	Action   string
	UpdateID int64
}

var (
	//等待推送列队
	waitPushUpdateDateList     []waitPushUpdateDateListType
	waitPushUpdateDateListLock sync.Mutex
)

//PushUpdateData 全局通知业务层面的更新处理
/**
1. 根据组织ID来进行通告
2. 业务设备收到信息后，根据ID或进行业务层面的强制更新数据
*/
func PushUpdateData(orgID int64, action string, updateID int64) (err error) {
	waitPushUpdateDateListLock.Lock()
	defer waitPushUpdateDateListLock.Unlock()
	for _, v := range waitPushUpdateDateList {
		if v.OrgID == orgID && v.Action == action && v.UpdateID == updateID {
			return
		}
	}
	waitPushUpdateDateList = append(waitPushUpdateDateList, waitPushUpdateDateListType{
		OrgID:    orgID,
		Action:   action,
		UpdateID: updateID,
	})
	return
}
