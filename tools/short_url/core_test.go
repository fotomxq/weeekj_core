package ToolsShortURL

import (
	"testing"
	"time"

	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
)

var (
	isInit  = false
	newData FieldsShortURL
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
}

func TestCreate(t *testing.T) {
	var err error
	newData, err = Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Minute * 5),
		OrgID:    1,
		UserID:   2,
		IsPublic: true,
		Data:     "测试数据集合",
		Params:   nil,
	})
	ToolsTest.ReportData(t, err, newData)
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:  0,
		UserID: 0,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetByID(t *testing.T) {
	data, err := GetByID(&ArgsGetByID{
		ID:     newData.ID,
		OrgID:  0,
		UserID: 0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetByKey(t *testing.T) {
	data, err := GetByKey(&ArgsGetByKey{
		Key:    newData.Key,
		OrgID:  0,
		UserID: 0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetByKeyAndPrivate(t *testing.T) {
	var err error
	newData, err = Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Minute * 5),
		OrgID:    1,
		UserID:   2,
		IsPublic: false,
		Data:     "测试数据集合私人数据private",
		Params:   nil,
	})
	ToolsTest.ReportData(t, err, newData)
	data, err := GetByKey(&ArgsGetByKey{
		Key:    newData.Key,
		OrgID:  newData.OrgID,
		UserID: newData.UserID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteByID(t *testing.T) {
	err := DeleteByID(&ArgsDeleteByID{
		ID:     newData.ID,
		OrgID:  0,
		UserID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestClear(t *testing.T) {
	go Run()
	time.Sleep(time.Second * 3)
}
