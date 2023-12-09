package ERPWarehouse

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL2 "gitee.com/weeekj/weeekj_core/v5/core/sql2"
	"testing"
)

var (
	newBatchData FieldsBatch
)

// 批次测试初始化
func TestBatchInit(t *testing.T) {
	TestLocationInit(t)
	TestCreateProduct(t)
	t.Log("new product data: ", newERPProductData)
}

func TestCreateBatch(t *testing.T) {
	var errCode string
	var err error
	newBatchData, errCode, err = CreateBatch(&ArgsBatchCreate{
		Sn:           CoreFilter.GetRandStr4(10),
		OrgID:        newWarehouseData.OrgID,
		WarehouseID:  newWarehouseData.ID,
		AreaID:       newAreaData.ID,
		LocationID:   newLocationData.ID,
		ProductID:    newERPProductData.ID,
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDays(30).Time,
		FactoryBatch: "factoryBatch",
		SystemBatch:  "systemBatch",
		CostPrice:    newERPProductData.CostPrice,
		CostPriceTax: newERPProductData.TaxCostPrice,
		Count:        20,
		Des:          "测试批次入库",
	})
	if err != nil {
		t.Error(errCode, err)
		return
	}
	t.Log("new batch, ", newBatchData.ID)
	productStoreCount := GetStoreProductCount(newWarehouseData.OrgID, newWarehouseData.ID, newAreaData.ID, newERPProductData.ID)
	t.Log("now product store count: ", productStoreCount)
}

func TestGetBatchListByProductID(t *testing.T) {
	dataList, err := GetBatchListByProductID(newERPProductData.OrgID, newERPProductData.ID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("get batch by product id: ", newERPProductData.ID, ", batch len: ", len(dataList))
}

func TestGetBatchByID(t *testing.T) {
	data := getBatchByID(newBatchData.ID)
	t.Log("get batch by id: ", newBatchData.ID, ", data: ", data.ID)
	newBatchData = data
}

// 批次列表测试
func TestGetBatchList(t *testing.T) {
	dataList, dataCount, err := GetBatchList(&ArgsGetBatchList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:        -1,
		WarehouseID:  -1,
		AreaID:       -1,
		LocationID:   -1,
		ProductID:    -1,
		FactoryBatch: "",
		SystemBatch:  "",
		IsRemove:     false,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("batch list, len： ", len(dataList), ", count: ", dataCount)
}

// 批次出库测试
func TestBatchOut(t *testing.T) {
	productStoreCount := GetStoreProductCount(newWarehouseData.OrgID, newWarehouseData.ID, newAreaData.ID, newERPProductData.ID)
	t.Log("now product store count: ", productStoreCount)
	errCode, err := BatchOut(&ArgsBatchOut{
		ID:           newBatchData.ID,
		OrgID:        newBatchData.OrgID,
		Count:        9,
		ActionSystem: "",
		ActionID:     0,
	})
	if err != nil {
		t.Error(errCode, err)
		return
	}
	t.Log("batch out success")
	productStoreCount2 := GetStoreProductCount(newWarehouseData.OrgID, newWarehouseData.ID, newAreaData.ID, newERPProductData.ID)
	t.Log("now product store count: ", productStoreCount2)
}

func TestBatchAutoOut(t *testing.T) {
	errCode, err := BatchOutAuto(&ArgsBatchOutAuto{
		OrgID:        newWarehouseData.OrgID,
		WarehouseID:  newWarehouseData.ID,
		AreaID:       newAreaData.ID,
		ProductID:    newERPProductData.ID,
		Count:        3,
		ActionSystem: "simple_sell",
		ActionID:     1,
	})
	if err != nil {
		t.Error("err code: ", errCode, ", err: ", err)
		return
	}
	productStoreCount := GetStoreProductCount(newWarehouseData.OrgID, newWarehouseData.ID, newAreaData.ID, newERPProductData.ID)
	t.Log("now product store count: ", productStoreCount)
}

func TestBatchClear(t *testing.T) {
	_, _ = BatchOutAuto(&ArgsBatchOutAuto{
		OrgID:        newWarehouseData.OrgID,
		WarehouseID:  newWarehouseData.ID,
		AreaID:       newAreaData.ID,
		ProductID:    newERPProductData.ID,
		Count:        8,
		ActionSystem: "simple_sell_clear",
		ActionID:     1,
	})
	TestLocationClear(t)
}
