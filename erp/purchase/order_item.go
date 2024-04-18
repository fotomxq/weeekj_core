package ERPPurchase

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetPurchaseItemList 获取PurchaseItem列表参数
type ArgsGetPurchaseItemList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//关联订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
}

// GetPurchaseItemList 获取PurchaseItem列表
func GetPurchaseItemList(args *ArgsGetPurchaseItemList) (dataList []FieldsOrderItem, dataCount int64, err error) {
	dataCount, err = purchaseItemDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("order_id", args.OrderID).SetIDQuery("product_id", args.ProductID).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getPurchaseItemByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetPurchaseItemByID 获取PurchaseItem数据包参数
type ArgsGetPurchaseItemByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetPurchaseItemByID 获取PurchaseItem数
func GetPurchaseItemByID(args *ArgsGetPurchaseItemByID) (data FieldsOrderItem, err error) {
	data = getPurchaseItemByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreatePurchaseItem 创建PurchaseItem参数
type ArgsCreatePurchaseItem struct {
	//关联订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id"`
	//采购需求行ID
	PurchaseItemID int64 `db:"purchase_item_id" json:"purchaseItemID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
}

// CreatePurchaseItem 创建PurchaseItem
func CreatePurchaseItem(args *ArgsCreatePurchaseItem) (id int64, err error) {
	//创建数据
	id, err = purchaseItemDB.Insert().SetFields([]string{"order_id", "purchase_item_id", "product_id", "count"}).Add(map[string]any{
		"order_id":         args.OrderID,
		"purchase_item_id": args.PurchaseItemID,
		"product_id":       args.ProductID,
		"count":            args.Count,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdatePurchaseItem 修改PurchaseItem参数
type ArgsUpdatePurchaseItem struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//关联订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id"`
	//采购需求行ID
	PurchaseItemID int64 `db:"purchase_item_id" json:"purchaseItemID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
}

// UpdatePurchaseItem 修改PurchaseItem
func UpdatePurchaseItem(args *ArgsUpdatePurchaseItem) (err error) {
	//更新数据
	err = purchaseItemDB.Update().SetFields([]string{"order_id", "purchase_item_id", "product_id", "count"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"order_id":         args.OrderID,
		"purchase_item_id": args.PurchaseItemID,
		"product_id":       args.ProductID,
		"count":            args.Count,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deletePurchaseItemCache(args.ID)
	//反馈
	return
}

// ArgsDeletePurchaseItem 删除PurchaseItem参数
type ArgsDeletePurchaseItem struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeletePurchaseItem 删除PurchaseItem
func DeletePurchaseItem(args *ArgsDeletePurchaseItem) (err error) {
	//删除数据
	err = purchaseItemDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deletePurchaseItemCache(args.ID)
	//反馈
	return
}

// getPurchaseItemByID 通过ID获取PurchaseItem数据包
func getPurchaseItemByID(id int64) (data FieldsOrderItem) {
	cacheMark := getPurchaseItemCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := purchaseItemDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "order_id", "purchase_item_id", "product_id", "count"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cachePurchaseItemTime)
	return
}

// 缓冲
func getPurchaseItemCacheMark(id int64) string {
	return fmt.Sprint("erp:purchase_item:id.", id)
}

func deletePurchaseItemCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getPurchaseItemCacheMark(id))
}
