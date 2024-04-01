package RestaurantPurchase

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetRestaurantPurchaseList 获取RestaurantPurchase列表参数
type ArgsGetRestaurantPurchaseList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetRestaurantPurchaseList 获取RestaurantPurchase列表
func GetRestaurantPurchaseList(args *ArgsGetRestaurantPurchaseList) (dataList []FieldsPurchaseAnalysis, dataCount int64, err error) {
	dataCount, err = restaurantPurchaseDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetIDQuery("store_id", args.StoreID).SetSearchQuery([]string{"name", "remark"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getRestaurantPurchaseByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetRestaurantPurchaseByID 获取RestaurantPurchase数据包参数
type ArgsGetRestaurantPurchaseByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
}

// GetRestaurantPurchaseByID 获取RestaurantPurchase数
func GetRestaurantPurchaseByID(args *ArgsGetRestaurantPurchaseByID) (data FieldsPurchaseAnalysis, err error) {
	data = getRestaurantPurchaseByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.StoreID, data.StoreID) {
		err = errors.New("no data")
		return
	}
	return
}

// GetRestaurantPurchaseMargeByID 获取打包数据
func GetRestaurantPurchaseMargeByID(args *ArgsGetRestaurantPurchaseByID) (headData FieldsPurchaseAnalysis, itemList []FieldsPurchaseAnalysisItem, err error) {
	headData, err = GetRestaurantPurchaseByID(args)
	if err != nil {
		return
	}
	itemList, _, _ = GetRestaurantPurchaseItemList(&ArgsGetRestaurantPurchaseItemList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  999,
			Sort: "id",
			Desc: false,
		},
		OrgID:              -1,
		StoreID:            -1,
		PurchaseAnalysisID: headData.ID,
		IsRemove:           false,
		Search:             "",
	})
	return
}

// GetRestaurantPurchaseNameByID 获取菜品名称
func GetRestaurantPurchaseNameByID(id int64) (name string) {
	data := getRestaurantPurchaseByID(id)
	if data.ID < 1 {
		return
	}
	return data.Name
}

// ArgsCreateRestaurantPurchase 创建RestaurantPurchase参数
type ArgsCreateRestaurantPurchase struct {
	//发生采购时间
	PurchaseAt time.Time `db:"purchase_at" json:"purchaseAt"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

// CreateRestaurantPurchase 创建RestaurantPurchase
func CreateRestaurantPurchase(args *ArgsCreateRestaurantPurchase) (id int64, err error) {
	//创建数据
	id, err = restaurantPurchaseDB.Insert().SetFields([]string{"purchase_at", "org_id", "store_id", "name"}).Add(map[string]any{
		"purchase_at": args.PurchaseAt,
		"org_id":      args.OrgID,
		"store_id":    args.StoreID,
		"name":        args.Name,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsCreateRestaurantPurchaseMargeByDay 打包创建某一天的菜谱参数
type ArgsCreateRestaurantPurchaseMargeByDay struct {
	//头部关键信息
	HeaderData ArgsCreateRestaurantPurchase `json:"headerData"`
	//行信息列
	RowData []ArgsCreateRestaurantPurchaseMargeByDayItem `json:"rowData"`
}

type ArgsCreateRestaurantPurchaseMargeByDayItem struct {
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

// CreateRestaurantPurchaseMargeByDay 打包创建某一天的菜谱
func CreateRestaurantPurchaseMargeByDay(args *ArgsCreateRestaurantPurchaseMargeByDay) (id int64, err error) {
	//创建头部
	id, err = CreateRestaurantPurchase(&args.HeaderData)
	if err != nil {
		return
	}
	//批量创建行
	for _, v := range args.RowData {
		_, err = CreateRestaurantPurchaseItem(&ArgsCreateRestaurantPurchaseItem{
			OrgID:              args.HeaderData.OrgID,
			StoreID:            args.HeaderData.StoreID,
			PurchaseAnalysisID: id,
			Name:               v.Name,
			Weight:             v.Weight,
			Price:              v.Price,
			TotalPrice:         v.TotalPrice,
		})
		if err != nil {
			return
		}
	}
	//反馈
	return
}

// ArgsUpdateRestaurantPurchase 修改RestaurantPurchase参数
type ArgsUpdateRestaurantPurchase struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//发生采购时间
	PurchaseAt time.Time `db:"purchase_at" json:"purchaseAt"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}

// UpdateRestaurantPurchase 修改RestaurantPurchase
func UpdateRestaurantPurchase(args *ArgsUpdateRestaurantPurchase) (err error) {
	//更新数据
	err = restaurantPurchaseDB.Update().SetFields([]string{"purchase_at", "org_id", "store_id", "name"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"purchase_at": args.PurchaseAt,
		"org_id":      args.OrgID,
		"store_id":    args.StoreID,
		"name":        args.Name,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteRestaurantPurchaseCache(args.ID)
	//反馈
	return
}

// ArgsDeleteRestaurantPurchase 删除RestaurantPurchase参数
type ArgsDeleteRestaurantPurchase struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteRestaurantPurchase 删除RestaurantPurchase
func DeleteRestaurantPurchase(args *ArgsDeleteRestaurantPurchase) (err error) {
	//删除数据
	err = restaurantPurchaseDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteRestaurantPurchaseCache(args.ID)
	//反馈
	return
}

// getRestaurantPurchaseByID 通过ID获取RestaurantPurchase数据包
func getRestaurantPurchaseByID(id int64) (data FieldsPurchaseAnalysis) {
	cacheMark := getRestaurantPurchaseCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := restaurantPurchaseDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "purchase_at", "org_id", "store_id", "name"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheRestaurantPurchaseTime)
	return
}

// 缓冲
func getRestaurantPurchaseCacheMark(id int64) string {
	return fmt.Sprint("restaurant:weekly_recipe:id.", id)
}

func deleteRestaurantPurchaseCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getRestaurantPurchaseCacheMark(id))
}
