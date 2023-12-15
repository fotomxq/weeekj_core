package UserCore

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
)

func subNats() {
	//为用户创建自动化头像
	CoreNats.SubDataByteNoErr("/user/core/create_avatar", subNatsCreateAvatar)
	//请求发送用户邮件验证等待列队
	CoreNats.SubDataByteNoErr("/user/core/push_email_wait", subNatsPushEmailWait)
	//请求发送用户邮件验证
	CoreNats.SubDataByteNoErr("/user/core/push_email", subNatsPushEmail)
	//创建新用户的后续处理
	CoreNats.SubDataByteNoErr("/user/core/create_user", subNatsCreateNewUser)
}
