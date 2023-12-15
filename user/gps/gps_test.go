package UserGPS

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitGPS(t *testing.T) {
	TestInit(t)
}

func TestCreate(t *testing.T) {
	err := Create(&ArgsCreate{
		UserID:    userID,
		Country:   86,
		City:      10010,
		MapType:   0,
		Longitude: 1.2,
		Latitude:  1.3,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		UserID:  0,
		Country: 0,
		City:    0,
		MapType: 0,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetLast(t *testing.T) {
	data, err := GetLast(&ArgsGetLast{
		UserID: userID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteByUser(t *testing.T) {
	err := DeleteByUser(&ArgsDeleteByUser{
		UserID: userID,
	})
	ToolsTest.ReportError(t, err)
}
