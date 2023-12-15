package ToolsInstall

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseEmail "github.com/fotomxq/weeekj_core/v5/base/email"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func InstallEmail() error {
	configFileName := "email_server.json"
	if !checkConfigFile(configFileName) {
		return nil
	}
	if Router2SystemConfig.Debug {
		//清理email
		dataList3, _, err := BaseEmail.GetServerList(&BaseEmail.ArgsGetServerList{
			Pages: CoreSQLPages.ArgsDataList{
				Page: 1,
				Max:  99999,
				Sort: "id",
				Desc: false,
			},
			OrgID:  -1,
			Search: "",
		})
		if err != nil {
			return err
		}
		for _, v := range dataList3 {
			if err := BaseEmail.DeleteServerByID(&BaseEmail.ArgsDeleteServerByID{
				ID:    v.ID,
				OrgID: -1,
			}); err != nil {
				return errors.New("delete service by id, " + err.Error())
			}
		}
	}
	//获取当前所有配置表
	dataList, _, err := BaseEmail.GetServerList(&BaseEmail.ArgsGetServerList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  99999,
			Sort: "id",
			Desc: false,
		},
		OrgID:  0,
		Search: "",
	})
	//如果存在内容，则跳出
	if err == nil && len(dataList) > 0 {
		return nil
	}
	//声明结构
	type dataInstallValueType struct {
		Host     string `json:"Host"`
		Port     string `json:"Port"`
		Email    string `json:"Email"`
		Password string `json:"Password"`
		IsSSL    bool   `json:"IsSSL"`
		Name     string `json:"Name"`
	}
	type dataInstallConfigType struct {
		//要建立的配置
		EmailServer []dataInstallValueType `json:"EmailServer"`
		//使用的默认邮箱配置
		EmailDefaultServerKey int `json:"EmailDefaultServerKey"`
	}
	//获取文件数据
	var data dataInstallConfigType
	if err := loadConfigFile(configFileName, &data); err != nil {
		return nil
	}
	//遍历数据
	var newDataList []BaseEmail.FieldsEmailServerType
	for _, v := range data.EmailServer {
		//创建新的
		if newData, err := BaseEmail.CreateServer(&BaseEmail.ArgsCreateServer{
			OrgID:    0,
			Name:     v.Name,
			Host:     v.Host,
			Port:     v.Port,
			IsSSL:    v.IsSSL,
			Email:    v.Email,
			Password: v.Password,
			Params:   nil,
		}); err == nil {
			newDataList = append(newDataList, newData)
		} else {
			if err.Error() == "data is exist" {
				continue
			}
			return errors.New("cannot create email config, " + err.Error())
		}
	}
	//修改默认的短信发送模版
	if len(newDataList) > 0 {
		for k, v := range newDataList {
			//如果是EmailDefaultServerID
			// 用于默认模版
			if k == data.EmailDefaultServerKey {
				configData, err := BaseConfig.GetByMark(&BaseConfig.ArgsGetByMark{
					Mark: "EmailDefaultServerID",
				})
				hash := ""
				if err == nil {
					hash = configData.UpdateHash
				}
				if err := BaseConfig.UpdateByMark(&BaseConfig.ArgsUpdateByMark{
					UpdateHash: hash,
					Mark:       "EmailDefaultServerID",
					Value:      fmt.Sprint(v.ID),
				}); err != nil {
					return errors.New("set config is error, EmailDefaultServerID, " + err.Error())
				}
				break
			}
		}
	}
	//反馈
	return nil
}
