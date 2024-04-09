package UserCore

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
)

func subNats() {
	//为用户创建自动化头像
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "用户创建新的头像",
		Description:  "",
		EventSubType: "all",
		Code:         "user_core_create_avatar",
		EventType:    "nats",
		EventURL:     "/user/core/create_avatar",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("user_core_create_avatar", "/user/core/create_avatar", subNatsCreateAvatar)
	//请求发送用户邮件验证等待列队
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "给用户发送新的验证邮件",
		Description:  "",
		EventSubType: "all",
		Code:         "user_core_push_email_wait",
		EventType:    "nats",
		EventURL:     "/user/core/push_email_wait",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("user_core_push_email_wait", "/user/core/push_email_wait", subNatsPushEmailWait)
	//请求发送用户邮件验证
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "给用户发送新的验证邮件并完成验证",
		Description:  "",
		EventSubType: "all",
		Code:         "user_core_push_email",
		EventType:    "nats",
		EventURL:     "/user/core/push_email",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("user_core_push_email", "/user/core/push_email", subNatsPushEmail)
	//创建新用户的后续处理
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "创建新的用户",
		Description:  "",
		EventSubType: "all",
		Code:         "user_core_create_user",
		EventType:    "nats",
		EventURL:     "/user/core/create_user",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("user_core_create_user", "/user/core/create_user", subNatsCreateNewUser)
	//推送服务
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "用户更换新的手机号",
		Description:  "",
		EventSubType: "all",
		Code:         "user_core_new_phone",
		EventType:    "nats",
		EventURL:     "/user/core/new_phone",
		//TODO:待补充
		EventParams: "",
	})
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "用户被删除通知",
		Description:  "",
		EventSubType: "all",
		Code:         "user_core_delete",
		EventType:    "nats",
		EventURL:     "/user/core/delete",
		//TODO:待补充
		EventParams: "",
	})
}
