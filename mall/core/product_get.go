package MallCore

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLIDs "gitee.com/weeekj/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"strings"
)

// ArgsGetProductList 获取商品列表参数
type ArgsGetProductList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//商品类型
	// 0 普通商品; 1 关联选项商品; 2 虚拟商品
	ProductType int `db:"product_type" json:"productType" check:"intThan0" empty:"true"`
	//是否为虚拟商品
	// 不会发生配送处理
	NeedIsVirtual bool `db:"need_is_virtual" json:"needIsVirtual" check:"bool"`
	IsVirtual     bool `db:"is_virtual" json:"isVirtual" check:"bool"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//是否已经发布
	NeedIsPublish bool `db:"need_is_publish" json:"needIsPublish" check:"bool"`
	IsPublish     bool `db:"is_publish" json:"isPublish" check:"bool"`
	//价格区间
	PriceMin int64 `db:"price_min" json:"priceMin" check:"price" empty:"true"`
	PriceMax int64 `db:"price_max" json:"priceMax" check:"price" empty:"true"`
	//是否包含会员
	NeedHaveVIP bool `json:"needHaveVIP" check:"bool"`
	HaveVIP     bool `json:"haveVIP" check:"bool"`
	//可用票据
	Tickets pq.Int64Array `json:"tickets" check:"ids" empty:"true"`
	//是否有库存
	NeedHaveCount bool `json:"needHaveCount" check:"bool"`
	HaveCount     bool `json:"haveCount" check:"bool"`
	//是否存在积分消费
	NeedHaveIntegral bool `json:"needHaveIntegral" check:"bool"`
	HaveIntegral     bool `json:"needIntegral" check:"bool"`
	//上级ID
	// 如果不给予，则强制忽略历史记录
	// 标记为历史记录，<1为最高级，其他级别代表历史记录
	ParentID int64 `db:"parent_id" json:"parentID" check:"int64Than0" empty:"true"`
	//配送费计费模版ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id" empty:"true"`
	//关联仓库的产品
	// ERPProduct模块
	WarehouseProductID int64 `db:"warehouse_product_id" json:"warehouseProductID"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetProductList 获取商品列表
func GetProductList(args *ArgsGetProductList) (dataList []FieldsCore, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ProductType > -1 {
		where = where + " AND product_type = :product_type"
		maps["product_type"] = args.ProductType
	}
	if args.NeedIsVirtual {
		where = where + " AND is_virtual = :is_virtual"
		maps["is_virtual"] = args.IsVirtual
	}
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.NeedIsPublish {
		if args.IsPublish {
			where = where + " AND publish_at > to_timestamp(1000000)"
		} else {
			where = where + " AND publish_at <= to_timestamp(1000000)"
		}
	}
	if args.PriceMin > 0 {
		where = where + " AND price_real >= :price_min"
		maps["price_min"] = args.PriceMin
	}
	if args.PriceMax > 0 {
		where = where + " AND price_real <= :price_max"
		maps["price_max"] = args.PriceMax
	}
	if args.NeedHaveVIP {
		if args.HaveVIP {
			where = where + " AND json_array_length(user_sub_price) > 0"
		} else {
			where = where + " AND json_array_length(user_sub_price) < 1"
		}
	}
	if len(args.Tickets) > 0 {
		where = where + " AND user_ticket @> :tickets"
		maps["tickets"] = args.Tickets
	}
	if args.NeedHaveCount {
		if args.HaveCount {
			where = where + " AND count > 0"
		} else {
			where = where + " AND count < 1"
		}
	}
	if args.NeedHaveIntegral {
		if args.HaveIntegral {
			where = where + " AND integral > 0"
		} else {
			where = where + " AND integral = 0"
		}
	}
	if args.ParentID > 0 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	} else {
		where = where + " AND parent_id = 0"
	}
	if args.TransportID > 0 {
		where = where + " AND transport_id = :transport_id"
		maps["transport_id"] = args.TransportID
	}
	if args.WarehouseProductID > -1 {
		where = where + " AND warehouse_product_id = :warehouse_product_id"
		maps["warehouse_product_id"] = args.WarehouseProductID
	}
	if args.Search != "" {
		args.Search = strings.ReplaceAll(args.Search, " ", "")
		where = where + " AND (REPLACE(title, ' ', '') ILIKE '%' || :search || '%' OR REPLACE(title_des, ' ', '') ILIKE '%' || :search || '%' OR REPLACE(des, ' ', '') ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if args.Pages.Sort == "sort" {
		args.Pages.Sort = "sort, update_at"
	}
	tableName := "mall_core"
	var rawList []FieldsCore
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "sort", "sort, update_at", "expire_at", "audit_at", "price", "price_real", "integral", "buy_count"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getProductByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetProduct 获取指定商品参数
type ArgsGetProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetProduct 获取指定商品
func GetProduct(args *ArgsGetProduct) (data FieldsCore, err error) {
	data = getProductByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data = FieldsCore{}
		err = errors.New("no data")
		return
	}
	return
}

func GetProductNoErr(id int64) (data FieldsCore) {
	data = getProductByID(id)
	return
}

// GetERPProductBindMallProductCount 获取的ERP产品数量关联的商品数量
func GetERPProductBindMallProductCount(orgID int64, erpProductID int64, erpProductCode string, haveNoPublish bool) (count int64) {
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM mall_core WHERE org_id = $1 AND parent_id = 0 AND delete_at < to_timestamp(1000000) AND (warehouse_product_id = $2 OR code = $3 OR other_options -> 'dataList' @> ('[{\"code\": \"' || $3 || '\"}]')::jsonb) AND ($4 = true OR ($4 = false AND publish_at > to_timestamp(1000000)))", orgID, erpProductID, erpProductCode, haveNoPublish)
	if err != nil {
		return
	}
	return
}

// GetProductTop 是否采用追溯的商品获取产品
func GetProductTop(args *ArgsGetProduct) (data FieldsCore, err error) {
	data, err = GetProduct(args)
	if err != nil {
		return
	}
	if data.ParentID < 1 {
		return
	}
	data = getProductByID(data.ParentID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetProducts 获取指定商品组参数
type ArgsGetProducts struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetProducts 获取指定商品组
func GetProducts(args *ArgsGetProducts) (dataList []FieldsCore, err error) {
	for _, v := range args.IDs {
		vData := getProductByID(v)
		if vData.ID < 1 || !CoreFilter.EqID2(args.OrgID, vData.OrgID) || (!args.HaveRemove && CoreSQL.CheckTimeHaveData(vData.DeleteAt)) {
			continue
		}
		if vData.ParentID > 0 {
			vDataParent := getProductByID(vData.ParentID)
			if vDataParent.ID < 1 || !CoreFilter.EqID2(args.OrgID, vDataParent.OrgID) || (!args.HaveRemove && CoreSQL.CheckTimeHaveData(vDataParent.DeleteAt)) {
				continue
			}
			dataList = append(dataList, vDataParent)
		} else {
			dataList = append(dataList, vData)
		}
	}
	if len(dataList) < 1 {
		err = errors.New("data is empty")
		return
	}
	return
}

// GetProductsName 获取指定商品组名称
func GetProductsName(args *ArgsGetProducts) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsOrgTitleAndDelete("mall_core", args.IDs, args.OrgID, args.HaveRemove)
	return
}

func GetProductName(productID int64) string {
	data := getProductByID(productID)
	if data.ID < 1 {
		return ""
	}
	return data.Title
}

// GetProductOtherOptionsName 获取商品附加选线的名称
func GetProductOtherOptionsName(productID int64, key string) (optionName string) {
	productData := getProductByID(productID)
	if productData.ID < 1 {
		return
	}
	var sort1, sort2 = 0, 0
	for _, v := range productData.OtherOptions.DataList {
		if v.Key == key {
			sort1 = v.Sort1
			sort2 = v.Sort2
			break
		}
	}
	for k, v := range productData.OtherOptions.Sort1.Options {
		if k == sort1 {
			optionName = v
			break
		}
	}
	for k, v := range productData.OtherOptions.Sort2.Options {
		if k == sort2 {
			optionName = fmt.Sprint(optionName, ",", v)
			break
		}
	}
	return
}

// GetProductLastPrice 获取商品最终价格
func GetProductLastPrice(productID int64, optionKey string) (price int64) {
	productData := getProductByID(productID)
	if productData.ID < 1 {
		return
	}
	if optionKey == "" {
		if productData.PriceExpireAt.Unix() >= CoreFilter.GetNowTime().Unix() {
			return productData.PriceReal
		}
		return productData.Price
	}
	for _, v := range productData.OtherOptions.DataList {
		if v.Key == optionKey {
			if v.PriceExpireAt.Unix() >= CoreFilter.GetNowTime().Unix() {
				return v.PriceReal
			}
			return v.Price
		}
	}
	return
}

// ArgsGetProductsByTickets 通过票据列货物商品列参数
type ArgsGetProductsByTickets struct {
	//票据ID列
	TicketIDs pq.Int64Array `db:"ticket_ids" json:"ticketIDs" check:"ids"`
}

type DataGetProductsByTickets struct {
	//商品ID
	ID int64 `json:"id"`
	//商品名称
	Name string `json:"name"`
	//票据配置id列
	TicketIDs pq.Int64Array `json:"ticketIDs"`
}

// GetProductsByTickets 通过票据列货物商品列
func GetProductsByTickets(args *ArgsGetProductsByTickets) (dataList []DataGetProductsByTickets, err error) {
	var rawList []FieldsCore
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM mall_core WHERE user_ticket @> $1 AND publish_at > to_timestamp(1000000) AND delete_at < to_timestamp(1000000) AND product_type = 0 AND parent_id = 0 LIMIT 10", args.TicketIDs)
	if err != nil || len(rawList) < 1 {
		err = errors.New("data is empty")
		return
	}
	for _, v := range rawList {
		vData := getProductByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, DataGetProductsByTickets{
			ID:        vData.ID,
			Name:      vData.Title,
			TicketIDs: vData.UserTicket,
		})
	}
	return
}

// ArgsGetProductCountBySort 获取指定分类下的商品数量参数
type ArgsGetProductCountBySort struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id"`
}

// GetProductCountBySort 获取指定分类下的商品数量
func GetProductCountBySort(args *ArgsGetProductCountBySort) (count int64) {
	cacheMark := getProductSortCountCacheMark(args.SortID)
	var err error
	count, err = Router2SystemConfig.MainCache.GetInt64(cacheMark)
	if err == nil && count > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT count(id) FROM mall_core WHERE org_id = $1 AND sort_id = $2 AND delete_at < to_timestamp(1000000) AND publish_at > to_timestamp(1000000)", args.OrgID, args.SortID)
	if count < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetInt64(cacheMark, count, 86400)
	return
}

// GetProductListByWarehouseProductID 通过仓储产品ID获取商品
func GetProductListByWarehouseProductID(warehouseProductID int64) (dataList []FieldsCore) {
	var rawList []FieldsCore
	err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM mall_core WHERE delete_at < to_timestamp(1000000) AND warehouse_product_id = $1", warehouseProductID)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getProductByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// getProductAddress 获取商品配送出发地
func getProductAddress(id int64) (data CoreSQLAddress.FieldsAddress, err error) {
	productData := getProductByID(id)
	if productData.ID < 1 {
		err = errors.New("no data")
		return
	}
	if productData.Address.Address != "" {
		//删除缓冲
		deleteProductCache(productData.ID)
		//反馈
		return productData.Address, nil
	}
	//TODO: 等待仓储模块完成后对接
	err = errors.New("not find")
	return
}

// 获取商品是否可用
func getProductByIDCanUse(id int64, orgID int64) (data FieldsCore, b bool) {
	data = getProductByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) || !CoreSQL.CheckTimeHaveData(data.PublishAt) || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		return
	}
	b = true
	return
}

// 获取商品信息
func getProductByID(id int64) (data FieldsCore) {
	cacheMark := getProductCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, parent_id, org_id, product_type, is_virtual, sort, publish_at, sort_id, tags, code, title, title_des, des, cover_file_ids, des_files, currency, price_real, price_expire_at, price, price_no_tax, integral, integral_price, integral_transport_free, user_sub_price, user_ticket, transport_id, address, warehouse_product_id, weight, count, buy_count, other_options, giving_tickets, params FROM mall_core WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 86400)
	return
}
