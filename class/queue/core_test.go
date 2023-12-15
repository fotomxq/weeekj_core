package ClassQueue

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	queue     Queue
	queueData FieldsQueue
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
	queue.TableName = "test_queue"
}

func TestQueue_Append(t *testing.T) {
	err := queue.Append(&ArgsAppend{
		ModID:  123,
		Status: 0,
		Params: []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportError(t, err)
}

func TestQueue_GetList(t *testing.T) {
	dataList, dataCount, err := queue.GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		Status: -1,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataList, dataCount)
		queueData = dataList[0]
	}
}

func TestQueue_GetByModID(t *testing.T) {
	data, err := queue.GetByModID(&ArgsGetByModID{
		ModID: queueData.ModID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestQueue_GetByID(t *testing.T) {
	data, err := queue.GetByID(&ArgsGetByID{
		ID: queueData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestQueue_Pick(t *testing.T) {
	data, err := queue.Pick()
	ToolsTest.ReportData(t, err, data)
}

func TestQueue_UpdateStatus(t *testing.T) {
	TestQueue_Append(t)
	TestQueue_GetList(t)
	err := queue.UpdateStatus(&ArgsUpdateStatus{
		ID:     queueData.ID,
		Status: 1,
	})
	ToolsTest.ReportError(t, err)
}

func TestQueue_Delete(t *testing.T) {
	err := queue.Delete(&ArgsDelete{
		ID: queueData.ID,
	})
	ToolsTest.ReportError(t, err)
}
