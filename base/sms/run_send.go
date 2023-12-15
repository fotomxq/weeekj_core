package BaseSMS

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

func runSend() {
	limit := 100
	step := 0
	for {
		//获取需要发送的数据集合
		var smsList []FieldsSMS
		if err := Router2SystemConfig.MainDB.Select(&smsList, "SELECT id, create_at, expire_at, send_at, failed_msg, is_check, config_id, token, nation_code, phone, params, from_info FROM core_sms WHERE expire_at >= NOW() AND send_at < to_timestamp(1000000) LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(smsList) < 1 {
			break
		}
		//存在数据
		for _, v := range smsList {
			//获取基础配置
			if v.ConfigID < 1 {
				//记录错误日志
				CoreLog.Error("sms run error, cannot get config data, config lost: ", v.ConfigID)
				//标记失败完成
				runChildOver(v.ID, false, "没有配置，无法发送短信")
				//跳过本条数据
				continue
			}
			configData, err := GetConfigByID(&ArgsGetConfigByID{
				ID:    v.ConfigID,
				OrgID: -1,
			})
			if err != nil {
				//记录错误日志
				CoreLog.Error("sms run error, get config data, config id: ", v.ConfigID, ", err: ", err)
				//标记失败完成
				runChildOver(v.ID, false, "配置无效，无法发送短信")
				//跳过本条数据
				continue
			}
			//根据系统类型，发送短信数据
			isSuccess := false
			failedMsg := ""
			switch configData.System {
			case "tencent":
				//发送腾讯短信
				if err := createSMSToTencent(&configData, &v); err != nil {
					isSuccess = false
					failedMsg = err.Error()
				} else {
					isSuccess = true
					failedMsg = ""
				}
			case "aliyun":
				//发送阿里短信
				if err := createSMSToAliyun(&configData, &v); err != nil {
					isSuccess = false
					failedMsg = err.Error()
				} else {
					isSuccess = true
					failedMsg = ""
				}
			default:
				isSuccess = false
				failedMsg = "无法识别短信系统"
			}
			//标记完成
			runChildOver(v.ID, isSuccess, failedMsg)
			//强制延迟10毫秒，消峰处理
			time.Sleep(time.Millisecond * 5)
		}
		//下一步
		step += limit
	}
}

// 将短信请求作废处理
func runChildOver(smsID int64, isSuccess bool, failedMsg string) {
	defer func() {
		if e := recover(); e != nil {
			CoreLog.Error("sms run error, recover, ", e)
			return
		}
	}()
	if isSuccess {
		if _, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_sms SET send_at = NOW() WHERE id = :id", map[string]interface{}{
			"id": smsID,
		}); err != nil {
			CoreLog.Error("sms run error, update sms send, ", err)
			return
		}
	} else {
		if _, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_sms SET send_at = NOW(), failed_msg = :failed_msg WHERE id = :id", map[string]interface{}{
			"id":         smsID,
			"failed_msg": failedMsg,
		}); err != nil {
			CoreLog.Error("sms run error, update sms failed, ", err)
			return
		}
	}
}
