package ToolsHolidaySeason

import (
	"testing"
	"time"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
)

var (
	isInit  = false
	newData FieldsHolidaySeason
)

func TestInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	ToolsTest.Init(t)
}

func TestSet(t *testing.T) {
	err := Set(&ArgsSet{
		DateAt:    CoreFilter.GetNowTime(),
		Status:    0,
		IsHoliday: false,
		Name:      "星期二",
		Wage:      3,
		IsForce:   false,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestCheckIsWork(t *testing.T) {
	b := CheckIsWork(&ArgsCheckIsWork{
		DateAt: CoreFilter.GetNowTimeCarbon().Time,
	})
	t.Log("is work: ", b)
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		DateMin:  CoreFilter.GetNowTimeCarbon().AddDays(-1).Time,
		DateMax:  CoreFilter.GetNowTimeCarbon().AddDays(10).Time,
		HaveWork: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newData = dataList[0]
	}
}

func TestRun(t *testing.T) {
	TestInit(t)
	go Run()
	time.Sleep(time.Second * 5)
}

func TestDelete(t *testing.T) {
	err := Delete(&ArgsDelete{
		ID: newData.ID,
	})
	if err != nil {
		t.Error(err)
	}
}
