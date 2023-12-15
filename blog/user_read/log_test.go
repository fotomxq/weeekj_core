package BlogUserRead

import (
	BlogCore "github.com/fotomxq/weeekj_core/v5/blog/core"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"github.com/golang-module/carbon"
	"testing"
)

var (
	newContentData BlogCore.FieldsContent
)

func TestInitLog(t *testing.T) {
	TestInit(t)
	newKey, err := CoreFilter.GetRandStr3(10)
	if err != nil {
		t.Error(err)
		return
	}
	newContentData, err = BlogCore.CreateContent(&BlogCore.ArgsCreateContent{
		OrgID:       TestOrg.OrgData.ID,
		Key:         newKey,
		IsTop:       false,
		SortID:      0,
		Tags:        []int64{},
		Title:       "测试标题",
		CoverFileID: 0,
		Des:         "测试内容姑娘",
		Params:      CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newContentData)
	if err != nil {
		return
	}
	err = BlogCore.UpdatePublish(&BlogCore.ArgsUpdatePublish{
		ID:    newContentData.ID,
		OrgID: -1,
	})
	ToolsTest.ReportData(t, err, newContentData)
	if err != nil {
		return
	}
}

func TestCreateLog(t *testing.T) {
	err := CreateLog(&ArgsCreateLog{
		UserID:    TestOrg.UserInfo.ID,
		FromMark:  "web",
		FromName:  "web_002",
		Name:      "测试姓名2",
		IP:        "0.0.0.1",
		ContentID: newContentData.ID,
		CreateAt:  CoreFilter.GetISOByTime(carbon.Now().SubHours(1).Time),
		LeaveAt:   "",
	})
	ToolsTest.ReportError(t, err)
	err = CreateLog(&ArgsCreateLog{
		UserID:    TestOrg.UserInfo.ID,
		FromMark:  "web",
		FromName:  "web_002",
		Name:      "测试姓名2",
		IP:        "0.0.0.1",
		ContentID: newContentData.ID,
		CreateAt:  CoreFilter.GetISOByTime(carbon.Now().SubHours(1).Time),
		LeaveAt:   CoreFilter.GetISOByTime(carbon.Now().SubMinutes(1).Time),
	})
	ToolsTest.ReportError(t, err)
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:       -1,
		ChildOrgID:  -1,
		UserID:      -1,
		FromMark:    "",
		FromName:    "",
		IP:          "",
		ContentID:   -1,
		SortID:      -1,
		ReadTimeMin: -1,
		ReadTimeMax: -1,
		TimeBetween: CoreSQLTime.DataCoreTime{},
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetLogCount(t *testing.T) {
	count, err := GetLogCount(&ArgsGetLogCount{
		OrgID:      -1,
		ChildOrgID: -1,
		ContentIDs: []int64{newContentData.ID},
		TimeBetween: CoreSQLTime.DataCoreTime{
			MinTime: CoreFilter.GetISOByTime(carbon.Now().SubHours(3).Time),
			MaxTime: CoreFilter.GetISOByTime(carbon.Now().AddSeconds(1).Time),
		},
		SkipTime: false,
	})
	ToolsTest.ReportData(t, err, count)
}

func TestClearLog(t *testing.T) {
	err := BlogCore.DeleteContent(&BlogCore.ArgsDeleteContent{
		ID:    newContentData.ID,
		OrgID: -1,
	})
	ToolsTest.ReportError(t, err)
	TestClear(t)
}
