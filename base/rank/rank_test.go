package BaseRank

import (
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"testing"
	"time"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
)

var (
	rankData FieldsRank
)

func TestInit2(t *testing.T) {
	TestInit(t)
}

func TestAppendRank(t *testing.T) {
	var err error
	rankData, err = AppendRank(&ArgsAppendRank{
		ServiceMark: "a1",
		ExpireAt:    CoreFilter.GetNowTime().Add(time.Second * 30),
		PickMin:     10,
		MissionMark: "b1",
		MissionData: []byte("abc01,1600924914,1601014914,false"),
	})
	ToolsTest.ReportData(t, err, rankData)
}

func TestGetRankList(t *testing.T) {
	dataList, dataCount, err := GetRankList(&ArgsGetRankList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		ServiceMark: "a1",
		MissionMark: "b1",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestPickRank(t *testing.T) {
	dataList, err := PickRank(&ArgsPickRank{
		ServiceMark: "a1",
		MissionMark: "b1",
		Max:         5,
	})
	ToolsTest.ReportDataList(t, err, dataList, int64(len(dataList)))
}

func TestOverRank(t *testing.T) {
	err := OverRank(&ArgsOverRank{
		ID:     rankData.ID,
		Result: []byte("ok"),
	})
	ToolsTest.ReportError(t, err)
}

func TestGetRankOverList(t *testing.T) {
	dataList, dataCount, err := GetRankOverList(&ArgsGetRankOverList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		ServiceMark: "a1",
		MissionMark: "b1",
		MissionData: []byte("abc01,1600924914,1601014914,false"),
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	//检查不能存在同样的数据
	dataList2, err := PickRank(&ArgsPickRank{
		ServiceMark: "a1",
		MissionMark: "b1",
		Max:         5,
	})
	ToolsTest.ReportDataList(t, err, dataList2, int64(len(dataList2)))
	if len(dataList2) > 0 {
		t.Error("dataList2 > 0")
	}
}

// 清理
func TestRemoveAll(t *testing.T) {
	t.Skip()
	_, err := CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_rank", "true", nil)
	ToolsTest.ReportError(t, err)
}
