package ServiceUserInfo

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//请求删除档案
	CoreNats.SubDataByteNoErr("/service/user_info/post_update", subNatsDeleteInfo)
	//请求统计数据
	CoreNats.SubDataByteNoErr("/service/user_info/analysis", subNatsAnalysis)
	//添加日志
	CoreNats.SubDataByteNoErr("/service/user_info/append_log", subNatsAppendLog)
}

// 请求删除档案
func subNatsDeleteInfo(_ *nats.Msg, action string, id int64, _ string, data []byte) {
	//获取档案
	infoData := getInfoID(id)
	if infoData.ID < 1 {
		return
	}
	//识别处理动作
	switch action {
	case "die":
		dieAtStr := gjson.GetBytes(data, "atTime").String()
		dieAt, err := CoreFilter.GetTimeByDefault(dieAtStr)
		if err != nil {
			CoreLog.Warn("service user info, sub nats, post update. info id: ", infoData.ID, ", die at time, err: ", err)
			return
		}
		if err := UpdateInfoDie(&ArgsUpdateInfoDie{
			ID:     infoData.ID,
			OrgID:  infoData.OrgID,
			DieAt:  CoreFilter.GetISOByTime(dieAt),
			Params: []CoreSQLConfig.FieldsConfigType{},
		}); err != nil {
			CoreLog.Warn("service user info, sub nats, post update. info id: ", infoData.ID, ", update die, err: ", err)
		}
	case "out":
		outAtStr := gjson.GetBytes(data, "atTime").String()
		outAt, err := CoreFilter.GetTimeByDefault(outAtStr)
		if err != nil {
			CoreLog.Warn("service user info, sub nats, post update. info id: ", infoData.ID, ", out at time, err: ", err)
			return
		}
		if err := UpdateInfoOut(&ArgsUpdateInfoOut{
			ID:     infoData.ID,
			OrgID:  infoData.OrgID,
			OutAt:  CoreFilter.GetISOByTime(outAt),
			Params: []CoreSQLConfig.FieldsConfigType{},
		}); err != nil {
			CoreLog.Warn("service user info, sub nats, post update. info id: ", infoData.ID, ", update out, err: ", err)
		}
	case "delete":
		if err := DeleteInfo(&ArgsDeleteInfo{
			ID:    infoData.ID,
			OrgID: infoData.OrgID,
		}); err != nil {
			CoreLog.Warn("service user info, sub nats, post update. info id: ", infoData.ID, ", update delete, err: ", err)
		}
	case "return":
		if err := ReturnInfo(&ArgsReturnInfo{
			ID:    infoData.ID,
			OrgID: infoData.OrgID,
		}); err != nil {
			CoreLog.Warn("service user info, sub nats, post update. info id: ", infoData.ID, ", update return, err: ", err)
		}
	}
}

// 请求统计数据
func subNatsAnalysis(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
	analysisBlockerWait.CheckWait(id, "", func(modID int64, _ string) {
		updateAnalysisOrg(modID)
	})
}

// 添加日志
func subNatsAppendLog(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	var args argsAppendLog
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		return
	}
	appendLog(&args)
}

// 推送请求统计数据
func pushNatsAnalysis(orgID int64) {
	CoreNats.PushDataNoErr("/service/user_info/analysis", "", orgID, "", nil)
}

// 推送档案状态变更
func pushNatsInfoStatus(action string, id int64) {
	CoreNats.PushDataNoErr("/service/user_info/status", action, id, "", nil)
}
