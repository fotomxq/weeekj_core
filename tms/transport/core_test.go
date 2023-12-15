package TMSTransport

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLGPS "github.com/fotomxq/weeekj_core/v5/core/sql/gps"
	MapArea "github.com/fotomxq/weeekj_core/v5/map/area"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	isInit     = false
	newMapArea MapArea.FieldsArea
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
	TestOrg.LocalCreateBind(t)
	var err error
	newMapArea, _, err = MapArea.Create(&MapArea.ArgsCreate{
		OrgID:    TestOrg.OrgData.ID,
		Mark:     "tms",
		ParentID: 0,
		Name:     "测试配送分区",
		Des:      "",
		Country:  86,
		City:     10000,
		MapType:  0,
		Points: []CoreSQLGPS.FieldsPoint{
			{
				Longitude: 0,
				Latitude:  0,
			},
			{
				Longitude: 100,
				Latitude:  0,
			},
			{
				Longitude: 100,
				Latitude:  100,
			},
			{
				Longitude: 0,
				Latitude:  100,
			},
		},
		Params: []CoreSQLConfig.FieldsConfigType{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new map area: ", newMapArea.ID)
	}
}

func TestClear(t *testing.T) {
	var err error
	err = MapArea.Delete(&MapArea.ArgsDelete{
		ID:    newMapArea.ID,
		OrgID: newMapArea.OrgID,
	})
	if err != nil {
		t.Error(err)
	}
	TestOrg.LocalClear(t)
}
