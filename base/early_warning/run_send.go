package BaseEarlyWarning

import (
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BaseEmail "gitee.com/weeekj/weeekj_core/v5/base/email"
	BaseSMS "gitee.com/weeekj/weeekj_core/v5/base/sms"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"strings"
)

// 推送消息
func runSend() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("early warning run, ", r)
		}
	}()
	var err error
	//获取等待发送的通知
	limit := 100
	step := 0
	for {
		//获取数据
		var waitSendList []FieldsWaitType
		err = Router2SystemConfig.MainDB.Select(&waitSendList, "SELECT id, create_at, update_at, bind_id, level, to_id, is_send, is_read, expire_finish, expire_at, template_id, content, bind_data, need_phone, need_sms, need_email, need_app FROM core_ew_wait WHERE is_send = false AND is_read = false AND expire_finish = false LIMIT $1 OFFSET $2", limit, step)
		if err != nil {
			break
		}
		if len(waitSendList) < 1 {
			break
		}
		//遍历数据，开始推送数据
		for _, v := range waitSendList {
			//获取模版数据
			templateData, err := GetTemplateByID(&ArgsGetTemplateByID{
				ID: v.TemplateID,
			})
			if err != nil {
				err = updateSendExpiredAndNext(&v)
				if err != nil {
					CoreLog.Error("early warning run, cannot get template data, ", err)
				}
				continue
			}
			//联系人信息
			toData, err := GetToByID(&ArgsGetToByID{
				ID: v.ToID,
			})
			if err != nil {
				err = updateSendExpiredAndNext(&v)
				if err != nil {
					CoreLog.Error("early warning run, cannot get to data, ", err)
				}
				continue
			}
			//根据模版类型，发送不同数据
			if v.NeedSMS {
				smsConfigID, err := BaseConfig.GetDataInt64("EarlyWarningSMSConfigID")
				if err != nil {
					CoreLog.Error("early warning run, load config by sms config id, ", err)
				} else {
					if smsConfigID > 0 && toData.Phone != "" {
						if _, err := BaseSMS.CreateSMS(&BaseSMS.ArgsCreateSMS{
							OrgID:      0,
							ConfigID:   smsConfigID,
							Token:      0,
							NationCode: toData.PhoneNationCode,
							Phone:      toData.Phone,
							Params:     nil,
							FromInfo: CoreSQLFrom.FieldsFrom{
								System: "early-warning",
								ID:     0,
								Mark:   templateData.Mark,
								Name:   templateData.Name,
							},
						}); err != nil {
							CoreLog.Error("early warning run, cannot send sms, ", err)
							continue
						}
					}
				}
			}
			if v.NeedEmail && toData.Email != "" {
				content := v.Content
				for _, v2 := range templateData.BindData {
					replaceContent := ""
					isFind := false
					for k3Mark := range v.BindData {
						if k3Mark == v2 {
							replaceContent = v2
							isFind = true
						}
					}
					if !isFind {
						CoreLog.Error("early warning run, cannot send email, template config is error, ", err)
						break
					}
					content = strings.ReplaceAll(content, v2, replaceContent)
				}
				if _, err := BaseEmail.Send(&BaseEmail.ArgsSend{
					CreateInfo: CoreSQLFrom.FieldsFrom{
						System: "early-warning",
						ID:     0,
						Mark:   "system",
					},
					SendAt:  CoreFilter.GetNowTime(),
					ToEmail: toData.Email,
					Title:   templateData.Title,
					Content: content,
				}); err != nil {
					CoreLog.Error("early warning run, cannot send email, send error, ", err)
				}
			}
		}
		//下一页
		step += limit
	}
}
