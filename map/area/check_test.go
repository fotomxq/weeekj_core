package MapArea

import (
	CoreSQLGPS "gitee.com/weeekj/weeekj_core/v5/core/sql/gps"
	"testing"
)

func TestInit2(t *testing.T) {
	TestInit(t)
	TestCreate(t)
}

func TestCheckPointInAreas(t *testing.T) {
	dataList, err := CheckPointInAreas(&ArgsCheckPointInAreas{
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
		t.Log(dataList)
	}
}

func TestDelete2(t *testing.T) {
	TestDelete(t)
}
