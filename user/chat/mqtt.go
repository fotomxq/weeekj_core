package UserChat

import (
	"encoding/json"
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	IOTMQTT "github.com/fotomxq/weeekj_core/v5/iot/mqtt"
)

// MQTT订阅
func initSub() {
	//暂无订阅处理
}

// 推送给用户消息更新了
func pushUpdateByUser(userID int64, groupID int64, action int) {
	//构建数据
	type newDataType struct {
		//组ID
		GroupID int64 `json:"groupID"`
		//行为特征
		// 0 创建; 1 有更新; 2 删除
		Action int `json:"action"`
	}
	newData := newDataType{
		GroupID: groupID,
		Action:  action,
	}
	//获取和打包数据
	dataByte, err := json.Marshal(newData)
	if err != nil {
		CoreLog.Error("user chat push mqtt, push update by user, json err: ", err)
		return
	}
	topic := fmt.Sprint("user/chat/user/", userID)
	if err = IOTMQTT.MQTTClient.PublishWait(topic, 0, false, dataByte); err != nil {
		CoreLog.Warn("user chat push mqtt, push update by user, push data: ", err)
		return
	}
}

// 推送更新给聊天室
func pushUpdateByGroup(groupID int64, action int) {
	//构建数据
	type newDataType struct {
		//组ID
		GroupID int64 `json:"groupID"`
		//行为特征
		// 0 新人加入; 1 有人退出; 2 新的消息; 3 更新成员姓名
		Action int `json:"action"`
	}
	newData := newDataType{
		GroupID: groupID,
		Action:  action,
	}
	//获取和打包数据
	dataByte, err := json.Marshal(newData)
	if err != nil {
		CoreLog.Error("user chat push mqtt, push update by group, json err: ", err)
		return
	}
	topic := fmt.Sprint("user/chat/group/", groupID)
	if err = IOTMQTT.MQTTClient.PublishWait(topic, 0, false, dataByte); err != nil {
		CoreLog.Warn("user chat push mqtt, push update by group, push data: ", err)
		return
	}
}
