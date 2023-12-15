package UserMessageMod

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// ArgsCreate 创建新的消息参数
type ArgsCreate struct {
	//预计发送时间
	WaitSendAt time.Time `json:"waitSendAt" check:"isoTime" empty:"true"`
	//发送人
	// 如果为0则为系统消息，同时自动跳过时间差验证机制
	SendUserID int64 `json:"sendUserID" check:"id"`
	//接收人
	ReceiveUserID int64 `json:"receiveUserID" check:"id"`
	//标题
	Title string `json:"title" check:"des" min:"1" max:"300"`
	//内容
	Content string `json:"content" check:"des" min:"1" max:"1500"`
	//附件文件列
	Files []int64 `json:"files" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `json:"params"`
}

// Create 创建新的消息
func Create(args ArgsCreate) {
	CoreNats.PushDataNoErr("/user/message/create", "user", 0, "", args)
}

// CreateSystemToUser 发送系统消息
func CreateSystemToUser(waitSendAt time.Time, receiveUserID int64, title, content string, fileID []int64, params CoreSQLConfig.FieldsConfigsType) {
	if waitSendAt.Unix() < 1000000 {
		waitSendAt = CoreFilter.GetNowTime()
	}
	if fileID == nil || len(fileID) < 1 {
		fileID = []int64{}
	}
	if params == nil || len(params) < 1 {
		params = CoreSQLConfig.FieldsConfigsType{}
	}
	CoreNats.PushDataNoErr("/user/message/create", "user", 0, "", ArgsCreate{
		WaitSendAt:    waitSendAt,
		SendUserID:    0,
		ReceiveUserID: receiveUserID,
		Title:         title,
		Content:       content,
		Files:         fileID,
		Params:        params,
	})
}

// CreateSystemToAllUser 发送全局用户消息
func CreateSystemToAllUser(waitSendAt time.Time, title, content string, fileID []int64, params CoreSQLConfig.FieldsConfigsType) {
	if waitSendAt.Unix() < 1000000 {
		waitSendAt = CoreFilter.GetNowTime()
	}
	CoreNats.PushDataNoErr("/user/message/create", "all", 0, "", ArgsCreate{
		WaitSendAt:    waitSendAt,
		SendUserID:    0,
		ReceiveUserID: 0,
		Title:         title,
		Content:       content,
		Files:         fileID,
		Params:        params,
	})
}
