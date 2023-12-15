package ToolsInstall

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseSMS "github.com/fotomxq/weeekj_core/v5/base/sms"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func InstallSMS() error {
	//配置文件名称
	configFileName := "sms_template.json"
	//检查配置文件是否存在
	if !checkConfigFile(configFileName) {
		return nil
	}
	//获取当前所有配置表
	dataList, _, err := BaseSMS.GetConfigList(&BaseSMS.ArgsGetConfigList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  100,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		IsRemove: false,
		Search:   "",
	})
	//如果存在数据，则跳出
	if err == nil && len(dataList) > 0 {
		return nil
	}
	//声明结构
	type dataInstallValueType struct {
		//组织ID
		OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
		//使用哪一家？
		// tencent / aliyun
		System string `db:"system" json:"system" check:"mark"`
		//来源系统的显示名称
		Name string `db:"name" json:"name" check:"name"`
		//应用ID
		AppID string `db:"app_id" json:"appID"`
		//应用密钥
		AppKey string `db:"app_key" json:"appKey"`
		//默认过期时间
		DefaultExpire string `db:"default_expire" json:"defaultExpire"`
		//获取间隔时间 秒
		TimeSpacing int64 `db:"time_spacing" json:"timeSpacing"`
		//模版ID
		TemplateID string `db:"template_id" json:"templateID"`
		//签名名称
		TemplateSign string `db:"template_sign" json:"templateSign"`
		//扩展参数
		TemplateParams CoreSQLConfig.FieldsConfigsType `db:"template_params" json:"templateParams"`
		//默认参数
		Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	}
	type dataInstallConfigType struct {
		//要建立的短信模版
		SMSConfig []dataInstallValueType `json:"smsConfig"`
		//对应的是SMSConfig的key值
		VerificationCodeSMSDefaultKey int `json:"VerificationCodeSMSDefaultKey"`
	}
	if Router2SystemConfig.Debug {
		//清理sms配置
		dataList2, _, err := BaseSMS.GetConfigList(&BaseSMS.ArgsGetConfigList{
			Pages: CoreSQLPages.ArgsDataList{
				Page: 1,
				Max:  1,
				Sort: "id",
				Desc: false,
			},
			OrgID:    -1,
			IsRemove: false,
			Search:   "",
		})
		if err == nil {
			for _, v := range dataList2 {
				if err := BaseSMS.DeleteConfig(&BaseSMS.ArgsDeleteConfig{
					ID:    v.ID,
					OrgID: -1,
				}); err != nil {
					return errors.New("delete config, " + err.Error())
				}
			}
		}
	}
	//获取文件数据
	var data dataInstallConfigType
	if err := loadConfigFile(configFileName, &data); err != nil {
		return nil
	}
	//遍历数据创建
	var newDataList []BaseSMS.FieldsConfig
	for _, v := range data.SMSConfig {
		if newData, err := BaseSMS.CreateConfig(&BaseSMS.ArgsCreateConfig{
			OrgID:          v.OrgID,
			System:         v.System,
			Name:           v.Name,
			AppID:          v.AppID,
			AppKey:         v.AppKey,
			DefaultExpire:  v.DefaultExpire,
			TimeSpacing:    v.TimeSpacing,
			TemplateID:     v.TemplateID,
			TemplateSign:   v.TemplateSign,
			TemplateParams: v.TemplateParams,
			Params:         v.Params,
		}); err == nil {
			newDataList = append(newDataList, newData)
		} else {
			if err.Error() == "data is exist" {
				continue
			}
			return errors.New("create sms config, " + err.Error())
		}
	}
	//修改默认的短信发送模版
	if len(newDataList) > 0 {
		for k, v := range newDataList {
			//如果是VerificationCodeSMSDefault
			// 用于短信验证码
			if k == data.VerificationCodeSMSDefaultKey {
				configData, err := BaseConfig.GetByMark(&BaseConfig.ArgsGetByMark{
					Mark: "VerificationCodeSMSDefault",
				})
				hash := ""
				if err == nil {
					hash = configData.UpdateHash
				}
				if err := BaseConfig.UpdateByMark(&BaseConfig.ArgsUpdateByMark{
					UpdateHash: hash,
					Mark:       "VerificationCodeSMSDefault",
					Value:      fmt.Sprint(v.ID),
				}); err != nil {
					return errors.New("set config is error, VerificationCodeSMSDefault, " + err.Error())
				}
				break
			}
		}
	}
	//反馈
	return nil
}
