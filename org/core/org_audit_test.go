package OrgCoreCore

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newOrgAuditData FieldsOrgAudit
)

func TestInitOrgAudit(t *testing.T) {
	TestInit(t)
	TestCreateUser(t)
}

func TestCreateOrgAudit(t *testing.T) {
	data, errCode, err := CreateOrgAudit(&ArgsCreateOrgAudit{
		UserID:     newUserInfo.ID,
		Key:        "",
		Name:       "测试名称",
		Des:        "测试描述",
		ParentID:   0,
		ParentFunc: []string{},
		OpenFunc:   []string{"all"},
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
		newOrgAuditData = data
	}
}

func TestGetOrgAuditList(t *testing.T) {
	dataList, dataCount, err := GetOrgAuditList(&ArgsGetOrgAuditList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		UserID:      -1,
		ParentID:    0,
		ParentFunc:  []string{},
		OpenFunc:    []string{"all"},
		NeedIsAudit: false,
		IsAudit:     false,
		IsBan:       false,
		AuditUserID: -1,
		IsRemove:    false,
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateOrgAuditPass(t *testing.T) {
	errCode, err := UpdateOrgAuditPass(&ArgsUpdateOrgAuditPass{
		ID:          newOrgAuditData.ID,
		ParentID:    0,
		AuditUserID: newOrgAuditData.UserID,
	})
	ToolsTest.ReportError(t, err)
	if err != nil {
		t.Error(errCode)
	}
}

func TestUpdateOrgAuditBan(t *testing.T) {
	err := UpdateOrgAuditBan(&ArgsUpdateOrgAuditBan{
		ID:          newOrgAuditData.ID,
		ParentID:    0,
		AuditUserID: newOrgAuditData.UserID,
		BanDes:      "测试拉黑",
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteOrgAudit(t *testing.T) {
	err := DeleteOrgAudit(&ArgsDeleteOrgAudit{
		ID:       newOrgAuditData.ID,
		ParentID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearOrgAudit(t *testing.T) {
	TestClear(t)
}
