package UserMessage

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	newData       FieldsMessage
	sendUserID    int64 = 123
	receiveUserID int64 = 234
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
}

func TestCreate(t *testing.T) {
	var err error
	newData, err = Create(&ArgsCreate{
		WaitSendAt:    CoreFilter.GetNowTime(),
		SendUserID:    sendUserID,
		ReceiveUserID: receiveUserID,
		Title:         "测试标题1",
		Content:       "测试内容X",
		Files:         []int64{},
		Params:        CoreSQLConfig.FieldsConfigsType{},
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
		SendUserID:      0,
		WaitSendAt:      time.Time{},
		ReceiveUserID:   0,
		ReceiveReadAt:   time.Time{},
		ReceiveDeleteAt: time.Time{},
		IsRemove:        false,
		Search:          "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetByID(t *testing.T) {
	data, err := GetByID(&ArgsGetByID{
		ID:            newData.ID,
		SendUserID:    0,
		ReceiveUserID: 0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateByID(t *testing.T) {
	err := UpdateByID(&ArgsUpdateByID{
		ID:            newData.ID,
		SendUserID:    0,
		ReceiveUserID: 0,
		WaitSendAt:    newData.WaitSendAt,
		Title:         newData.Title,
		Content:       newData.Content,
		Files:         newData.Files,
		Params:        newData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdatePost(t *testing.T) {
	err := UpdatePost(&ArgsUpdatePost{
		ID:         newData.ID,
		SendUserID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateAudit(t *testing.T) {
	err := UpdateAudit(&ArgsUpdateAudit{
		ID: newData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateReceiveRead(t *testing.T) {
	t.Skip()
	err := UpdateReceiveRead(&ArgsUpdateReceiveRead{
		ID:            newData.ID,
		ReceiveUserID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateReceiveReads(t *testing.T) {
	err := UpdateReceiveReads(&ArgsUpdateReceiveReads{
		IDs:           []int64{newData.ID},
		ReceiveUserID: newData.ReceiveUserID,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteBySend(t *testing.T) {
	err := DeleteBySend(&ArgsDeleteBySend{
		IDs:        []int64{newData.ID},
		SendUserID: -1,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteByReceive(t *testing.T) {
	err := DeleteByReceive(&ArgsDeleteByReceive{
		IDs:           []int64{newData.ID},
		ReceiveUserID: 0,
	})
	ToolsTest.ReportError(t, err)
}
