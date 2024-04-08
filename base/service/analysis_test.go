package BaseService

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestAnalysisInit(t *testing.T) {
	TestServiceInit(t)
}

func TestAnalysisCreate(t *testing.T) {
	err := appendAnalysisData(&argsAppendAnalysisData{
		ServiceID:    newServiceData.ID,
		SendCount:    1,
		ReceiveCount: 1,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetAnalysisList(t *testing.T) {
	nowAt := CoreFilter.GetNowTimeCarbon()
	minAt := nowAt.StartOfMonth()
	maxAt := nowAt.EndOfMonth()
	dataList, dataCount, err := GetAnalysisList(&ArgsGetAnalysisList{
		ServiceID: newServiceData.ID,
		BetweenAt: CoreSQL2.ArgsTimeBetween{
			MinTime: CoreFilter.GetTimeToDefaultTime(minAt.Time),
			MaxTime: CoreFilter.GetTimeToDefaultTime(maxAt.Time),
		},
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)

}

func TestAnalysisClear(t *testing.T) {
	TestServiceClear(t)
}
