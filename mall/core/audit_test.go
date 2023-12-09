package MallCore

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newAuditData FieldsAudit
)

func TestInitAudit(t *testing.T) {
	TestInitProduct(t)
	TestCreateProduct(t)
}

func TestCreateAudit(t *testing.T) {
	var err error
	newAuditData, err = CreateAudit(&ArgsCreateAudit{
		OrgID:     newProductData.OrgID,
		ProductID: newProductData.ID,
	})
	ToolsTest.ReportData(t, err, newAuditData)
}

func TestGetAuditByProduct(t *testing.T) {
	data, err := GetAuditByProduct(&ArgsGetAuditByProduct{
		OrgID:     newProductData.OrgID,
		ProductID: newProductData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAuditByProducts(t *testing.T) {
	data, err := GetAuditByProducts(&ArgsGetAuditByProducts{
		OrgID:      newProductData.OrgID,
		ProductIDs: []int64{newProductData.ID},
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAuditList(t *testing.T) {
	dataList, dataCount, err := GetAuditList(&ArgsGetAuditList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:         -1,
		ProductID:     0,
		NeedIsExpire:  false,
		IsExpire:      false,
		NeedIsPassing: false,
		IsPassing:     false,
		IsRemove:      false,
		Search:        "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateAuditPassing(t *testing.T) {
	err := UpdateAuditPassing(&ArgsUpdateAuditPassing{
		ID: newAuditData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateAuditBan(t *testing.T) {
	err := UpdateAuditBan(&ArgsUpdateAuditBan{
		ID:          newAuditData.ID,
		BanDes:      "测试拒绝",
		BanDesFiles: []int64{},
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteAudit(t *testing.T) {
	err := DeleteAudit(&ArgsDeleteAudit{
		ID:    newAuditData.ID,
		OrgID: newAuditData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearAudit(t *testing.T) {
	//TestDeleteProduct(t)
	TestClearProduct(t)
}
