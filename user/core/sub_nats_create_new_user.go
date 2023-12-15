package UserCore

import (
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	OrgUserMod "github.com/fotomxq/weeekj_core/v5/org/user/mod"
	"github.com/nats-io/nats.go"
	"time"
)

func subNatsCreateNewUser(_ *nats.Msg, _ string, userID int64, _ string, rawData []byte) {
	appendLog := "sub nats create new user, "
	//获取参数集合
	var userInfo FieldsUserType
	if err := CoreNats.ReflectDataByte(rawData, &userInfo); err != nil {
		CoreLog.Error(appendLog, "get params by raw data, ", err)
		return
	}
	//请求构建用户头像
	pushNatsCreateAvatar(userID)
	//推送email验证处理
	if !CoreSQL.CheckTimeHaveData(userInfo.EmailVerify) && userInfo.Email != "" {
		pushNatsUserEmailWait(userInfo.ID)
	}
	//推送sms验证处理
	if !CoreSQL.CheckTimeHaveData(userInfo.PhoneVerify) && userInfo.NationCode != "" && userInfo.Phone != "" {
		pushNatsNewPhone(userInfo.ID, userInfo.NationCode, userInfo.Phone)
	}
	//通知组织用户更新集合
	OrgUserMod.PushUpdateUserData(userInfo.OrgID, userInfo.ID)
	//统计行为
	AnalysisAny2.AppendData("add", "user_core_new_count", time.Time{}, userInfo.OrgID, 0, 0, 0, 0, 1)
}
