package OrgCert

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newWarningData FieldsWarning
)

func TestInitWarning(t *testing.T) {
	TestInitCert(t)
	TestCreateCert(t)
}

func TestCreateWarning(t *testing.T) {
	err := createWarning(&argsCreateWarning{
		OrgID:      TestOrg.OrgData.ID,
		ChildOrgID: 0,
		CertID:     newCertData.ID,
		Msg:        "测试异常",
		Params:     nil,
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetWarningList(t *testing.T) {
	dataList, dataCount, err := GetWarningList(&ArgsGetWarningList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:        -1,
		ChildOrgID:   -1,
		NeedIsFinish: false,
		IsFinish:     false,
		TimeBetween:  CoreSQLTime.DataCoreTime{},
		Search:       "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil {
		newWarningData = dataList[0]
	}
}

func TestUpdateWarningFinish(t *testing.T) {
	err := UpdateWarningFinish(&ArgsUpdateWarningFinish{
		ID:         newWarningData.ID,
		OrgID:      newWarningData.OrgID,
		ChildOrgID: newWarningData.ChildOrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearWarning(t *testing.T) {
	TestDeleteCert(t)
	TestClearCert(t)
}
