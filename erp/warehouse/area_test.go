package ERPWarehouse

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	"testing"
)

var (
	newAreaData FieldsArea
)

func TestAreaInit(t *testing.T) {
	TestInit(t)
	TestCreateWarehouse(t)
	TestGetWarehouseList(t)
}

func TestCreateArea(t *testing.T) {
	err := CreateArea(&ArgsCreateArea{
		OrgID:       newWarehouseData.OrgID,
		WarehouseID: newWarehouseData.ID,
		Name:        "测试仓储区域A",
		Location:    "位置A",
		Weight:      0,
		SizeW:       0,
		SizeH:       0,
		SizeZ:       0,
		MapType:     0,
		Longitude:   0,
		Latitude:    0,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new area data")
	}
}

func TestGetAreaList(t *testing.T) {
	dataList, dataCount, err := GetAreaList(&ArgsGetAreaList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:       newWarehouseData.OrgID,
		WarehouseID: newWarehouseData.ID,
		IsRemove:    false,
		Search:      "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("get area list: ", len(dataList), ", count: ", dataCount)
	if len(dataList) < 1 {
		t.Error("no data")
		return
	}
	newAreaData = dataList[0]
}

func TestGetAreaName(t *testing.T) {
	data := GetAreaName(newAreaData.ID)
	t.Log("get area name: ", data)
}

func TestUpdateArea(t *testing.T) {
	err := UpdateArea(&ArgsUpdateArea{
		ID:        newAreaData.ID,
		OrgID:     newAreaData.OrgID,
		Name:      newAreaData.Name,
		Location:  newAreaData.Location,
		Weight:    newAreaData.Weight,
		SizeW:     newAreaData.SizeW,
		SizeH:     newAreaData.SizeH,
		SizeZ:     newAreaData.SizeZ,
		MapType:   newAreaData.MapType,
		Longitude: newAreaData.Longitude,
		Latitude:  newAreaData.Latitude,
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteArea(t *testing.T) {
	err := DeleteArea(&ArgsDeleteArea{
		ID:    newAreaData.ID,
		OrgID: newAreaData.OrgID,
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAreaClear(t *testing.T) {
	TestDeleteWarehouse(t)
	TestClear(t)
}
