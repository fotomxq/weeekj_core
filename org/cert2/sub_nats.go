package OrgCert2

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//缴费成功
	CoreNats.SubDataByteNoErr("/finance/pay/finish", subNatsPayFinish)
	//请求审核证件
	CoreNats.SubDataByteNoErr("/org/cert/audit", subNatsAutoAudit)
	//过期提醒
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsExpire)
	//删除用户操作
	CoreNats.SubDataByteNoErr("/user/core/delete", subNatsDeleteUser)
	//删除组织
	CoreNats.SubDataByteNoErr("/org/core/org", subNatsDeleteOrg)
}

// 通知已经缴费
func subNatsPayFinish(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	logAppend := "org cert2 sub nats update pay finish, "
	switch action {
	case "finish":
		//缴费完成
		// 根据ID标记完成缴费
		if err := updateCertPayFinish(id); err != nil {
			if err.Error() == "no data" {
				//不记录错误
				return
			}
			CoreLog.Warn(logAppend, "pay id: ", id, ", err: ", err)
		}
	}
}

// 请求审核证件
func subNatsAutoAudit(_ *nats.Msg, _ string, certID int64, _ string, _ []byte) {
	logAppend := "org cert2 sub nats audit cert auto, "
	if err := auditCertAuto(certID); err != nil {
		CoreLog.Error(logAppend, "update, ", err)
	}
}

// 过期提醒
func subNatsExpire(_ *nats.Msg, action string, certID int64, _ string, _ []byte) {
	if action != "org_cert" {
		return
	}
	//检查证件是否已经过期
	certData := getCertByID(certID)
	if certData.ID < 1 || CoreSQL.CheckTimeThanNow(certData.ExpireAt) {
		return
	}
	//检查是否为提前30天，还是已经过期数据
	isExpireBefore30 := false
	if CoreSQL.CheckTimeThanNow(certData.ExpireAt) {
		if CoreSQL.CheckTimeThanNow(CoreFilter.GetCarbonByTime(certData.ExpireAt).SubDays(30).Time) {
			return
		}
		isExpireBefore30 = true
	} else {
		//已经过期数据
	}
	//日志前缀
	logAppend := "org cert2 sub nats expire, "
	//获取配置
	configData := getConfigByID(certData.ConfigID)
	if configData.ID < 1 {
		CoreLog.Error(logAppend, "get config id: ", certData.ConfigID)
		return
	}
	var params CoreSQLConfig.FieldsConfigsType
	needUserMessage, b := configData.Params.GetValBool("needUserMessage")
	if b {
		needUserMessageStr, err := CoreFilter.GetStringByInterface(needUserMessage)
		if err != nil {
			needUserMessageStr = "false"
		}
		params = append(params, CoreSQLConfig.FieldsConfigType{
			Mark: "needUserMessage",
			Val:  needUserMessageStr,
		})
	}
	needSMS, b := configData.Params.GetValBool("needSMS")
	if b {
		needSMSStr, err := CoreFilter.GetStringByInterface(needSMS)
		if err != nil {
			needSMSStr = "false"
		}
		params = append(params, CoreSQLConfig.FieldsConfigType{
			Mark: "needSMS",
			Val:  needSMSStr,
		})
	}
	smsConfigID, b := configData.Params.GetValInt64("sms30ConfigID")
	if b && smsConfigID > 0 {
		//写入预警消息后，自动归类为短信配置ID
		params = append(params, CoreSQLConfig.FieldsConfigType{
			Mark: "smsConfigID",
			Val:  fmt.Sprint(smsConfigID),
		})
	}
	var err error
	if isExpireBefore30 {
		err = createWarning(&argsCreateWarning{
			OrgID:      certData.OrgID,
			CertID:     certData.ID,
			Msg:        fmt.Sprintf("证件《%s》还有1个月过期，请尽快处理", configData.Name),
			Params:     params,
			ConfigName: configData.Name,
			BindFrom:   configData.BindFrom,
			BindID:     certData.BindID,
		})
	} else {
		err = createWarning(&argsCreateWarning{
			OrgID:      certData.OrgID,
			CertID:     certData.ID,
			Msg:        fmt.Sprintf("证件《%s》已经过期", configData.Name),
			Params:     params,
			ConfigName: configData.Name,
			BindFrom:   configData.BindFrom,
			BindID:     certData.BindID,
		})
	}
	if err != nil {
		CoreLog.Error(logAppend, "create warning failed, ", err)
	}
}

// 删除用户操作
func subNatsDeleteUser(_ *nats.Msg, _ string, userID int64, _ string, _ []byte) {
	//日志前缀
	logAppend := "org cert2 sub nats delete user, "
	//删除该用户的所有证件信息
	err := DeleteAllCertByBindID("user", userID)
	if err != nil {
		CoreLog.Error(logAppend, "delete user all cert failed, user id: ", userID, ", err: ", err)
		return
	}
}

// 删除组织操作
func subNatsDeleteOrg(_ *nats.Msg, action string, orgID int64, _ string, _ []byte) {
	if action != "delete" {
		return
	}
	//日志前缀
	logAppend := "org cert2 sub nats delete org, "
	//删除该用户的所有证件信息
	err := DeleteAllCertByBindID("org", orgID)
	if err != nil {
		CoreLog.Error(logAppend, "delete org all cert failed, org id: ", orgID, ", err: ", err)
		return
	}
}
