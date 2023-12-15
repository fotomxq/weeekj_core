package MallCore

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	ERPProductMod "github.com/fotomxq/weeekj_core/v5/erp/product/mod"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsUpdateProduct 修改商品参数
type ArgsUpdateProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否为虚拟商品
	// 不会发生配送处理
	IsVirtual bool `db:"is_virtual" json:"isVirtual" check:"bool"`
	//排序
	Sort int `db:"sort" json:"sort"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//商品条形码编码
	Code string `db:"code" json:"code"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//副标题
	TitleDes string `db:"title_des" json:"titleDes" check:"title" min:"1" max:"300" empty:"true"`
	//商品描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//封面ID
	CoverFileIDs pq.Int64Array `db:"cover_file_ids" json:"coverFileIDs" check:"ids" empty:"true"`
	//描述图组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	// 货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	// 实际费用
	PriceReal int64 `db:"price_real" json:"priceReal" check:"price"`
	//折扣截止
	PriceExpireAt string `db:"price_expire_at" json:"priceExpireAt" check:"isoTime" empty:"true"`
	//折扣前费用
	Price int64 `db:"price" json:"price" check:"price"`
	//不含税价格
	PriceNoTax int64 `db:"price_no_tax" json:"priceNoTax" check:"price" empty:"true"`
	//积分价格
	Integral int64 `db:"integral" json:"integral" check:"price" empty:"true"`
	//积分最多抵扣费用
	IntegralPrice int64 `db:"integral_price" json:"integralPrice" check:"price" empty:"true"`
	//积分包含配送费免费
	IntegralTransportFree bool `db:"integral_transport_free" json:"integralTransportFree"`
	//会员价格
	// 会员配置分平台和商户，平台会员需参与活动才能使用，否则将禁止设置和后期使用
	UserSubPrice FieldsUserSubPrices `db:"user_sub_price" json:"userSubPrice"`
	//票据
	// 可以使用的票据列，具体的配置在票据配置内进行设置
	// 票据分平台和商户，平台票据需参与活动才能使用，否则将自动禁止设置和后期使用
	UserTicket pq.Int64Array `db:"user_ticket" json:"userTicket" check:"ids" empty:"true"`
	//配送费计费模版ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id" empty:"true"`
	//唯一送货地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//关联仓库的产品
	WarehouseProductID int64 `db:"warehouse_product_id" json:"warehouseProductID" check:"id" empty:"true"`
	//货物重量
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//总库存
	Count int `db:"count" json:"count" check:"intThan0" empty:"true"`
	//关联附加选项
	OtherOptions DataOtherOptions `db:"other_options" json:"otherOptions"`
	//给与票据列
	// 和赠礼区别在于，赠礼不可退，此票据会跟随订单取消设置是否退还
	GivingTickets FieldsGivingTickets `db:"giving_tickets" json:"givingTickets"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//是否同步ERP产品
	SyncERPProduct bool `json:"syncERPProduct" check:"bool"`
}

// UpdateProduct 修改商品
func UpdateProduct(args *ArgsUpdateProduct) (errCode string, err error) {
	//处理时间
	var priceExpireAt time.Time
	if args.PriceExpireAt != "" {
		priceExpireAt, err = CoreFilter.GetTimeByISO(args.PriceExpireAt)
		if err != nil {
			errCode = "price_expire_at"
			return
		}
	}
	//计算其他选项
	var otherOptions FieldsOtherOption
	otherOptions, err = args.OtherOptions.GetFields()
	if err != nil {
		errCode = "other_expire_at"
		return
	}
	if len(otherOptions.DataList) > 0 {
		var mallProductOtherOptionsMax int64
		mallProductOtherOptionsMax, err = BaseConfig.GetDataInt64("MallProductOtherOptionsMax")
		if err != nil {
			mallProductOtherOptionsMax = 0
		}
		if mallProductOtherOptionsMax > 0 {
			if int64(len(otherOptions.DataList)) > mallProductOtherOptionsMax {
				errCode = "other_option_too_many"
				err = errors.New("other option too many")
				return
			}
		}
	}
	//获取原始数据
	var oldData FieldsCore
	oldData, err = GetProduct(&ArgsGetProduct{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		errCode = "not_exist"
		err = errors.New("product not exist")
		return
	}
	//检查当前数据是否为上级了
	if oldData.ParentID > 0 {
		errCode = "is_history"
		err = errors.New("data have parent")
		return
	}
	//修改商品并回到等待提交状态
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE mall_core SET update_at = NOW(), is_virtual = :is_virtual, sort = :sort, publish_at = TO_TIMESTAMP(0), sort_id = :sort_id, tags = tags, code = :code, title = :title, title_des = :title_des, des = :des, cover_file_ids = :cover_file_ids, des_files = :des_files, currency = :currency, price_real = :price_real, price_expire_at = :price_expire_at, price = :price, price_no_tax = :price_no_tax, integral = :integral, integral_price = :integral_price, integral_transport_free = :integral_transport_free, user_sub_price = :user_sub_price, user_ticket = :user_ticket, transport_id = :transport_id, address = :address, warehouse_product_id = :warehouse_product_id, weight = :weight, count = :count, other_options = :other_options, giving_tickets = :giving_tickets, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"is_virtual":              args.IsVirtual,
		"sort":                    args.Sort,
		"sort_id":                 args.SortID,
		"tags":                    args.Tags,
		"code":                    args.Code,
		"title":                   args.Title,
		"title_des":               args.TitleDes,
		"des":                     args.Des,
		"cover_file_ids":          args.CoverFileIDs,
		"des_files":               args.DesFiles,
		"currency":                args.Currency,
		"price_real":              args.PriceReal,
		"price_expire_at":         priceExpireAt,
		"price":                   args.Price,
		"price_no_tax":            args.PriceNoTax,
		"integral":                args.Integral,
		"integral_price":          args.IntegralPrice,
		"integral_transport_free": args.IntegralTransportFree,
		"user_sub_price":          args.UserSubPrice,
		"user_ticket":             args.UserTicket,
		"transport_id":            args.TransportID,
		"address":                 args.Address,
		"warehouse_product_id":    args.WarehouseProductID,
		"weight":                  args.Weight,
		"count":                   args.Count,
		"other_options":           otherOptions,
		"giving_tickets":          args.GivingTickets,
		"params":                  args.Params,
		"id":                      args.ID,
		"org_id":                  args.OrgID,
	})
	if err == nil {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO mall_core (create_at, update_at, sort, parent_id, org_id, product_type, is_virtual, sort_id, tags, code, title, title_des, des, cover_file_ids, des_files, currency, price_real, price_expire_at, price, price_no_tax, integral, integral_price, integral_transport_free, user_sub_price, user_ticket, transport_id, address, warehouse_product_id, weight, count, buy_count, other_options, giving_tickets, params) VALUES (:create_at, :update_at, :sort, :parent_id, :org_id, :product_type, :is_virtual, :sort_id, :tags, :code, :title, :title_des, :des, :cover_file_ids, :des_files, :currency, :price_real, :price_expire_at, :price, :price_no_tax, :integral, :integral_price, :integral_transport_free, :user_sub_price, :user_ticket, :transport_id, :address, :warehouse_product_id, :weight, :count, :buy_count, :other_options, :giving_tickets, :params)", map[string]interface{}{
			"create_at":               oldData.CreateAt,
			"update_at":               CoreFilter.GetNowTime(),
			"sort":                    oldData.Sort,
			"parent_id":               oldData.ID,
			"org_id":                  oldData.OrgID,
			"product_type":            oldData.ProductType,
			"is_virtual":              oldData.IsVirtual,
			"sort_id":                 oldData.SortID,
			"tags":                    oldData.Tags,
			"code":                    args.Code,
			"title":                   oldData.Title,
			"title_des":               oldData.TitleDes,
			"des":                     oldData.Des,
			"cover_file_ids":          oldData.CoverFileIDs,
			"des_files":               oldData.DesFiles,
			"currency":                oldData.Currency,
			"price_real":              oldData.PriceReal,
			"price_expire_at":         oldData.PriceExpireAt,
			"price":                   oldData.Price,
			"price_no_tax":            oldData.PriceNoTax,
			"integral":                oldData.Integral,
			"integral_price":          oldData.IntegralPrice,
			"integral_transport_free": oldData.IntegralTransportFree,
			"user_sub_price":          oldData.UserSubPrice,
			"user_ticket":             oldData.UserTicket,
			"transport_id":            oldData.TransportID,
			"address":                 oldData.Address,
			"warehouse_product_id":    oldData.WarehouseProductID,
			"weight":                  oldData.Weight,
			"count":                   oldData.Count,
			"buy_count":               oldData.BuyCount,
			"other_options":           oldData.OtherOptions,
			"giving_tickets":          oldData.GivingTickets,
			"params":                  oldData.Params,
		})
		if err != nil {
			errCode = "insert"
			return
		}
	} else {
		errCode = "update"
		return
	}
	//删除缓冲
	deleteProductCache(args.ID)
	if args.SortID > 0 {
		deleteProductSortCache(args.SortID)
	}
	if oldData.SortID > 0 {
		deleteProductSortCache(oldData.SortID)
	}
	//同步ERP产品
	if args.SyncERPProduct {
		newData := getProductByID(args.ID)
		if newData.ID > 0 && newData.WarehouseProductID > 0 {
			ERPProductMod.UpdateProduct(ERPProductMod.ArgsUpdateProduct{
				ID:           newData.WarehouseProductID,
				Title:        newData.Title,
				TitleDes:     newData.TitleDes,
				Des:          newData.Des,
				CoverFileIDs: newData.CoverFileIDs,
				Weight:       newData.Weight,
				Price:        newData.Price,
			})
		}
	}
	//反馈
	return
}

// ArgsUpdateProductCount 修改库存参数
type ArgsUpdateProductCount struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//库存数量
	Count int `db:"count" json:"count" check:"intThan0" empty:"true"`
}

// UpdateProductCount 修改库存
func UpdateProductCount(args *ArgsUpdateProductCount) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE mall_core SET update_at = NOW(), count = :count WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteProductCache(args.ID)
	//反馈
	return
}

// UpdateProductAddCount 修改库存
func UpdateProductAddCount(productID int64, optionKey string, count int) (err error) {
	//获取商品
	productData := getProductByID(productID)
	//如果存在选项
	if len(productData.OtherOptions.DataList) > 0 {
		for k, v := range productData.OtherOptions.DataList {
			if v.Key != optionKey {
				continue
			}
			productData.OtherOptions.DataList[k].Count += count
			if productData.OtherOptions.DataList[k].Count < 0 {
				err = errors.New("no count")
				return
			}
		}
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE mall_core SET update_at = NOW(), other_options = :other_options WHERE id = :id", map[string]interface{}{
			"id":            productID,
			"other_options": productData.OtherOptions,
		})
		if err != nil {
			return
		}
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE mall_core SET update_at = NOW(), count = count + :count WHERE id = :id AND 0 <= count + :count", map[string]interface{}{
			"id":    productID,
			"count": count,
		})
		if err != nil {
			return
		}
	}
	//删除缓冲
	deleteProductCache(productID)
	//反馈
	return
}

// ArgsUpdateProductPublish 发布商品参数
type argsUpdateProductPublish struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// UpdateProductPublish 发布商品
func updateProductPublish(args *argsUpdateProductPublish) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE mall_core SET update_at = NOW(), publish_at = NOW() WHERE id = :id", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteProductCache(args.ID)
	//反馈
	return
}

// UpdateProductPublishDown 下架产品
func UpdateProductPublishDown(orgID int64, productID int64) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE mall_core SET update_at = NOW(), publish_at = to_timestamp(0) WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":     productID,
		"org_id": orgID,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteProductCache(productID)
	//反馈
	return
}

func UpdateProductBuy(productID int64, count int) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE mall_core SET update_at = NOW(), buy_count = buy_count + :count WHERE id = :id", map[string]interface{}{
		"id":    productID,
		"count": count,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update product buy count, id: ", productID, ", count: ", count, ", err: ", err))
		return
	}
	//删除缓冲
	deleteProductCache(productID)
	//反馈
	return
}
