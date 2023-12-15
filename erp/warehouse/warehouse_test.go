package ERPWarehouse

import (
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newWarehouseData FieldsWarehouse
)

func TestWarehouseInit(t *testing.T) {
	TestInit(t)
}

func TestCreateWarehouse(t *testing.T) {
	err := CreateWarehouse(&ArgsCreateWarehouse{
		OrgID:       TestOrg.OrgData.ID,
		Name:        "测试仓储",
		Weight:      0,
		SizeW:       0,
		SizeH:       0,
		SizeZ:       0,
		AddressData: CoreSQLAddress.FieldsAddress{},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("create new warehouse")
}

func TestGetWarehouseList(t *testing.T) {
	dataList, dataCount, err := GetWarehouseList(&ArgsGetWarehouseList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:    -1,
		IsRemove: false,
		Search:   "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("warehouse list, ", dataCount, ", len: ", len(dataList))
	if len(dataList) < 1 {
		t.Error("no data")
		return
	}
	newWarehouseData = dataList[0]
}

func TestGetWarehouseName(t *testing.T) {
	data := GetWarehouseName(newWarehouseData.ID)
	t.Log("new warehouse name: ", data)
}

func TestUpdateWarehouse(t *testing.T) {
	err := UpdateWarehouse(&ArgsUpdateWarehouse{
		ID:          newWarehouseData.ID,
		OrgID:       newWarehouseData.OrgID,
		Name:        newWarehouseData.Name,
		Weight:      newWarehouseData.Weight,
		SizeW:       newWarehouseData.SizeW,
		SizeH:       newWarehouseData.SizeH,
		SizeZ:       newWarehouseData.SizeZ,
		AddressData: newWarehouseData.AddressData,
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteWarehouse(t *testing.T) {
	err := DeleteWarehouse(&ArgsDeleteWarehouse{
		ID:    newWarehouseData.ID,
		OrgID: newWarehouseData.OrgID,
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestWarehouseClear(t *testing.T) {
	TestClear(t)
}
