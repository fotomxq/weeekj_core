package OrgUser

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//请求更新用户数据
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "组织用户提交更新",
		Description:  "",
		EventSubType: "all",
		Code:         "org_user_post_update",
		EventType:    "nats",
		EventURL:     "/org/user/post_update",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("org_user_post_update", "/org/user/post_update", subNatsUpdateUserData)
	//删除过期数据
	CoreNats.SubDataByteNoErr("base_expire_tip_expire", "/base/expire_tip/expire", subNatsDeleteExpire)
}

// subNatsUpdateUserData 请求更新用户数据
func subNatsUpdateUserData(_ *nats.Msg, _ string, id int64, mark string, _ []byte) {
	waitUpdateBlockerWait.CheckWait(id, mark, func(modID int64, modMark string) {
		orgID, _ := CoreFilter.GetInt64ByString(modMark)
		if orgID < 1 {
			userData, _ := UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
				ID:    modID,
				OrgID: -1,
			})
			if userData.ID > 0 {
				orgID = userData.OrgID
			}
		}
		if orgID > 0 && modID > 0 {
			if err := updateByUserID(orgID, modID); err != nil {
				CoreLog.Error("org user data sub nats update user data, org id: ", orgID, ", user id: ", modID, ", err: ", err)
			}
		}
	})
}

// subNatsDeleteExpire 删除过期数据
func subNatsDeleteExpire(_ *nats.Msg, action string, _ int64, _ string, data []byte) {
	if action != "org_user_data" {
		return
	}
	orgID := gjson.GetBytes(data, "orgID").Int()
	userID := gjson.GetBytes(data, "userID").Int()
	if orgID < 1 || userID < 1 {
		return
	}
	_ = DeleteDataByUserID(&ArgsDeleteDataByUserID{
		UserID: userID,
		OrgID:  orgID,
	})
}
