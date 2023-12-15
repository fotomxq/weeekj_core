package ERPWarehouse

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
)

// ArgsBatchOutAuto 按照规则自动完成产品出库参数
type ArgsBatchOutAuto struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域ID
	AreaID int64 `db:"area_id" json:"areaID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//调用模块
	ActionSystem string `db:"action_system" json:"actionSystem" check:"mark" empty:"true"`
	//调用模块ID
	ActionID int64 `db:"action_id" json:"actionID" check:"id" empty:"true"`
}

// BatchOutAuto 按照规则自动完成产品出库
func BatchOutAuto(args *ArgsBatchOutAuto) (errCode string, err error) {
	//锁定机制
	batchWriteLock.Lock()
	defer batchWriteLock.Unlock()
	//检查产品的在库数量
	productStoreCount := GetStoreProductCount(args.OrgID, args.WarehouseID, args.AreaID, args.ProductID)
	if productStoreCount < args.Count {
		errCode = "err_erp_warehouse_product_no_store"
		err = errors.New(fmt.Sprint("too more batch count, product store count: ", productStoreCount, ", need count: ", args.Count))
		return
	}
	//剩余处理的批次
	remainCount := args.Count + 0
	//等待出库的批次
	// 其中v.Count在最后一轮如果超出剩余待出库数量，则会被覆盖为需出库数量
	var waitOutBatchList []FieldsBatch
	//获取批次列表
	var page int64 = 1
	for {
		batchList, _, _ := GetBatchList(&ArgsGetBatchList{
			Pages: CoreSQL2.ArgsPages{
				Page: page,
				Max:  100,
				Sort: "id",
				Desc: false,
			},
			OrgID:        args.OrgID,
			WarehouseID:  args.WarehouseID,
			AreaID:       args.AreaID,
			LocationID:   -1,
			ProductID:    args.ProductID,
			FactoryBatch: "",
			SystemBatch:  "",
			IsRemove:     false,
			Search:       "",
		})
		if len(batchList) < 1 {
			break
		}
		for k := 0; k < len(batchList); k++ {
			v := batchList[k]
			//该批次全部出库
			if v.Count <= remainCount {
				waitOutBatchList = append(waitOutBatchList, v)
				remainCount = remainCount - v.Count
				continue
			}
			//该批次部分出库
			if v.Count > remainCount {
				v.Count = remainCount + 0
				waitOutBatchList = append(waitOutBatchList, v)
				remainCount = remainCount - v.Count
				break
			}
		}
		//没有需要出库的物资则跳出
		if remainCount < 1 {
			break
		}
		page += 1
	}
	//如果批次遍历后，还是不足，说明库存不足
	if remainCount > 0 {
		errCode = "err_erp_warehouse_product_no_store"
		err = errors.New("too more batch count")
		return
	}
	//开始批次出库操作
	for _, v := range waitOutBatchList {
		errCode, err = BatchOut(&ArgsBatchOut{
			ID:    v.ID,
			OrgID: v.OrgID,
			Count: v.Count,
		})
		if err != nil {
			return
		}
	}
	//反馈
	return
}
