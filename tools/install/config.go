package ToolsInstall

import (
	"encoding/json"
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
)

func InstallConfig() error {
	//声明结构
	type dataInstallConfigValueType struct {
		Mark        string `json:"Mark"`
		AllowPublic bool   `json:"AllowPublic"`
		Name        string `json:"Name"`
		ValueType   string `json:"ValueType"`
		Value       string `json:"Value"`
		GroupMark   string `json:"GroupMark"`
		Des         string `json:"Des"`
		NeedDelete  bool   `json:"needDelete"`
	}
	type dataInstallConfigType struct {
		Configs []dataInstallConfigValueType
	}
	//所有配置
	var configAllConfig dataInstallConfigType
	//获取config目录下所有的json文件，并遍历处理
	fileList, err := CoreFile.GetFileList(fmt.Sprint(configDir, "config"), []string{"json"}, true)
	if err != nil {
		return nil
	}
	for _, vSrc := range fileList {
		//获取文件数据
		var data dataInstallConfigType
		fileByte, err := CoreFile.LoadFile(vSrc)
		if err != nil {
			return errors.New("load file, " + err.Error())
		}
		err = json.Unmarshal(fileByte, &data)
		if err != nil {
			return errors.New("get byte to json, " + err.Error())
		}
		for _, v := range data.Configs {
			if v.Mark == "" {
				continue
			}
			//写入一切配置，后续需要用到
			configAllConfig.Configs = append(configAllConfig.Configs, v)
			//检查是否已经存在，如果存在则跳过
			data, err := BaseConfig.GetByMark(&BaseConfig.ArgsGetByMark{
				Mark: v.Mark,
			})
			//检查是否为删除模式
			if v.NeedDelete {
				if err == nil && data.Mark != "" {
					err = BaseConfig.DeleteByMark(&BaseConfig.ArgsDeleteByMark{
						Mark: data.Mark,
					})
					if err != nil {
						return err
					}
				}
			}
			//如果存在配置
			if err == nil {
				//变更配置描述信息
				err = BaseConfig.UpdateInfo(&BaseConfig.ArgsUpdateInfo{
					Mark:        v.Mark,
					AllowPublic: v.AllowPublic,
					Name:        v.Name,
					GroupMark:   v.GroupMark,
					Des:         v.Des,
					ValueType:   v.ValueType,
				})
				if err != nil {
					return err
				}
				continue
			} else {
				//创建新的
				err = BaseConfig.Create(&BaseConfig.ArgsCreate{
					Mark:        v.Mark,
					AllowPublic: v.AllowPublic,
					Name:        v.Name,
					ValueType:   v.ValueType,
					Value:       v.Value,
					GroupMark:   v.GroupMark,
					Des:         v.Des,
				})
				if err != nil {
					return errors.New(fmt.Sprint("cannot create config, mark: ", v.Mark, ", err: "+err.Error()))
				}
			}
		}
		/**
		vSrcBase, err := CoreFile.GetFileInfo(vSrc)
		if err != nil{
			CoreLog.Error("安装系统配置，但无法获取配置文件数据: ", vSrc)
			continue
		}
		CoreLog.Info("安装系统配置，配置分组配置文件: ", vSrcBase.Name() ,"，设置了: ", len(data.Configs), "条数据。")
		*/
	}
	//去掉该设计，改为配置中会出现NeedDelete，用于删除版本迭代的旧配置
	//获取所有配置，反向判断配置是否应该被删除？
	//allConfig, err := BaseConfig.GetAll()
	//if err != nil {
	//	return errors.New("cannot get all config, err: " + err.Error())
	//}
	//deleteCount := 0
	//for _, v := range allConfig {
	//	isFind := false
	//	for _, v2 := range configAllConfig.Configs {
	//		if v.Mark == v2.Mark {
	//			isFind = true
	//			break
	//		}
	//	}
	//	if !isFind && !installAppend {
	//		err = BaseConfig.DeleteByMark(&BaseConfig.ArgsDeleteByMark{
	//			Mark: v.Mark,
	//		})
	//		if err != nil {
	//			return errors.New("config not in db, but cannot delete config by mark: " + v.Mark + ", err: " + err.Error())
	//		}
	//		deleteCount += 1
	//	}
	//}
	/**
	CoreLog.Info("安装系统配置，发现多余配置项，共计：", deleteCount, "条，已经删除。")
	*/
	//全部完成后，重新加载所有数据到缓冲
	if _, err := BaseConfig.GetAll(); err != nil {
		return errors.New("get all config, " + err.Error())
	}
	//反馈
	return nil
}
