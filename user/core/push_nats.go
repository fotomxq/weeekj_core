package UserCore

import CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"

// 为用户创建自动化头像
func pushNatsCreateAvatar(userID int64) {
	CoreNats.PushDataNoErr("/user/core/create_avatar", "", userID, "", nil)
}

// 请求发送用户邮件验证等待列队
func pushNatsUserEmailWait(userID int64) {
	CoreNats.PushDataNoErr("/user/core/push_email_wait", "", userID, "", nil)
}

// 请求发送用户邮件
func pushNatsUserEmail(emailWaitID int64) {
	CoreNats.PushDataNoErr("/user/core/push_email", "", emailWaitID, "", nil)
}

// 通知新增用户
func pushNatsCreateUser(userInfo FieldsUserType) {
	CoreNats.PushDataNoErr("/user/core/create_user", "", userInfo.ID, "", userInfo)
}

// 通知用户绑定了手机号
func pushNatsNewPhone(userID int64, nationCode string, phone string) {
	CoreNats.PushDataNoErr("/user/core/new_phone", "", userID, "", map[string]interface{}{
		"nationCode": nationCode,
		"phone":      phone,
	})
}
