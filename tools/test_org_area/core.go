package TestOrgArea

import (
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	"testing"

	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLGPS "github.com/fotomxq/weeekj_core/v5/core/sql/gps"
	MapArea "github.com/fotomxq/weeekj_core/v5/map/area"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
)

var (
	isInit = false
	//ParentAreaData 行政分区
	ParentAreaData MapArea.FieldsArea
	//AreaData 业务分区
	AreaData MapArea.FieldsArea
)

// Init mark: 业务分区的标识码
func Init(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
		//OrgCore.Init(true, true)
		_ = OrgCore.Init()
	}
	isInit = true
	TestOrg.LocalCreateOrg(t)
}

func CreateParentArea(t *testing.T) {
	Init(t)
	var err error
	var errCode string
	ParentAreaData, errCode, err = MapArea.Create(&MapArea.ArgsCreate{
		OrgID:    TestOrg.OrgData.ID,
		Mark:     "",
		ParentID: 0,
		Name:     "测试分区",
		Des:      "测试分区描述",
		Country:  86,
		City:     10010,
		MapType:  0,
		Points: []CoreSQLGPS.FieldsPoint{
			{
				Longitude: 0,
				Latitude:  0,
			},
			{
				Longitude: 0,
				Latitude:  15,
			},
			{
				Longitude: 15,
				Latitude:  15,
			},
			{
				Longitude: 15,
				Latitude:  0,
			},
		},
		Params: []CoreSQLConfig.FieldsConfigType{},
	})
	if err != nil {
		t.Error("无法创建行政分区, ", err, ", ", errCode)
		return
	}
	t.Log("创建行政分区, ", ParentAreaData)
}

func CreateChildArea(t *testing.T, mark string) {
	if ParentAreaData.ID < 1 {
		CreateParentArea(t)
	}
	var err error
	var errCode string
	AreaData, errCode, err = MapArea.Create(&MapArea.ArgsCreate{
		OrgID:    TestOrg.OrgData.ID,
		Mark:     mark,
		ParentID: ParentAreaData.ID,
		Name:     "测试分区",
		Des:      "测试分区描述",
		Country:  86,
		City:     10010,
		MapType:  0,
		Points: []CoreSQLGPS.FieldsPoint{
			{
				Longitude: 1,
				Latitude:  1,
			},
			{
				Longitude: 1,
				Latitude:  10,
			},
			{
				Longitude: 10,
				Latitude:  10,
			},
			{
				Longitude: 10,
				Latitude:  1,
			},
		},
		Params: []CoreSQLConfig.FieldsConfigType{},
	})
	if err != nil {
		t.Error("无法创建业务分区, ", err, ", errCode: ", errCode)
		return
	}
	t.Log("创建业务分区, ", AreaData)
}

func Clear(t *testing.T) {
	if AreaData.ID > 0 {
		err := MapArea.Delete(&MapArea.ArgsDelete{
			ID:    AreaData.ID,
			OrgID: 0,
		})
		ToolsTest.ReportError(t, err)
		AreaData.ID = 0
	}
	if ParentAreaData.ID > 0 {
		err := MapArea.Delete(&MapArea.ArgsDelete{
			ID:    ParentAreaData.ID,
			OrgID: 0,
		})
		ToolsTest.ReportError(t, err)
		ParentAreaData.ID = 0
	}
	TestOrg.LocalClear(t)
}
