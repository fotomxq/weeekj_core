package FinancePhysicalPay

import (
	"testing"

	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
)

func TestInitLog(t *testing.T) {
	TestInit(t)
	TestCreatePhysical(t)
}

func TestCreateLog(t *testing.T) {
	_, err := CreateLog(&ArgsCreateLog{
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
		UserID: TestOrg.UserInfo.ID,
		System: "order",
		Data: []ArgsCreateLogData{
			{
				PhysicalCount: 2,
				BindFrom: CoreSQLFrom.FieldsFrom{
					System: "mall",
					ID:     1,
					Mark:   "",
					Name:   "",
				},
				BindCount: 1,
			},
		},
	})
	ToolsTest.ReportError(t, err)
	//测试不可行
	_, err = CreateLog(&ArgsCreateLog{
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
		UserID: TestOrg.UserInfo.ID,
		System: "order",
		Data: []ArgsCreateLogData{
			{
				PhysicalCount: 1,
				BindFrom: CoreSQLFrom.FieldsFrom{
					System: "mall",
					ID:     1,
					Mark:   "",
					Name:   "",
				},
				BindCount: 1,
			},
		},
	})
	if err == nil {
		t.Error("create failed limit 2, ", err)
	} else {
		t.Log("create no failed, limit 2")
	}
	_, err = CreateLog(&ArgsCreateLog{
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
		UserID: TestOrg.UserInfo.ID,
		System: "order",
		Data: []ArgsCreateLogData{
			{
				PhysicalCount: 100,
				BindFrom: CoreSQLFrom.FieldsFrom{
					System: "mall",
					ID:     1,
					Mark:   "",
					Name:   "",
				},
				BindCount: 50,
			},
		},
	})
	if err == nil {
		t.Error("create failed limit 1, ", err)
	} else {
		t.Log("create no failed, limit 1")
	}
	//测试减少多个
	_, err = CreateLog(&ArgsCreateLog{
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
		UserID: TestOrg.UserInfo.ID,
		System: "order",
		Data: []ArgsCreateLogData{
			{
				PhysicalCount: 4,
				BindFrom: CoreSQLFrom.FieldsFrom{
					System: "mall",
					ID:     1,
					Mark:   "",
					Name:   "",
				},
				BindCount: 2,
			},
		},
	})
	ToolsTest.ReportError(t, err)
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages:       CoreSQLPages.ArgsDataList{Page: 1, Max: 10, Sort: "id", Desc: false},
		OrgID:       -1,
		BindID:      -1,
		UserID:      -1,
		System:      "",
		PhysicalID:  0,
		TimeBetween: CoreSQLTime.DataCoreTime{},
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestClearLog(t *testing.T) {
	TestDeletePhysical(t)
	TestClear(t)
}
