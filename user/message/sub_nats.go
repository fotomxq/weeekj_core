package UserMessage

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	UserCoreMod "github.com/fotomxq/weeekj_core/v5/user/core/mod"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//通知发送消息
	CoreNats.SubDataByteNoErr("base_expire_tip_expire", "/base/expire_tip/expire", subNatsWaitSend)
	//请求审核消息
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "用户消息审核通知",
		Description:  "",
		EventSubType: "all",
		Code:         "user_message_audit",
		EventType:    "nats",
		EventURL:     "/user/message/audit",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("user_message_audit", "/user/message/audit", subNatsAudit)
	//创建消息
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "用户消息创建通知",
		Description:  "",
		EventSubType: "all",
		Code:         "user_message_create",
		EventType:    "nats",
		EventURL:     "/user/message/create",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("user_message_create", "/user/message/create", subNatsCreate)
}

// 通知发送消息
func subNatsWaitSend(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	//如果系统不符合，跳出
	if action != "user_message" {
		return
	}
	//修改状态为等待审核
	if err := UpdatePost(&ArgsUpdatePost{
		ID:         id,
		SendUserID: -1,
	}); err != nil {
		CoreLog.Warn("user message sub nats wait send, message id: ", id, ", update post failed, ", err)
		return
	}
	//请求自动审核消息
	pushNatsAutoAudit(id)
}

// 请求审核消息
func subNatsAudit(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
	if err := UpdateAudit(&ArgsUpdateAudit{
		ID: id,
	}); err != nil {
		CoreLog.Warn("user message sub nats wait send, message id: ", id, ", update audit failed, ", err)
		return
	}
}

// 创建消息
func subNatsCreate(_ *nats.Msg, action string, _ int64, _ string, data []byte) {
	appendLog := "user message sub nats wait send, "
	//解析数据
	var args ArgsCreate
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Warn(appendLog, "get params, ", err)
		return
	}
	switch action {
	case "all":
		step := 0
		for {
			userList := UserCoreMod.GetAllUserList(-1, 2, -1, []int64{}, step, 1000)
			if len(userList) < 1 {
				break
			}
			for _, vUser := range userList {
				args.SendUserID = vUser.ID
				if _, err := create(&args); err != nil {
					CoreLog.Warn("user message sub nats wait send, create message, ", err)
					continue
				}
			}
			step += 1000
		}
	case "user":
		if _, err := create(&args); err != nil {
			CoreLog.Warn("user message sub nats wait send, create message, ", err)
			return
		}
	}
}
