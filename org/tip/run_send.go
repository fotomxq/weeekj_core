package OrgTip

import (
	BaseSMS "github.com/fotomxq/weeekj_core/v5/base/sms"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserMessage "github.com/fotomxq/weeekj_core/v5/user/message"
	"time"
)

// runSend 遍历数据并推送消息
func runSend() {
	//获取所有到达推送时间节点的数据
	limit := 100
	step := 0
	for {
		var dataList []FieldsTipType
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_info, from_info, content, params FROM org_tip WHERE tip_at >= NOW() AND delete_at <= to_timestamp(1000000) AND allow_send = false LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			//支持user模块数据推送
			switch v.FromInfo.System {
			case "user":
				//给该用户发送一个消息
				runSendUser(v.FromInfo.ID, &v)
			case "org":
				//给商户的负责人推送消息
				orgData, err := OrgCore.GetOrg(&OrgCore.ArgsGetOrg{
					ID: v.FromInfo.ID,
				})
				if err != nil {
					CoreLog.Error("org tip run error, org not exist, ", err)
					continue
				}
				if orgData.UserID > 0 {
					runSendUser(orgData.UserID, &v)
				}
			default:
				//推送一条日志
				CoreLog.Error("org tip run error, from system not support, tip id: ", v.ID, ", tip data: ", v)
			}
			//标记完成
			if _, err := CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_tip SET allow_send = true WHERE id = :id", map[string]interface{}{
				"id": v.ID,
			}); err != nil {
				CoreLog.Error("org tip run, update allow send failed, ", err)
			}
		}
		step += limit
		//强制延迟10毫秒，消峰处理
		time.Sleep(time.Millisecond * 10)
	}
}

func runSendUser(userID int64, data *FieldsTipType) {
	if _, err := UserMessage.Create(&UserMessage.ArgsCreate{
		WaitSendAt:    CoreFilter.GetNowTime(),
		SendUserID:    data.CreateInfo.ID,
		ReceiveUserID: userID,
		Title:         data.Title,
		Content:       data.Content,
		Files:         data.Files,
		Params:        data.Params,
	}); err != nil {
		CoreLog.Error("org tip run error, send user message, err: ", err)
	}
	if data.NeedSMS && data.SMSConfigID > 0 {
		userData, err := UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
			ID:    userID,
			OrgID: -1,
		})
		if err != nil {
			return
		}
		if userData.NationCode == "" && userData.Phone == "" {
			return
		}
		//发送短信
		_, err = BaseSMS.CreateSMS(&BaseSMS.ArgsCreateSMS{
			OrgID:      data.OrgID,
			ConfigID:   data.SMSConfigID,
			Token:      0,
			NationCode: userData.NationCode,
			Phone:      userData.Phone,
			Params:     nil,
			FromInfo: CoreSQLFrom.FieldsFrom{
				System: "org_tip",
				ID:     0,
				Mark:   "",
				Name:   "商户提醒服务",
			},
		})
		if err != nil {
			CoreLog.Error("org tip run, send sms failed, ", err)
		}
	}
}
