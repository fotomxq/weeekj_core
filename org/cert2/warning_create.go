package OrgCert2

import (
	"errors"
	"fmt"
	BaseSMS "github.com/fotomxq/weeekj_core/v5/base/sms"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserMessage "github.com/fotomxq/weeekj_core/v5/user/message"
	"time"
)

// argsCreateWarning 创建新的记录参数
type argsCreateWarning struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//异常证件ID
	CertID int64 `db:"cert_id" json:"certID" check:"id"`
	//消息
	Msg string `db:"msg" json:"msg" check:"des" min:"1" max:"600"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//配置标题
	ConfigName string `json:"configName"`
	//绑定来源
	// user 用户 / org 商户 / org_bind 商户成员 / finance_assets 财务资产 /
	BindFrom string `db:"bind_from" json:"bindFrom"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID"`
}

// createWarning 创建新的记录
func createWarning(args *argsCreateWarning) (err error) {
	//获取证件信息
	var certData FieldsCert
	certData, err = GetCert(&ArgsGetCert{
		ID:    args.CertID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	//获取配置
	var configData FieldsConfig
	configData, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    certData.ConfigID,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	//检查是否存在未处理
	var data FieldsWarning
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at FROM org_cert_warning2 WHERE org_id = $1 AND cert_id = $2 AND finish_at < to_timestamp(1000000)", args.OrgID, args.CertID)
	if err == nil && data.ID > 0 {
		//检查是否存在每年提醒开关？
		needTipEveryYear, b := configData.Params.GetValBool("needTipEveryYear")
		if !b {
			needTipEveryYear = false
		}
		if !needTipEveryYear {
			err = errors.New("have warning")
			return
		} else {
			//检查是否达到一年
			if data.CreateAt.Unix() >= CoreFilter.GetNowTimeCarbon().SubMonths(11).Time.Unix() {
				err = errors.New("have warning")
				return
			}
		}
	}
	//创建记录
	var warningData FieldsWarning
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_cert_warning2", "INSERT INTO org_cert_warning2 (org_id, cert_id, config_id, config_mark, msg, send_msg_at, send_sms_at, params) VALUES (:org_id, :cert_id, :config_id, :config_mark, :msg, to_timestamp(0), to_timestamp(0), :params)", map[string]interface{}{
		"org_id":      args.OrgID,
		"cert_id":     args.CertID,
		"config_id":   configData.ID,
		"config_mark": configData.Mark,
		"msg":         args.Msg,
		"params":      args.Params,
	}, &warningData)
	if err == nil {
		//根据需求，发送用户消息和短信消息
		needUserMessage, b := args.Params.GetValBool("needUserMessage")
		if b && needUserMessage {
			switch args.BindFrom {
			case "org":
				var orgData OrgCoreCore.FieldsOrg
				orgData, err = OrgCoreCore.GetOrg(&OrgCoreCore.ArgsGetOrg{
					ID: args.BindID,
				})
				if err != nil {
					break
				}
				var params CoreSQLConfig.FieldsConfigsType
				params = append(params, CoreSQLConfig.FieldsConfigType{
					Mark: "orgCertWarningID",
					Val:  fmt.Sprint(warningData.ID),
				})
				_, err = UserMessage.Create(&UserMessage.ArgsCreate{
					WaitSendAt:    time.Now(),
					SendUserID:    0,
					ReceiveUserID: orgData.UserID,
					Title:         args.Msg,
					Content:       args.Msg,
					Files:         []int64{},
					Params:        params,
				})
				if err != nil {
					break
				}
			case "user":
				var params CoreSQLConfig.FieldsConfigsType
				params = append(params, CoreSQLConfig.FieldsConfigType{
					Mark: "orgCertWarningID",
					Val:  fmt.Sprint(warningData.ID),
				})
				_, err = UserMessage.Create(&UserMessage.ArgsCreate{
					WaitSendAt:    time.Now(),
					SendUserID:    0,
					ReceiveUserID: args.BindID,
					Title:         args.Msg,
					Content:       args.Msg,
					Files:         []int64{},
					Params:        params,
				})
				if err != nil {
					break
				}
			}
		}
		needSMS, b := args.Params.GetValBool("needSMS")
		if b && needSMS {
			smsConfigID, b := args.Params.GetValInt64("smsConfigID")
			if b && smsConfigID > 0 {
				switch args.BindFrom {
				case "org":
					var orgData OrgCoreCore.FieldsOrg
					orgData, err = OrgCoreCore.GetOrg(&OrgCoreCore.ArgsGetOrg{
						ID: args.BindID,
					})
					if err != nil {
						break
					}
					var userData UserCore.FieldsUserType
					userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
						ID:    orgData.UserID,
						OrgID: -1,
					})
					if err != nil {
						break
					}
					if userData.NationCode == "" || userData.Phone == "" {
						break
					}
					var params CoreSQLConfig.FieldsConfigsType
					params = append(params, CoreSQLConfig.FieldsConfigType{
						Mark: "$1",
						Val:  args.ConfigName,
					})
					_, err = BaseSMS.CreateSMS(&BaseSMS.ArgsCreateSMS{
						OrgID:      orgData.ID,
						ConfigID:   smsConfigID,
						Token:      0,
						NationCode: userData.NationCode,
						Phone:      userData.Phone,
						Params:     params,
						FromInfo: CoreSQLFrom.FieldsFrom{
							System: "",
							ID:     0,
							Mark:   "",
							Name:   "组织证件",
						},
					})
					if err != nil {
						break
					}
				case "user":
					var userData UserCore.FieldsUserType
					userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
						ID:    args.BindID,
						OrgID: -1,
					})
					if err != nil {
						break
					}
					if userData.NationCode == "" || userData.Phone == "" {
						break
					}
					var params CoreSQLConfig.FieldsConfigsType
					params = append(params, CoreSQLConfig.FieldsConfigType{
						Mark: "$1",
						Val:  args.ConfigName,
					})
					_, err = BaseSMS.CreateSMS(&BaseSMS.ArgsCreateSMS{
						OrgID:      userData.OrgID,
						ConfigID:   smsConfigID,
						Token:      0,
						NationCode: userData.NationCode,
						Phone:      userData.Phone,
						Params:     params,
						FromInfo: CoreSQLFrom.FieldsFrom{
							System: "",
							ID:     0,
							Mark:   "",
							Name:   "组织证件",
						},
					})
					if err != nil {
						break
					}
				}
			}
		}
	}
	return
}
