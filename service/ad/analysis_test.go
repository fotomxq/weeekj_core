package ServiceAD

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	TestOrgArea "github.com/fotomxq/weeekj_core/v5/tools/test_org_area"
	"testing"
)

func TestInitAnalysis(t *testing.T) {
	TestInit(t)
	TestPutAD(t)
}

func TestGetAnalysis(t *testing.T) {
	dataList, err := GetAnalysis(&ArgsGetAnalysis{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubMinute().Time,
			MaxTime: CoreFilter.GetNowTimeCarbon().Time,
		},
		TimeType:  "hour",
		OrgID:     TestOrg.OrgData.ID,
		AreaID:    TestOrgArea.AreaData.ID,
		AdID:      adData.ID,
		IsHistory: false,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("analysis data: ", dataList)
	}
}

func TestClearAnalysis(t *testing.T) {
	TestClearPut(t)
	TestClear(t)
}
