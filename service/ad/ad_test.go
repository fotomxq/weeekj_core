package ServiceAD

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	adData FieldsAD
)

func TestInitAD(t *testing.T) {
	TestInit(t)
}

func TestCreateAD(t *testing.T) {
	var err error
	adData, err = CreateAD(&ArgsCreateAD{
		OrgID:       TestOrg.OrgData.ID,
		Mark:        "test",
		Name:        "测试广告",
		Des:         "测试广告描述",
		CoverFileID: 0,
		DesFiles:    []int64{},
		Params:      []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, adData)
}

func TestGetADByID(t *testing.T) {
	data, err := GetADByID(&ArgsGetADByID{
		ID:    adData.ID,
		OrgID: 0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetADList(t *testing.T) {
	dataList, dataCount, err := GetADList(&ArgsGetADList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		Mark:     "",
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateAD(t *testing.T) {
	err := UpdateAD(&ArgsUpdateAD{
		ID:          adData.ID,
		OrgID:       TestOrg.OrgData.ID,
		Mark:        "test",
		Name:        adData.Name,
		Des:         adData.Des,
		CoverFileID: adData.CoverFileID,
		DesFiles:    adData.DesFiles,
		Params:      adData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteAD(t *testing.T) {
	err := DeleteAD(&ArgsDeleteAD{
		ID:    adData.ID,
		OrgID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearAD(t *testing.T) {
	TestClear(t)
}
