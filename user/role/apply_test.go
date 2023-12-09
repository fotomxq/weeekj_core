package UserRole

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newApplyData FieldsApply
)

func TestInitApply(t *testing.T) {
	TestInitType(t)
	TestCreateType(t)
}

func TestCreateApply(t *testing.T) {
	var err error
	newApplyData, err = CreateApply(&ArgsCreateApply{
		AuditDes:    "申请测试",
		RoleType:    newTypeData.ID,
		UserID:      1,
		Name:        "测试角色",
		Country:     86,
		City:        "10010",
		Gender:      1,
		Phone:       "17777777777",
		CoverFileID: 0,
		CertFiles:   []int64{},
		Params:      CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newApplyData)
}

func TestGetApplyID(t *testing.T) {
	data, err := GetApplyID(&ArgsGetApplyID{
		ID: newApplyData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetApplyList(t *testing.T) {
	dataList, dataCount, err := GetApplyList(&ArgsGetApplyList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		RoleType:    -1,
		UserID:      -1,
		NeedIsAudit: false,
		IsAudit:     false,
		IsRemove:    false,
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestAuditApply(t *testing.T) {
	errCode, err := AuditApply(&ArgsAuditApply{
		ID:          newApplyData.ID,
		IsAudit:     false,
		AuditBanDes: "拒绝测试",
	})
	ToolsTest.ReportError(t, err)
	if err != nil {
		t.Error("errCode: ", errCode)
	}
	TestDeleteApply(t)
	TestCreateApply(t)
	errCode, err = AuditApply(&ArgsAuditApply{
		ID:          newApplyData.ID,
		IsAudit:     true,
		AuditBanDes: "",
	})
	ToolsTest.ReportError(t, err)
	if err != nil {
		t.Error("errCode: ", errCode)
	}
}

func TestDeleteApply(t *testing.T) {

}

func TestClearApply(t *testing.T) {
	TestDeleteType(t)
}
