package MapArea

import (
	"testing"

	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLGPS "gitee.com/weeekj/weeekj_core/v5/core/sql/gps"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
)

var (
	newDataChild FieldsArea
)

func TestInit3(t *testing.T) {
	TestInit(t)
	TestCreate(t)
}

func TestCheckPointInAreasRand(t *testing.T) {
	data, err := CheckPointInAreasRand(&ArgsCheckPointInAreas{
		MapType: 0,
		Point: CoreSQLGPS.FieldsPoint{
			Latitude:  37.83948,
			Longitude: 112.51696,
		},
		OrgID:    0,
		IsParent: true,
		Mark:     "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
	data, err = CheckPointInAreasRand(&ArgsCheckPointInAreas{
		MapType: 0,
		Point: CoreSQLGPS.FieldsPoint{
			Latitude:  38.83948,
			Longitude: 113.51696,
		},
		OrgID:    0,
		IsParent: true,
		Mark:     "",
	})
	if err == nil {
		t.Error("have area: ", data)
	}
}

func TestCreateChild(t *testing.T) {
	var err error
	newDataChild, _, err = Create(&ArgsCreate{
		OrgID:    1,
		Mark:     "mark2",
		ParentID: newData.ID,
		Name:     "测试c",
		Des:      "测试描述c",
		Country:  86,
		City:     10010,
		MapType:  0,
		Points: CoreSQLGPS.FieldsPoints{
			{
				Latitude:  37.945206,
				Longitude: 112.540460,
			},
			{
				Latitude:  37.791263,
				Longitude: 112.404504,
			},
			{
				Latitude:  112.404504,
				Longitude: 112.598138,
			},
			{
				Latitude:  37.79668,
				Longitude: 112.710748,
			},
			{
				Latitude:  37.93545,
				Longitude: 112.713495,
			},
		},
		Params: []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newDataChild)
}

func TestCheckPointInAreasRandChild(t *testing.T) {
	data, err := CheckPointInAreasRand(&ArgsCheckPointInAreas{
		MapType: 0,
		Point: CoreSQLGPS.FieldsPoint{
			Latitude:  37.83948,
			Longitude: 112.51696,
		},
		OrgID:    0,
		IsParent: false,
		Mark:     "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
	data, err = CheckPointInAreasRand(&ArgsCheckPointInAreas{
		MapType: 0,
		Point: CoreSQLGPS.FieldsPoint{
			Latitude:  37.83948,
			Longitude: 112.51696,
		},
		OrgID:    0,
		IsParent: false,
		Mark:     "",
	})
	if err == nil {
		t.Error("have area: ", data)
	}
}

func TestDelete3(t *testing.T) {
	err := Delete(&ArgsDelete{
		ID:    newDataChild.ID,
		OrgID: 0,
	})
	if err != nil {
		t.Error(err)
	}
	TestDelete(t)
}
