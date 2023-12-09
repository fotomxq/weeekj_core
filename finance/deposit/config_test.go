package FinanceDeposit

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	"testing"
)

func TestInit3(t *testing.T) {
	TestInit(t)
}

func TestSetConfig(t *testing.T) {
	data, err := SetConfig(&ArgsSetConfig{
		Mark:             configMark,
		Name:             "测试资金池",
		Des:              "测试备注",
		Currency:         86,
		TakeOut:          true,
		TakeLimit:        0,
		OnceSaveMinLimit: 0,
		OnceSaveMaxLimit: 0,
		OnceTakeMinLimit: 0,
		OnceTakeMaxLimit: 0,
		Configs:          nil,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetConfigList(t *testing.T) {
	data, count, err := GetConfigList(&ArgsGetConfigList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "mark",
			Desc: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, count)
	}
}

func TestGetConfigByMark(t *testing.T) {
	data, err := GetConfigByMark(&ArgsGetConfigByMark{
		Mark: configMark,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestDeleteConfigByMark(t *testing.T) {
	err := DeleteConfigByMark(&ArgsDeleteConfigByMark{
		Mark: configMark,
	})
	if err != nil {
		t.Error(err)
	}
}
