package ServiceAD

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newApplyData FieldsApply
)

func TestInitApply(t *testing.T) {
	TestInit(t)
}

func TestCreateApply(t *testing.T) {
	data, err := CreateApply(&ArgsCreateApply{
		StartAt:     CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().Time),
		EndAt:       CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().AddHour().Time),
		OrgID:       1,
		UserID:      1,
		AreaIDs:     []int64{},
		Mark:        "test",
		Name:        "测试名称",
		Des:         "测试内容",
		CoverFileID: 0,
		DesFiles:    []int64{},
		Params:      CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newApplyData = data
	}
}

func TestGetApplyList(t *testing.T) {
	dataList, dataCount, err := GetApplyList(&ArgsGetApplyList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:       -1,
		UserID:      -1,
		NeedIsAudit: false,
		IsAudit:     false,
		NeedIsStart: false,
		IsStart:     false,
		NeedIsEnd:   false,
		IsEnd:       false,
		Mark:        "",
		IsRemove:    false,
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetApplyID(t *testing.T) {
	data, err := GetApplyID(&ArgsGetApplyID{
		ID:     newApplyData.ID,
		OrgID:  newApplyData.OrgID,
		UserID: newApplyData.UserID,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newApplyData = data
	}
}

func TestAuditApply(t *testing.T) {
	err := AuditApply(&ArgsAuditApply{
		ID:          newApplyData.ID,
		OrgID:       newApplyData.OrgID,
		IsAudit:     true,
		AuditDes:    "审核通过",
		AuditBanDes: "",
		AreaIDs:     newApplyData.AreaIDs,
		Mark:        newApplyData.Mark,
	})
	ToolsTest.ReportError(t, err)
}

func TestAppendAnalysisClick(t *testing.T) {
	adData, err := getAdByAuto(newApplyData.OrgID, newApplyData.Mark)
	ToolsTest.ReportError(t, err)
	//检查次数是否非0
	TestGetApplyID(t)
	if newApplyData.Count < 1 {
		t.Error("analysis count less 1")
	} else {
		t.Log("analysis count pass")
	}
	//更新点击次数
	err = AppendAnalysisClick(&ArgsAppendAnalysisClick{
		OrgID:      newApplyData.OrgID,
		AreaID:     0,
		AdID:       adData.ID,
		ClickCount: 1,
	})
	ToolsTest.ReportError(t, err)
	//检查点击次数
	TestGetApplyID(t)
	if newApplyData.ClickCount < 1 {
		t.Error("analysis click count less 1")
	} else {
		t.Log("analysis click count pass")
	}
}

func TestDeleteApply(t *testing.T) {
	err := DeleteApply(&ArgsDeleteApply{
		ID:     newApplyData.ID,
		OrgID:  newApplyData.OrgID,
		UserID: newApplyData.UserID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearApply(t *testing.T) {
	TestClear(t)
}
