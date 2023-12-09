package UserCore

import (
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BaseEmail "gitee.com/weeekj/weeekj_core/v5/base/email"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/nats-io/nats.go"
	"strings"
	"time"
)

// 请求发送用户邮件验证等待列队
func subNatsPushEmailWait(_ *nats.Msg, _ string, userID int64, _ string, _ []byte) {
	appendLog := fmt.Sprint("sub nats push email wait, user id: ", userID, ", ")
	//获取用户信息
	userData := getUserByID(userID)
	if userData.ID < 1 || CoreSQL.CheckTimeHaveData(userData.DeleteAt) || CoreSQL.CheckTimeThanNow(userData.EmailVerify) {
		CoreLog.Error(appendLog, "user not exist")
		return
	}
	if userData.Email == "" {
		CoreLog.Error(appendLog, "user not have email")
		return
	}
	//检查该用户是否存在过短请求？
	// 默认5秒以内
	var lastID int64
	err := Router2SystemConfig.MainDB.Get(&lastID, "SELECT id FROM user_reg_wait_email WHERE user_id = $1 AND delete_at < to_timestamp(1000000) AND create_at > CURRENT_TIMESTAMP - INTERVAL '60 second' LIMIT 1;", userData.ID)
	if err == nil && lastID > 0 {
		CoreLog.Error(appendLog, "time too short")
		return
	}
	//如果不存在，则标记该用户其他所有请求为删除状态
	_, _ = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_reg_wait_email SET delete_at = NOW() WHERE user_id = :user_id", map[string]interface{}{
		"user_id": userData.ID,
	})
	//默认过期时间
	var userNewEmailRealTimeStr string
	userNewEmailRealTimeStr, err = BaseConfig.GetDataString("UserNewEmailRealTime")
	if err != nil {
		CoreLog.Error(appendLog, "load config UserNewEmailRealTime, err: ", err)
		return
	}
	var userNewEmailRealTime time.Time
	userNewEmailRealTime, err = CoreFilter.GetTimeByISO(userNewEmailRealTimeStr)
	if err != nil {
		userNewEmailRealTime = CoreFilter.GetNowTime().Add(time.Hour * 1)
	}
	//生成随机码
	var randStr string
	randStr, err = CoreFilter.GetRandStr3(10)
	if err != nil {
		CoreLog.Error(appendLog, "rand have, err: ", err)
		return
	}
	//推送列队
	emailWaitID, err := CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO user_reg_wait_email (expire_at, user_id, email, vcode, is_send) VALUES (:expire_at,:user_id,:email,:vcode,false)", map[string]interface{}{
		"expire_at": userNewEmailRealTime,
		"user_id":   userData.ID,
		"email":     userData.Email,
		"vcode":     randStr,
	})
	if err != nil {
		CoreLog.Error(appendLog, "insert, err: ", err)
		return
	}
	//请求创建
	pushNatsUserEmail(emailWaitID)
	//反馈
	return
}

// 请求发送用户邮件验证
func subNatsPushEmail(_ *nats.Msg, _ string, emailWaitID int64, _ string, _ []byte) {
	appendLog := "sub nats push email, "
	//获取等待发送的ID
	var emailWaitData FieldsRegWaitEmail
	err := Router2SystemConfig.MainDB.Get(&emailWaitData, "SELECT id, create_at, delete_at, expire_at, user_id, email, vcode, is_send FROM user_reg_wait_email WHERE id = $1", emailWaitID)
	if err != nil || emailWaitData.ID < 1 {
		return
	}
	//获取等待推送的邮件数据
	//获取默认发送配置
	userNewEmailRealTitle, err := BaseConfig.GetDataString("UserNewEmailRealTitle")
	if err != nil {
		CoreLog.Error(appendLog, "load config UserNewEmailRealTitle, ", err)
		return
	}
	userNewEmailRealContent, err := BaseConfig.GetDataString("UserNewEmailRealContent")
	if err != nil {
		CoreLog.Error(appendLog, "load config UserNewEmailRealContent, ", err)
		return
	}
	//生成验证连接
	userNewEmailRealSendURL, err := BaseConfig.GetDataString("LoginRegByEmailURL")
	if err != nil {
		CoreLog.Error(appendLog, "load config LoginRegByEmailURL, ", err)
		return
	}
	//修正内容
	sendURL := strings.ReplaceAll(userNewEmailRealSendURL, "{$userID}", fmt.Sprint(emailWaitData.UserID))
	sendURL = strings.ReplaceAll(sendURL, "{$rand}", emailWaitData.VCode)
	content := strings.ReplaceAll(userNewEmailRealContent, "{$url}", sendURL)
	//发送邮件
	if _, err := BaseEmail.Send(&BaseEmail.ArgsSend{
		ServerID: 0,
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     emailWaitData.UserID,
		},
		SendAt:  CoreFilter.GetNowTime(),
		ToEmail: emailWaitData.Email,
		Title:   userNewEmailRealTitle,
		Content: content,
		IsHtml:  true,
	}); err != nil {
		CoreLog.Error(appendLog, "send email, err: ", err)
	}
	//销毁数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_reg_wait_email SET is_send = true WHERE id = :id", map[string]interface{}{
		"id": emailWaitData.ID,
	})
	if err != nil {
		CoreLog.Error(appendLog, "update email wait data, err: ", err)
	}
	//标记所有过期数据
	_, _ = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_reg_wait_email SET delete_at = NOW() WHERE expire_at < NOW()", nil)
}
