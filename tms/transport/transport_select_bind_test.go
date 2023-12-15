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
	newMapArea2 MapArea.FieldsArea
)

// 配送员选择测试聚合
func TestTransportSelectBindInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
	TestOrg.LocalCreateBind(t)
	//创建两个分区，用于多分区混合测试
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
	newMapArea2, _, err = MapArea.Create(&MapArea.ArgsCreate{
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
				Longitude: 100,
				Latitude:  100,
			},
			{
				Longitude: 200,
				Latitude:  100,
			},
			{
				Longitude: 200,
				Latitude:  200,
			},
			{
				Longitude: 100,
				Latitude:  200,
			},
		},
		Params: []CoreSQLConfig.FieldsConfigType{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new map area: ", newMapArea2.ID)
	}
	//绑定不同的人员
	newBindData, err = SetBind(&ArgsSetBind{
		OrgID:          TestOrg.OrgData.ID,
		BindID:         TestOrg.BindData.ID,
		MapAreaID:      newMapArea.ID,
		MoreMapAreaIDs: []int64{},
		Params:         []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newBindData)
	t.Log("new bind data, bind id: ", newBindData.BindID)
	newBindData, err = SetBind(&ArgsSetBind{
		OrgID:          TestOrg.OrgData.ID,
		BindID:         TestOrg.BindData.ID,
		MapAreaID:      newMapArea.ID,
		MoreMapAreaIDs: []int64{},
		Params:         []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newBindData)
	t.Log("new bind data, bind id: ", newBindData.BindID)
}

func TestTransportSelectBind(t *testing.T) {

}

func TestTransportSelectBindClear(t *testing.T) {
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
