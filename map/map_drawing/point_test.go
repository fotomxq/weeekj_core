package MapMapDrawing

import (
	CoreSQLGPS "gitee.com/weeekj/weeekj_core/v5/core/sql/gps"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newPointData FieldsPoint
)

func TestInitPoint(t *testing.T) {
	TestInit(t)
}

func TestCreatePoint(t *testing.T) {
	TestCreatePic(t)
	data, err := CreatePoint(&ArgsCreatePoint{
		OrgID:       1,
		PicID:       newPicData.ID,
		PicPoint:    CoreSQLGPS.FieldsPoint{},
		GPSPoint:    CoreSQLGPS.FieldsPoint{},
		Radius:      0,
		CoverFileID: 0,
		CoverIcon:   "",
		CoverRGB:    "",
		BindMark:    "pic",
		BindID:      newPicData.ID,
		Params:      nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newPointData = data
	}
}

func TestGetPointByPic(t *testing.T) {
	dataList, err := GetPointByPic(&ArgsGetPointByPic{
		PicID: newPicData.ID,
		OrgID: newPointData.OrgID,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestUpdatePoint(t *testing.T) {
	err := UpdatePoint(&ArgsUpdatePoint{
		ID:          newPointData.ID,
		OrgID:       newPointData.OrgID,
		PicPoint:    newPointData.PicPoint,
		GPSPoint:    newPointData.GPSPoint,
		Radius:      newPointData.Radius,
		CoverFileID: newPointData.CoverFileID,
		CoverIcon:   newPointData.CoverIcon,
		CoverRGB:    newPointData.CoverRGB,
		BindMark:    newPointData.BindMark,
		BindID:      newPointData.BindID,
		Params:      newPointData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeletePoint(t *testing.T) {
	err := DeletePoint(&ArgsDeletePoint{
		ID:    newPointData.ID,
		OrgID: newPointData.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestDeletePic(t)
}
