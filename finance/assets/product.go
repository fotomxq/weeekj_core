package FinanceAssets

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// 获取产品列表
type ArgsGetProductList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品编码
	Code string `db:"code" json:"code" check:"mark" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetProductList(args *ArgsGetProductList) (dataList []FieldsProduct, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := CoreSQL.GetDeleteSQL(args.IsRemove, "")
	if args.OrgID > 0 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Code != "" {
		where = where + " AND code = :code"
		maps["code"] = args.Code
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"finance_assets_product",
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, name, des, cover_files, des_files, code, currency, price, warehouse_product_ids, mall_commodity_ids, count, params FROM finance_assets_product WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "name", "price", "count"},
	)
	return
}

// 获取产品
type ArgsGetProductByID struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//组织ID
	// 用于检查
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
}

func GetProductByID(args *ArgsGetProductByID) (data FieldsProduct, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, name, des, cover_files, des_files, code, currency, price, warehouse_product_ids, mall_commodity_ids, count, params FROM finance_assets_product WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND "+CoreSQL.GetDeleteSQL(false, ""), args.ID, args.OrgID)
	return
}

// 获取一组产品ID
type ArgsGetProducts struct {
	//一组ID
	IDs pq.Int64Array `db:"ids" json:"ids"`
}

func GetProductsName(args *ArgsGetProducts) (data map[int64]string, err error) {
	type dataType struct {
		//ID
		ID int64 `db:"id" json:"id"`
		//名称
		Name string `db:"name" json:"name"`
	}
	var dataList []dataType
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, name FROM finance_assets_product WHERE id = ANY($1)", args.IDs)
	if err == nil {
		data = map[int64]string{}
		for _, v := range dataList {
			data[v.ID] = v.Name
		}
	}
	return
}

// 创建产品
type ArgsCreateProduct struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//封面列
	// 第一张作为封面
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles" check:"ids" empty:"true"`
	//描述信息
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//产品编码
	Code string `db:"code" json:"code" check:"mark" empty:"true"`
	//储蓄货币类型
	// 采用CoreCurrency匹配
	Currency int `db:"currency" json:"currency" check:"currency"`
	//产品价值
	Price int64 `db:"price" json:"price" check:"int64Than0" empty:"true"`
	//关联的仓储产品
	// 用于仓储识别本资产时使用
	WarehouseProductIDs pq.Int64Array `db:"warehouse_product_ids" json:"warehouseProductIDs" check:"ids" empty:"true"`
	//关联的商品数据
	// 用于商品ID识别本资产时使用
	MallCommodityIDs pq.Int64Array `db:"mall_commodity_ids" json:"mallCommodityIDs" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func CreateProduct(args *ArgsCreateProduct) (data FieldsProduct, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "finance_assets_product", "INSERT INTO finance_assets_product (org_id, name, des, cover_files, des_files, code, currency, price, warehouse_product_ids, mall_commodity_ids, count, params) VALUES (:org_id, :name, :des, :cover_files, :des_files, :code, :currency, :price, :warehouse_product_ids, :mall_commodity_ids, 0, :params)", args, &data)
	return
}

// 修改产品
type ArgsUpdateProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 用于检查
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//封面列
	// 第一张作为封面
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles" check:"ids" empty:"true"`
	//描述信息
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//产品编码
	Code string `db:"code" json:"code" check:"mark" empty:"true"`
	//储蓄货币类型
	// 采用CoreCurrency匹配
	Currency int `db:"currency" json:"currency" check:"currency"`
	//产品价值
	Price int64 `db:"price" json:"price" check:"int64Than0" empty:"true"`
	//关联的仓储产品
	// 用于仓储识别本资产时使用
	WarehouseProductIDs pq.Int64Array `db:"warehouse_product_ids" json:"warehouseProductIDs" check:"ids" empty:"true"`
	//关联的商品数据
	// 用于商品ID识别本资产时使用
	MallCommodityIDs pq.Int64Array `db:"mall_commodity_ids" json:"mallCommodityIDs" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func UpdateProduct(args *ArgsUpdateProduct) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE finance_assets_product SET update_at = NOW(), name = :name, des = :des, cover_files = :cover_files, des_files = :des_files, code = :code, currency = :currency, price = :price, warehouse_product_ids = :warehouse_product_ids, mall_commodity_ids = :mall_commodity_ids, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// 删除产品
type ArgsDeleteProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 用于检查
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

func DeleteProduct(args *ArgsDeleteProduct) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "finance_assets_product", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// 修改总数
type argsUpdateProductCountInc struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//变更数字
	Count int64 `db:"count" json:"count"`
}

func updateProductCountInc(args *argsUpdateProductCountInc) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_assets_product SET update_at = NOW(), count = count + :count WHERE id = :id", args)
	return
}

// 检查产品是否为组织的？
func checkProductAndOrg(productID int64, orgID int64) (err error) {
	type dataType struct {
		//ID
		ID int64 `db:"id" json:"id"`
	}
	var data dataType
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM finance_assets_product WHERE id = $1 AND org_id = $2", productID, orgID)
	if err == nil && data.ID < 1 {
		err = errors.New("org not own product")
	}
	return
}
