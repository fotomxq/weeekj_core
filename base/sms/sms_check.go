package BaseSMS

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsCheckSMS 验证短信请求参数
type ArgsCheckSMS struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `json:"configID"`
	//会话
	Token int64 `json:"token"`
	//值
	Value string `json:"value"`
}

// CheckSMS 验证短信请求
func CheckSMS(args *ArgsCheckSMS) bool {
	_, b := CheckSMSAndData(&ArgsCheckSMSAndData{
		OrgID:    args.OrgID,
		ConfigID: args.ConfigID,
		Token:    args.Token,
		Value:    args.Value,
	})
	return b
}

// ArgsCheckSMSAndData 验证并反馈验证码的数据结构参数
type ArgsCheckSMSAndData struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `json:"configID"`
	//会话
	Token int64 `json:"token"`
	//值
	Value string `json:"value"`
}

// CheckSMSAndData 验证并反馈验证码的数据结构
func CheckSMSAndData(args *ArgsCheckSMSAndData) (data FieldsSMS, b bool) {
	//获取数据，将自动倒排获取最新一个数据
	if err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, send_at, failed_msg, is_check, config_id, token, nation_code, phone, params, from_info FROM core_sms WHERE is_check = false AND expire_at >= NOW() AND (config_id = $1 OR $1 < 1) AND token = $2 ORDER BY id DESC LIMIT 1", args.ConfigID, args.Token); err != nil {
		return
	}
	if _, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_sms SET is_check = true WHERE id = :id", map[string]interface{}{
		"id": data.ID,
	}); err != nil {
		CoreLog.Error("check sms and data, update history data, id: ", data.ID, ", err: ", err)
		return
	}
	var val string
	val, b = data.Params.GetVal("val")
	if !b {
		val, b = data.Params.GetVal("code")
		if !b {
			return
		}
	}
	if val == args.Value {
		b = true
		return
	}
	return
}
