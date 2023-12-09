package OrgCert

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newCertData FieldsCert
)

func TestInitCert(t *testing.T) {
	TestInitConfig(t)
	TestCreateConfig(t)
}

func TestCreateCert(t *testing.T) {
	data, errCode, err := CreateCert(&ArgsCreateCert{
		OrgID:      TestOrg.OrgData.ID,
		ChildOrgID: 0,
		BindFrom:   newConfigData.BindFrom,
		ConfigID:   newConfigData.ID,
		ConfigMark: "",
		BindID:     1,
		Name:       "测试",
		SN:         "SNSN",
		FileIDs:    []int64{},
		ExpireAt:   "",
		Params:     nil,
	})
	if err != nil {
		t.Error(err, errCode)
		return
	} else {
		t.Log(data)
		newCertData = data
	}
	data, errCode, err = CreateCert(&ArgsCreateCert{
		OrgID:      TestOrg.OrgData.ID,
		ChildOrgID: 0,
		BindFrom:   newConfigData.BindFrom,
		ConfigID:   newConfigData.ID,
		ConfigMark: "",
		BindID:     1,
		Name:       "测试",
		SN:         "SNSN",
		FileIDs:    []int64{},
		ExpireAt:   "",
		Params:     nil,
	})
	if err != nil {
		t.Error(err, errCode)
		return
	} else {
		t.Log(data)
		newCertData = data
	}
}

func TestGetCertList(t *testing.T) {
	dataList, dataCount, err := GetCertList(&ArgsGetCertList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:        -1,
		ChildOrgID:   -1,
		ConfigID:     -1,
		AuditBindID:  -1,
		BindID:       -1,
		NeedIsExpire: false,
		IsExpire:     false,
		NeedIsAudit:  false,
		IsAudit:      false,
		NeedIsPay:    false,
		IsPay:        false,
		BindFrom:     "",
		Mark:         "",
		IsRemove:     false,
		Search:       "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetCertChildGroupList(t *testing.T) {
	dataList, dataCount, err := GetCertChildGroupList(&ArgsGetCertChildGroupList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "child_org_id",
			Desc: false,
		},
		OrgID:        -1,
		ConfigID:     -1,
		NeedIsExpire: false,
		IsExpire:     false,
		NeedIsAudit:  false,
		IsAudit:      false,
		NeedIsPay:    false,
		IsPay:        false,
		BindFrom:     "",
		Mark:         "",
		IsRemove:     false,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetCert(t *testing.T) {
	data, err := GetCert(&ArgsGetCert{
		ID:         newCertData.ID,
		OrgID:      newCertData.OrgID,
		ChildOrgID: newCertData.ChildOrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateCert(t *testing.T) {
	err := UpdateCert(&ArgsUpdateCert{
		ID:         newCertData.ID,
		OrgID:      newCertData.OrgID,
		ChildOrgID: newCertData.ChildOrgID,
		BindID:     newCertData.BindID,
		Name:       newCertData.Name,
		SN:         newCertData.SN,
		ExpireAt:   "",
		FileIDs:    newCertData.FileIDs,
		Params:     newCertData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateCertAudit(t *testing.T) {
	err := UpdateCertAudit(&ArgsUpdateCertAudit{
		ID:          newCertData.ID,
		OrgID:       newCertData.OrgID,
		ChildOrgID:  newCertData.ChildOrgID,
		AuditBindID: TestOrg.BindData.ID,
		IsBan:       false,
		AuditDes:    "测试审核",
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteCert(t *testing.T) {
	err := DeleteCert(&ArgsDeleteCert{
		ID:         newCertData.ID,
		OrgID:      newCertData.OrgID,
		ChildOrgID: newCertData.ChildOrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearCert(t *testing.T) {
	TestDeleteConfig(t)
	TestClearConfig(t)
}
