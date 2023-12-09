package FinanceAssets

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	"testing"
)

func TestInit4(t *testing.T) {
	TestInit5(t)
	TestSetAssets(t)
}

func TestGetLogList(t *testing.T) {
	data, count, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:     0,
		BindID:    0,
		UserID:    0,
		ProductID: 0,
		IsHistory: false,
		Search:    "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, count)
	}
}

func TestClear4(t *testing.T) {
	TestClearAssetsByProductID(t)
	TestClear5(t)
}
