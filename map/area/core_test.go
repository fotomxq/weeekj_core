package MapArea

import (
	"testing"

	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLGPS "github.com/fotomxq/weeekj_core/v5/core/sql/gps"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
)

var (
	isInit  = false
	newData FieldsArea
)

func TestInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	ToolsTest.Init(t)
}

func TestCreate(t *testing.T) {
	var err error
	newData, _, err = Create(&ArgsCreate{
		OrgID:    1,
		Mark:     "mark1",
		ParentID: 0,
		Name:     "测试",
		Des:      "测试描述",
		Country:  86,
		City:     10010,
		MapType:  0,
		Points: CoreSQLGPS.FieldsPoints{
			{
				Latitude:  37.945207,
				Longitude: 112.540461,
			},
			{
				Latitude:  37.791264,
				Longitude: 112.404505,
			},
			{
				Latitude:  37.697874,
				Longitude: 112.598139,
			},
			{
				Latitude:  37.79669,
				Longitude: 112.710749,
			},
			{
				Latitude:  37.93546,
				Longitude: 112.713496,
			},
		},
		Params: []CoreSQLConfig.FieldsConfigType{},
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
		OrgID:    0,
		Mark:     "",
		ParentID: 0,
		Country:  0,
		City:     0,
		MapType:  0,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetByID(t *testing.T) {
	data, err := GetByID(&ArgsGetByID{
		ID:    newData.ID,
		OrgID: 0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdate(t *testing.T) {
	_, err := Update(&ArgsUpdate{
		ID:       newData.ID,
		OrgID:    newData.OrgID,
		Mark:     newData.Mark,
		ParentID: newData.ParentID,
		Name:     newData.Name,
		Des:      newData.Des,
		Country:  newData.Country,
		City:     newData.City,
		MapType:  newData.MapType,
		Points:   newData.Points,
		Params:   newData.Params,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	err := Delete(&ArgsDelete{
		ID:    newData.ID,
		OrgID: 0,
	})
	if err != nil {
		t.Error(err)
	}
}
