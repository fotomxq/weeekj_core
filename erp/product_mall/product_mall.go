package ERPProductMall

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetProductMallList 获取模板列表
type ArgsGetProductMallList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//所属分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetProductMallList 获取品牌列表
func GetProductMallList(args *ArgsGetProductMallList) (dataList []FieldsProductMall, dataCount int64, err error) {
	dataCount, err = productMallDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "price"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetSearchQuery([]string{"product_name"}, args.Search).SetIDQuery("product_id", args.ProductID).SetIDQuery("org_id", args.OrgID).SetIDQuery("category_id", args.CategoryID).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getProductMallData(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// GetProductMall 获取模板
func GetProductMall(id int64, orgID int64) (data FieldsProductMall) {
	data = getProductMallData(id)
	if data.ID < 1 {
		return
	}
	if !CoreFilter.EqID2(orgID, data.OrgID) {
		data = FieldsProductMall{}
		return
	}
	return
}

// ArgsCreateProductMall 创建模板参数
type ArgsCreateProductMall struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品名称
	ProductName string `db:"product_name" json:"productName" check:"des" min:"1" max:"300"`
	//挂出价格
	Price int64 `db:"price" json:"price" check:"int64Than0"`
	//所属分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id"`
}

// CreateProductMall 创建模板
func CreateProductMall(args *ArgsCreateProductMall) (id int64, err error) {
	//创建数据
	id, err = productMallDB.Insert().SetFields([]string{"org_id", "product_id", "product_name", "price", "category_id"}).Add(map[string]any{
		"org_id":       args.OrgID,
		"product_id":   args.ProductID,
		"product_name": args.ProductName,
		"price":        args.Price,
		"category_id":  args.CategoryID,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateProductMall 更新模板参数
type ArgsUpdateProductMall struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品名称
	ProductName string `db:"product_name" json:"productName" check:"des" min:"1" max:"300"`
	//挂出价格
	Price int64 `db:"price" json:"price" check:"int64Than0"`
	//所属分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id"`
}

// UpdateProductMall 更新模板
func UpdateProductMall(args *ArgsUpdateProductMall) (err error) {
	//更新数据
	err = productMallDB.Update().SetFields([]string{"product_id", "product_name", "price", "category_id"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]interface{}{
		"product_id":   args.ProductID,
		"product_name": args.ProductName,
		"price":        args.Price,
		"category_id":  args.CategoryID,
	})
	if err != nil {
		return
	}
	deleteProductMallCache(args.ID)
	return
}

// ArgsDeleteProductMall 删除模板参数
type ArgsDeleteProductMall struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteProductMall 删除模板
func DeleteProductMall(args *ArgsDeleteProductMall) (err error) {
	err = productMallDB.Delete().NeedSoft(true).AddWhereID(args.ID).AddWhereOrgID(args.OrgID).ExecNamed(nil)
	if err != nil {
		return
	}
	deleteProductMallCache(args.ID)
	return
}

// getProductMallData 获取模板数据
func getProductMallData(id int64) (data FieldsProductMall) {
	cacheMark := getProductMallCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := productMallDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "product_id", "product_name", "price", "category_id"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheProductMallTime)
	return
}

// 缓冲
func getProductMallCacheMark(id int64) string {
	return fmt.Sprint("erp:product:mall:id.", id)
}

func deleteProductMallCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getProductMallCacheMark(id))
}
