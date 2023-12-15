package ToolsInstall

import (
	"errors"
	BaseEarlyWarning "github.com/fotomxq/weeekj_core/v5/base/early_warning"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func InstallEarlyWarning() error {
	//配置文件名称
	configFileName := "early_warning.json"
	//检查配置文件是否存在
	if !checkConfigFile(configFileName) {
		return nil
	}
	//声明结构
	type dataInstallEarlyWarningValueType struct {
		Mark       string   `json:"Mark"`
		Name       string   `json:"Name"`
		ExpireTime string   `json:"ExpireTime"`
		Title      string   `json:"Title"`
		Content    string   `json:"Content"`
		TemplateID string   `json:"TemplateID"`
		BindData   []string `json:"BindData"`
	}
	type dataInstallEarlyWarningVType struct {
		Template []dataInstallEarlyWarningValueType `json:"Template"`
	}
	//获取文件数据
	var data dataInstallEarlyWarningVType
	if err := loadConfigFile(configFileName, &data); err != nil {
		return nil
	}
	for _, v := range data.Template {
		if v.Mark == "" {
			continue
		}
		//检查是否已经存在，如果存在则跳过
		data, err := BaseEarlyWarning.GetTemplateByMark(&BaseEarlyWarning.ArgsGetTemplateByMark{
			Mark: v.Mark,
		})
		if err == nil {
			if Router2SystemConfig.Debug {
				err = BaseEarlyWarning.DeleteTemplate(&BaseEarlyWarning.ArgsDeleteTemplate{
					ID: data.ID,
				})
				if err != nil {
					return errors.New("delete old template, " + err.Error())
				}
			} else {
				continue
			}
		}
		//创建新的
		_, err = BaseEarlyWarning.CreateTemplate(&BaseEarlyWarning.ArgsCreateTemplate{
			Mark:              v.Mark,
			Name:              v.Name,
			DefaultExpireTime: v.ExpireTime,
			Title:             v.Title,
			Content:           v.Content,
			TemplateID:        v.TemplateID,
			BindData:          v.BindData,
		})
		if err != nil {
			return errors.New("cannot create template, " + err.Error())
		}
	}
	return nil
}
