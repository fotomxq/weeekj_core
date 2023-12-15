package FinanceDeposit

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	"testing"
)

func TestInit2(t *testing.T) {
	TestInit(t)
}

func TestSetConfig2(t *testing.T) {
	TestSetConfig(t)
}

func TestSetByFrom(t *testing.T) {
	data, errCode, err := SetByFrom(&ArgsSetByFrom{
		UpdateHash: "",
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     1,
			Mark:   "",
			Name:   "测试用户",
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     2,
			Mark:   "",
			Name:   "测试组织",
		},
		ConfigMark:      configMark,
		AppendSavePrice: 1322,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
		depositData = data
	}
}

func TestSetByFrom2(t *testing.T) {
	data, errCode, err := SetByFrom(
		&ArgsSetByFrom{
			UpdateHash: depositData.UpdateHash,
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     depositData.CreateInfo.ID,
				Mark:   "",
				Name:   "测试用户",
			},
			FromInfo: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     3,
				Mark:   "",
				Name:   "测试组织",
			},
			ConfigMark:      configMark,
			AppendSavePrice: 1322,
		},
	)
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
	}
}

// 测试空资金账户建立
func TestSetByFrom3(t *testing.T) {
	data, errCode, err := SetByFrom(&ArgsSetByFrom{
		UpdateHash: "",
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     4,
			Mark:   "",
			Name:   "测试用户",
		},
		FromInfo:        CoreSQLFrom.FieldsFrom{},
		ConfigMark:      configMark,
		AppendSavePrice: 2,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
	}
}

func TestGetByID(t *testing.T) {
	data, err := GetByID(&ArgsGetByID{
		ID: depositData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetByFrom(t *testing.T) {
	data, err := GetByFrom(&ArgsGetByFrom{
		CreateInfo: depositData.CreateInfo,
		FromInfo:   depositData.FromInfo,
		ConfigMark: depositData.ConfigMark,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetList(t *testing.T) {
	data, count, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     depositData.CreateInfo.ID,
			Mark:   "",
			Name:   "",
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     depositData.FromInfo.ID,
			Mark:   "",
			Name:   "",
		},
		ConfigMark: "",
		MinPrice:   0,
		MaxPrice:   0,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, count)
		if count < 1 {
			t.Error("data: ", data, ", err: ", err)
		}
	}
}

func TestDeleteConfigByMark2(t *testing.T) {
	TestDeleteConfigByMark(t)
}

func TestClear(t *testing.T) {
}
