package ServiceDistribution

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newDistributionData FieldsDistribution
)

func TestInitDistribution(t *testing.T) {
	TestInit(t)
}

func TestCreateDistribution(t *testing.T) {
	var err error
	newDistributionData, err = CreateDistribution(&ArgsCreateDistribution{
		OrgID:  TestOrg.OrgData.ID,
		Name:   "测试商户",
		UserID: TestOrg.UserInfo.ID,
	})
	ToolsTest.ReportData(t, err, newDistributionData)
}

func TestGetDistributionList(t *testing.T) {
	dataList, dataCount, err := GetDistributionList(&ArgsGetDistributionList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:    -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateDistribution(t *testing.T) {
	err := UpdateDistribution(&ArgsUpdateDistribution{
		ID:     newDistributionData.ID,
		OrgID:  newDistributionData.OrgID,
		Name:   newDistributionData.Name,
		UserID: newDistributionData.UserID,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteDistribution(t *testing.T) {
	err := DeleteDistribution(&ArgsDeleteDistribution{
		ID:    newDistributionData.ID,
		OrgID: newDistributionData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearDistribution(t *testing.T) {
	TestClear(t)
}
