package BaseService

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newServiceData FieldsService
)

func TestServiceInit(t *testing.T) {
	TestInit(t)
}

func TestCreateService(t *testing.T) {
	err := setService(&argsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "测试服务",
		Description:  "测试服务描述",
		EventSubType: "server",
		Code:         "A001",
		EventType:    "nats",
		EventURL:     "/test/a1",
		EventParams:  "test",
	})
	ToolsTest.ReportError(t, err)
}
func TestGetServiceByCode(t *testing.T) {
	data := getServiceByCode("A001")
	if data.ID < 1 {
		t.Fatal("get service by code error")
		return
	}
	newServiceData = data
	t.Log(data)
}

func TestGetServiceByID(t *testing.T) {
	data := getServiceByID(newServiceData.ID)
	if data.ID < 1 {
		t.Fatal("get service by id error")
		return
	}
	t.Log(data)
}

func TestGetServiceList(t *testing.T) {
	dataList, dataCount, err := GetServiceList(&ArgsGetServiceList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		Code:     "",
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil {
		t.Log(dataList)
	}
}

func TestUpdateService(t *testing.T) {
	err := updateService(&argsUpdateService{
		ID:           newServiceData.ID,
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "测试服务",
		Description:  "测试服务描述",
		EventSubType: "server",
		Code:         "A001",
		EventType:    "nats",
		EventURL:     "/test/a1",
		EventParams:  "test",
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteService(t *testing.T) {
	err := deleteService(&argsDeleteService{
		ID: newServiceData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestServiceClear(t *testing.T) {
	TestClear(t)
}
