package MapMapDrawing

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newPicData FieldsPic
)

func TestInitPic(t *testing.T) {
	TestInit(t)
}

func TestCreatePic(t *testing.T) {
	data, err := CreatePic(&ArgsCreatePic{
		OrgID:      1,
		ParentID:   0,
		Name:       "测试模版",
		Des:        "test des",
		FileID:     0,
		FixHeight:  0,
		FixWidth:   0,
		ButtonName: "",
		BindAreaID: 0,
		Params:     nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newPicData = data
	}
}

func TestGetPicList(t *testing.T) {
	dataList, dataCount, err := GetPicList(&ArgsGetPicList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		ParentID: -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetPicChild(t *testing.T) {
	dataList, err := GetPicChild(&ArgsGetPicChild{
		ParentID: newPicData.ParentID,
		OrgID:    newPicData.OrgID,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestUpdatePic(t *testing.T) {
	err := UpdatePic(&ArgsUpdatePic{
		ID:         newPicData.ID,
		OrgID:      newPicData.OrgID,
		Name:       newPicData.Name,
		Des:        newPicData.Des,
		FileID:     newPicData.FileID,
		FixHeight:  newPicData.FixHeight,
		FixWidth:   newPicData.FixWidth,
		ButtonName: newPicData.ButtonName,
		BindAreaID: newPicData.BindAreaID,
		Params:     newPicData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeletePic(t *testing.T) {
	err := DeletePic(&ArgsDeletePic{
		ID:    newPicData.ID,
		OrgID: newPicData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}
