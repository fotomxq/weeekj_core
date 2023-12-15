package TMSTransport

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newBindData FieldsBind
)

func TestBindInit(t *testing.T) {
	TestInit(t)
}

func TestSetBind(t *testing.T) {
	var err error
	newBindData, err = SetBind(&ArgsSetBind{
		OrgID:          TestOrg.OrgData.ID,
		BindID:         TestOrg.BindData.ID,
		MapAreaID:      newMapArea.ID,
		MoreMapAreaIDs: []int64{},
		Params:         []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newBindData)
	t.Log("new bind data, bind id: ", newBindData.BindID)
}

func TestGetBindList(t *testing.T) {
	dataList, dataCount, err := GetBindList(&ArgsGetBindList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:     -1,
		BindID:    -1,
		MapAreaID: -1,
		IsRemove:  false,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetBindByBindID(t *testing.T) {
	data, err := GetBindByBindID(&ArgsGetBindByBindID{
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetBind(t *testing.T) {
	data, err := GetBind(&ArgsGetBind{
		ID:    newBindData.ID,
		OrgID: newBindData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteBind(t *testing.T) {
	err := DeleteBind(&ArgsDeleteBind{
		ID:    newBindData.ID,
		OrgID: newBindData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestBindClear(t *testing.T) {
	TestClear(t)
}
