package RestaurantPurchase

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetRestaurantPurchaseItemList 获取RestaurantPurchaseItem列表参数
type ArgsGetRestaurantPurchaseItemList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID" check:"id" empty:"true"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
	//原材料采购台账ID
	PurchaseAnalysisID int64 `db:"purchase_analysis_id" json:"purchaseAnalysisID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetRestaurantPurchaseItemList 获取RestaurantPurchaseItem列表
func GetRestaurantPurchaseItemList(args *ArgsGetRestaurantPurchaseItemList) (dataList []FieldsPurchaseAnalysisItem, dataCount int64, err error) {
	dataCount, err = restaurantPurchaseItemDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("raw_org_id", args.RawOrgID).SetIDQuery("org_id", args.OrgID).SetIDQuery("store_id", args.StoreID).SetIDQuery("purchase_analysis_id", args.PurchaseAnalysisID).SetSearchQuery([]string{"name"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getRestaurantPurchaseItemByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetRestaurantPurchaseItemByID 获取RestaurantPurchaseItem数据包参数
type ArgsGetRestaurantPurchaseItemByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetRestaurantPurchaseItemByID 获取RestaurantPurchaseItem数
func GetRestaurantPurchaseItemByID(args *ArgsGetRestaurantPurchaseItemByID) (data FieldsPurchaseAnalysisItem, err error) {
	data = getRestaurantPurchaseItemByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// GetRestaurantPurchaseItemNameByID 获取菜品名称
func GetRestaurantPurchaseItemNameByID(id int64) (name string) {
	data := getRestaurantPurchaseItemByID(id)
	if data.ID < 1 {
		return
	}
	return data.Name
}

// ArgsCreateRestaurantPurchaseItem 创建RestaurantPurchaseItem参数
type ArgsCreateRestaurantPurchaseItem struct {
	//组织ID
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	//原材料采购台账ID
	PurchaseAnalysisID int64 `db:"purchase_analysis_id" json:"purchaseAnalysisID" check:"id"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//原材料重量 默认kg
	Weight int64 `db:"weight" json:"weight" check:"int64Than0" empty:"true"`
	//单价
	Price int64 `db:"price" json:"price" check:"int64Than0" empty:"true"`
	//总价
	TotalPrice int64 `db:"total_price" json:"totalPrice" check:"int64Than0" empty:"true"`
}

// CreateRestaurantPurchaseItem 创建RestaurantPurchaseItem
func CreateRestaurantPurchaseItem(args *ArgsCreateRestaurantPurchaseItem) (id int64, err error) {
	//创建数据
	id, err = restaurantPurchaseItemDB.Insert().SetFields([]string{"raw_org_id", "org_id", "store_id", "purchase_analysis_id", "name", "weight", "price", "total_price"}).Add(map[string]any{
		"purchase_analysis_id": args.PurchaseAnalysisID,
		"name":                 args.Name,
		"weight":               args.Weight,
		"price":                args.Price,
		"total_price":          args.TotalPrice,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateRestaurantPurchaseItem 修改RestaurantPurchaseItem参数
type ArgsUpdateRestaurantPurchaseItem struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//原材料采购台账ID
	PurchaseAnalysisID int64 `db:"purchase_analysis_id" json:"purchaseAnalysisID" check:"id"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//原材料重量 默认kg
	Weight int64 `db:"weight" json:"weight" check:"int64Than0" empty:"true"`
	//单价
	Price int64 `db:"price" json:"price" check:"int64Than0" empty:"true"`
	//总价
	TotalPrice int64 `db:"total_price" json:"totalPrice" check:"int64Than0" empty:"true"`
}

// UpdateRestaurantPurchaseItem 修改RestaurantPurchaseItem
func UpdateRestaurantPurchaseItem(args *ArgsUpdateRestaurantPurchaseItem) (err error) {
	//更新数据
	err = restaurantPurchaseItemDB.Update().SetFields([]string{"purchase_analysis_id", "name", "weight", "price", "total_price"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"purchase_analysis_id": args.PurchaseAnalysisID,
		"name":                 args.Name,
		"weight":               args.Weight,
		"price":                args.Price,
		"total_price":          args.TotalPrice,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteRestaurantPurchaseItemCache(args.ID)
	//反馈
	return
}

// ArgsDeleteRestaurantPurchaseItem 删除RestaurantPurchaseItem参数
type ArgsDeleteRestaurantPurchaseItem struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteRestaurantPurchaseItem 删除RestaurantPurchaseItem
func DeleteRestaurantPurchaseItem(args *ArgsDeleteRestaurantPurchaseItem) (err error) {
	//删除数据
	err = restaurantPurchaseItemDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteRestaurantPurchaseItemCache(args.ID)
	//反馈
	return
}

// getRestaurantPurchaseItemByID 通过ID获取RestaurantPurchaseItem数据包
func getRestaurantPurchaseItemByID(id int64) (data FieldsPurchaseAnalysisItem) {
	cacheMark := getRestaurantPurchaseItemCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := restaurantPurchaseItemDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "raw_org_id", "org_id", "store_id", "purchase_analysis_id", "name", "weight", "price", "total_price"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheRestaurantPurchaseItemTime)
	return
}

// 缓冲
func getRestaurantPurchaseItemCacheMark(id int64) string {
	return fmt.Sprint("restaurant:restaurant_purchase_item:id.", id)
}

func deleteRestaurantPurchaseItemCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getRestaurantPurchaseItemCacheMark(id))
}
