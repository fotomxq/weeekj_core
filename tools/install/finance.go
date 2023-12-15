package ToolsInstall

import (
	"errors"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
)

// InstallFinance 财务系统初始化
func InstallFinance() error {
	//配置文件名称
	configFileName := "finance.json"
	//检查配置文件是否存在
	if !checkConfigFile(configFileName) {
		return nil
	}
	//声明结构
	type dataInstallValueType struct {
		//标识码
		// 可用于同一类货币下，多个用途，如赠送的储值额度、或用户自行充值的额度
		// user 用户自己储值 ; deposit 押金 ; free 免费赠送额度 ; ... 特定系统下的充值模块
		Mark string `bson:"Mark" json:"mark"`
		//显示名称
		Name string `bson:"Name" json:"name"`
		//备注
		Des string `bson:"Des" json:"des"`
		//储蓄货币类型
		// 采用CoreCurrency匹配
		Currency int `bson:"Currency" json:"currency"`
		//能否取出
		// 如果能，则允许用户使用取出接口
		TakeOut bool `bson:"TakeOut" json:"takeOut"`
		//取款最低限额
		// 低于该资金禁止取款，同时需启动是否可取
		TakeLimit int64 `bson:"TakeLimit" json:"takeLimit"`
		//单次存款最低限额
		OnceSaveMinLimit int64 `bson:"OnceSaveMinLimit" json:"onceSaveMinLimit"`
		//单次存款最大限额
		OnceSaveMaxLimit int64 `bson:"OnceSaveMaxLimit" json:"onceSaveMaxLimit"`
		//单次取款最低限额
		OnceTakeMinLimit int64 `bson:"OnceTakeMinLimit" json:"onceTakeMinLimit"`
		//单次取款最大限额
		OnceTakeMaxLimit int64 `bson:"OnceTakeMaxLimit" json:"onceTakeMaxLimit"`
		//扩展参数设计
		Configs []CoreSQLConfig.FieldsConfigType `bson:"Configs" json:"configs"`
	}
	type dataInstallConfigType struct {
		//增量配置表
		FinanceDepositConfig []dataInstallValueType `json:"FinanceDepositConfig"`
		//财务默认Save.Mark
		ConfigFinancePayToDefaultSavingsMark string `json:"ConfigFinancePayToDefaultSavingsMark"`
	}
	//获取文件数据
	var data dataInstallConfigType
	if err := loadConfigFile(configFileName, &data); err != nil {
		return nil
	}
	//获取当前所有配置表
	dataList, _, err := FinanceDeposit.GetConfigList(&FinanceDeposit.ArgsGetConfigList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  9999,
			Sort: "id",
			Desc: false,
		},
	})
	//遍历数据，选择性添加
	for _, vJSON := range data.FinanceDepositConfig {
		isFind := false
		if err == nil {
			//检查是否已经创建，创建则跳出
			for _, vData := range dataList {
				if vData.Mark == vJSON.Mark {
					isFind = true
					break
				}
			}
		}
		//如果没找到，则添加
		if !isFind {
			if _, err2 := FinanceDeposit.SetConfig(&FinanceDeposit.ArgsSetConfig{
				Mark:             vJSON.Mark,
				Name:             vJSON.Name,
				Des:              vJSON.Des,
				Currency:         vJSON.Currency,
				TakeOut:          vJSON.TakeOut,
				TakeLimit:        vJSON.TakeLimit,
				OnceSaveMinLimit: vJSON.OnceSaveMinLimit,
				OnceSaveMaxLimit: vJSON.OnceSaveMaxLimit,
				OnceTakeMinLimit: vJSON.OnceTakeMinLimit,
				OnceTakeMaxLimit: vJSON.OnceTakeMaxLimit,
				Configs:          vJSON.Configs,
			}); err2 != nil {
				return err2
			}
		}
	}
	//设置全局收款Save.Mark
	FinancePayToDefaultSavingsMark, err := BaseConfig.GetByMark(&BaseConfig.ArgsGetByMark{
		Mark: "FinancePayToDefaultSavingsMark",
	})
	if err != nil {
		return err
	}
	if FinancePayToDefaultSavingsMark.Value == "" {
		configData, err := BaseConfig.GetByMark(&BaseConfig.ArgsGetByMark{
			Mark: "FinancePayToDefaultSavingsMark",
		})
		hash := ""
		if err == nil {
			hash = configData.UpdateHash
		}
		if err := BaseConfig.UpdateByMark(&BaseConfig.ArgsUpdateByMark{
			UpdateHash: hash,
			Mark:       "FinancePayToDefaultSavingsMark",
			Value:      data.ConfigFinancePayToDefaultSavingsMark,
		}); err != nil {
			return errors.New("set config by FinancePayToDefaultSavingsMark, " + err.Error())
		}
	}
	//反馈
	return nil
}
