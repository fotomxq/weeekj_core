package OrgCoreCore

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newBindAuditData FieldsBindAudit
)

func TestInitBindAudit(t *testing.T) {
	TestInit(t)
	TestCreateUser(t)
	TestCreateOrg(t)
}

func TestCreateBindAudit(t *testing.T) {
	data, errCode, err := CreateBindAudit(&ArgsCreateBindAudit{
		UserID:   newUserInfo.ID,
		Name:     "测试名称",
		OrgID:    orgData.ID,
		GroupIDs: []int64{},
		Manager:  []string{"all"},
		Params:   nil,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
		newBindAuditData = data
	}
}

func TestGetBindAuditList(t *testing.T) {
	dataList, dataCount, err := GetBindAuditList(&ArgsGetBindAuditList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:       newBindAuditData.OrgID,
		UserID:      -1,
		GroupID:     -1,
		Manager:     "",
		NeedIsAudit: false,
		IsAudit:     false,
		IsBan:       false,
		AuditBindID: -1,
		IsRemove:    false,
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateBindAuditPass(t *testing.T) {
	err := UpdateBindAuditPass(&ArgsUpdateBindAuditPass{
		ID:          newBindAuditData.ID,
		OrgID:       0,
		AuditBindID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateBindAuditBan(t *testing.T) {
	err := UpdateBindAuditBan(&ArgsUpdateBindAuditBan{
		ID:          newBindAuditData.ID,
		OrgID:       0,
		AuditBindID: 0,
		BanDes:      "测试拉黑",
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteBindAudit(t *testing.T) {
	err := DeleteBindAudit(&ArgsDeleteBindAudit{
		ID:    newBindAuditData.ID,
		OrgID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearBindAudit(t *testing.T) {
	TestDeleteOrg(t)
	TestClear(t)
}
