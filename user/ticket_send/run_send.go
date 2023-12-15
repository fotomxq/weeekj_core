package UserTicketSend

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserSubscriptionMod "github.com/fotomxq/weeekj_core/v5/user/subscription/mod"
	UserTicket "github.com/fotomxq/weeekj_core/v5/user/ticket"
	"time"
)

func runSend() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("user ticket send run, ", r)
		}
	}()
	step := 0
	limit := 3
	for {
		//遍历配置
		var dataList []FieldsSend
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, org_id, send_count, need_user_sub_config_id, config_id, per_count FROM user_ticket_send WHERE finish_at < to_timestamp(1000000) AND need_auto = true LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		//遍历数据
		for _, vSend := range dataList {
			//检查配置是否还存在？
			_, err := UserTicket.GetConfigByID(&UserTicket.ArgsGetConfigByID{
				ID:    vSend.ConfigID,
				OrgID: vSend.OrgID,
			})
			if err != nil {
				//删除该配置并跳过
				_ = DeleteSend(&ArgsDeleteSend{
					ID:    vSend.ID,
					OrgID: vSend.OrgID,
				})
				continue
			}
			//赠送数量
			var sendCount int64 = 0
			//遍历用户
			var userPage int64 = 1
			for {
				userList, _, err := UserCore.GetUserList(&UserCore.ArgsGetUserList{
					Pages: CoreSQLPages.ArgsDataList{
						Page: userPage,
						Max:  1000,
						Sort: "id",
						Desc: false,
					},
					Status:       2,
					OrgID:        vSend.OrgID,
					ParentSystem: "",
					ParentID:     -1,
					SortID:       -1,
					Tags:         []int64{},
					IsRemove:     false,
					Search:       "",
				})
				if err != nil || len(userList) < 1 {
					break
				}
				for _, vUser := range userList {
					if b := runSendUser(vUser.ID, &vSend); b {
						sendCount += 1
					}
				}
				//更新发送量
				_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_ticket_send SET send_count = :send_count WHERE id = :id", map[string]interface{}{
					"id":         vSend.ID,
					"send_count": sendCount,
				})
				userPage += 1
				//消峰处理
				time.Sleep(time.Millisecond * 500)
			}
			//全部完成后，标记该自动化处理完成
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_ticket_send SET finish_at = NOW() WHERE id = :id", map[string]interface{}{
				"id": vSend.ID,
			})
			if err != nil {
				CoreLog.Error("user ticket send auto send failed, update send finish, send id: ", vSend.ID, ", err: ", err)
			}
			//消峰处理
			time.Sleep(time.Millisecond * 500)
		}
		//下一页
		step += limit
	}
}

// 赠送内部逻辑
func runSendUser(userID int64, sendData *FieldsSend) (b bool) {
	//检查是否赠送过
	var logID int64
	err := Router2SystemConfig.MainDB.Get(&logID, "SELECT id FROM user_ticket_send_log WHERE user_id = $1 AND send_id = $2", userID, sendData.ID)
	if err == nil && logID > 0 {
		return true
	}
	//检查用户的会员是否符合条件
	if sendData.NeedUserSubConfigID > 0 {
		b = UserSubscriptionMod.CheckSub(&UserSubscriptionMod.ArgsCheckSub{
			ConfigID: sendData.NeedUserSubConfigID,
			UserID:   userID,
		})
		if !b {
			return
		}
	}
	//赠送票据
	err = UserTicket.AddTicket(&UserTicket.ArgsAddTicket{
		OrgID:       sendData.OrgID,
		ConfigID:    sendData.ConfigID,
		UserID:      userID,
		Count:       sendData.PerCount,
		UseFromName: "auto",
	})
	if err != nil {
		CoreLog.Warn("user ticket send auto send failed, add ticket, ", err, ", user id: ", userID, ", send id: ", sendData.ID, ", ticket config id: ", sendData.ConfigID)
		return
	}
	//标记成功
	b = true
	//记录日志
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_ticket_send_log (send_id, user_id) VALUES (:send_id,:user_id)", map[string]interface{}{
		"send_id": sendData.ID,
		"user_id": userID,
	})
	if err != nil {
		CoreLog.Error("user ticket send auto send failed, insert log, ", err, ", user id: ", userID, ", send id: ", sendData.ID, ", ticket config id: ", sendData.ConfigID)
		return
	}
	return
}
