package OrgCert

import (
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

func runWarningCreate() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("org cert warning run, ", r)
		}
	}()
	//加载配置项
	orgCertWarningCreateTime, err := BaseConfig.GetDataInt("OrgCertWarningCreateTime")
	if err != nil {
		orgCertWarningCreateTime = 3
	}
	//预加载配置列
	var configList []FieldsConfig
	//初始化
	var dataList []FieldsCert
	limit := 100
	step := 0
	//检查即将过期的数据，提前30天警告
	for {
		if orgCertWarningCreateTime < 1 {
			if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, org_id, child_org_id, config_id, bind_id FROM org_cert WHERE delete_at < to_timestamp(1000000) AND expire_at < $1 AND expire_at >= NOW() LIMIT $2 OFFSET $3", CoreFilter.GetNowTimeCarbon().AddMonth().Time, limit, step); err != nil {
				break
			}
		} else {
			if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, org_id, child_org_id, config_id, bind_id FROM org_cert WHERE delete_at < to_timestamp(1000000) AND expire_at < $1 AND expire_at >= NOW() AND create_at < $2 LIMIT 100 OFFSET "+fmt.Sprint(step), CoreFilter.GetNowTimeCarbon().AddMonth().Time, CoreFilter.GetNowTimeCarbon().SubMonths(orgCertWarningCreateTime).Time); err != nil {
				break
			}
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			var vConfig FieldsConfig
			var b bool
			configList, vConfig, b = runWarningLoadConfig(configList, v.ConfigID)
			if !b {
				continue
			}
			var params CoreSQLConfig.FieldsConfigsType
			needUserMessage, b := vConfig.Params.GetValBool("needUserMessage")
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
			needSMS, b := vConfig.Params.GetValBool("needSMS")
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
			smsConfigID, b := vConfig.Params.GetValInt64("sms30ConfigID")
			if b && smsConfigID > 0 {
				//写入预警消息后，自动归类为短信配置ID
				params = append(params, CoreSQLConfig.FieldsConfigType{
					Mark: "smsConfigID",
					Val:  fmt.Sprint(smsConfigID),
				})
			}
			_ = createWarning(&argsCreateWarning{
				OrgID:      v.OrgID,
				ChildOrgID: v.ChildOrgID,
				CertID:     v.ID,
				Msg:        fmt.Sprintf("证件《%s》还有1个月过期，请尽快处理", vConfig.Name),
				Params:     params,
				ConfigName: vConfig.Name,
				BindFrom:   vConfig.BindFrom,
				BindID:     v.BindID,
			})
		}
		step += limit
		time.Sleep(time.Millisecond * 500)
	}
	//检查所有已经过期的数据
	step = 0
	for {
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, org_id, child_org_id, config_id, bind_id FROM org_cert WHERE delete_at < to_timestamp(1000000) AND expire_at < NOW() LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			var vConfig FieldsConfig
			var b bool
			configList, vConfig, b = runWarningLoadConfig(configList, v.ConfigID)
			if !b {
				continue
			}
			var params CoreSQLConfig.FieldsConfigsType
			needUserMessage, b := vConfig.Params.GetValBool("needUserMessage")
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
			needSMS, b := vConfig.Params.GetValBool("needSMS")
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
			smsConfigID, b := vConfig.Params.GetValInt64("smsExpireConfigID")
			if b && smsConfigID > 0 {
				//写入预警消息后，自动归类为短信配置ID
				params = append(params, CoreSQLConfig.FieldsConfigType{
					Mark: "smsConfigID",
					Val:  fmt.Sprint(smsConfigID),
				})
			}
			_ = createWarning(&argsCreateWarning{
				OrgID:      v.OrgID,
				ChildOrgID: v.ChildOrgID,
				CertID:     v.ID,
				Msg:        fmt.Sprintf("证件《%s》已经过期", vConfig.Name),
				Params:     params,
				ConfigName: vConfig.Name,
				BindFrom:   vConfig.BindFrom,
				BindID:     v.BindID,
			})
		}
		step += limit
		time.Sleep(time.Millisecond * 500)
	}
}

func runWarningLoadConfig(configList []FieldsConfig, addConfigID int64) ([]FieldsConfig, FieldsConfig, bool) {
	for _, v := range configList {
		if v.ID == addConfigID {
			return configList, v, true
		}
	}
	appendData, err := GetConfigByID(&ArgsGetConfigByID{
		ID:    addConfigID,
		OrgID: -1,
	})
	if err != nil {
		return configList, FieldsConfig{}, false
	}
	configList = append(configList, appendData)
	return configList, appendData, true
}
