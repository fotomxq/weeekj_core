package EAMCore

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
	"time"
)

var (
	newCoreData FieldsEAM
)

func TestCoreInit(t *testing.T) {
	TestInit(t)
	TestSetProduct(t)
}

func TestCreateCore(t *testing.T) {
	newID, err := CreateCore(&ArgsCreateCore{
		Code:               CoreFilter.GetRandStr4(10),
		OrgID:              TestOrg.OrgData.ID,
		ProductID:          newERPProductData.ID,
		WarehouseBatchID:   0,
		ERPPurchaseOrderID: 0,
		Status:             0,
		Price:              10,
		WarrantyAt:         time.Time{},
		Location:           "测试地点",
		Remark:             "测试备注",
	})
	ToolsTest.ReportData(t, err, newID)
	newCoreData.ID = newID
}

func TestGetCore(t *testing.T) {
	data, err := GetCore(&ArgsGetCore{
		ID:    newCoreData.ID,
		OrgID: newCoreData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
	newCoreData = data
}

func TestGetCoreByCode(t *testing.T) {
	data, err := GetCoreByCode(&ArgsGetCoreByCode{
		Code:  newCoreData.Code,
		OrgID: newCoreData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
	newCoreData = data
}

func TestGetCoreList(t *testing.T) {
	dataList, dataCount, err := GetCoreList(&ArgsGetCoreList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:              -1,
		Code:               "",
		ProductID:          -1,
		WarehouseBatchID:   -1,
		ERPPurchaseOrderID: -1,
		Status:             -1,
		IsRemove:           false,
		Search:             "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateCore(t *testing.T) {
	err := UpdateCore(&ArgsUpdateCore{
		ID:         newCoreData.ID,
		Code:       newCoreData.Code,
		OrgID:      newCoreData.OrgID,
		Status:     newCoreData.Status,
		Price:      newCoreData.Price,
		WarrantyAt: newCoreData.WarrantyAt,
		Location:   newCoreData.Location,
		Remark:     newCoreData.Remark,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateCoreStatus(t *testing.T) {
	err := UpdateCoreStatus(&ArgsUpdateCoreStatus{
		ID:     newCoreData.ID,
		Code:   newCoreData.Code,
		OrgID:  newCoreData.OrgID,
		Status: newCoreData.Status,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteCore(t *testing.T) {
	err := DeleteCore(&ArgsDeleteCore{
		ID:    newCoreData.ID,
		OrgID: newCoreData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestCoreClear(t *testing.T) {
	TestDeleteProduct(t)
	TestClear(t)
}
