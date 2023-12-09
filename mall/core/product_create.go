package MallCore

import (
	"errors"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserTicket "gitee.com/weeekj/weeekj_core/v5/user/ticket"
	"github.com/lib/pq"
	"time"
)

// ArgsCreateProduct 创建商品参数
type ArgsCreateProduct struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//商品类型
	// 0 普通商品; 1 关联选项商品; 2 虚拟商品
	ProductType int `db:"product_type" json:"productType" check:"intThan0" empty:"true"`
	//是否为虚拟商品
	// 不会发生配送处理
	IsVirtual bool `db:"is_virtual" json:"isVirtual" check:"bool"`
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
}

// CreateProduct 创建商品
func CreateProduct(args *ArgsCreateProduct) (data FieldsCore, errCode string, err error) {
	//修正参数
	if args.Tags == nil {
		args.Tags = []int64{}
	}
	if args.CoverFileIDs == nil {
		args.CoverFileIDs = []int64{}
	}
	if args.DesFiles == nil {
		args.DesFiles = []int64{}
	}
	if args.UserSubPrice == nil || len(args.UserSubPrice) < 1 {
		args.UserSubPrice = FieldsUserSubPrices{}
	}
	if args.UserTicket == nil {
		args.UserTicket = []int64{}
	}
	if args.GivingTickets == nil || len(args.GivingTickets) < 1 {
		args.GivingTickets = FieldsGivingTickets{}
	}
	if args.Params == nil || len(args.Params) < 1 {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
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
	//检查封面个数
	if len(args.CoverFileIDs) > 0 {
		if len(args.CoverFileIDs) > 5 {
			errCode = "cover_too_many"
			err = errors.New("cover too many")
			return
		}
	}
	//检查价格
	if args.PriceReal < 0 || args.Price < 0 {
		errCode = "price"
		err = errors.New("price real or price less 0")
		return
	}
	for _, vSub := range args.UserSubPrice {
		if vSub.Price < 0 {
			errCode = "sub_price"
			err = errors.New("sub price less 0")
			return
		}
	}
	//汇总会员ID列
	var userSubIDs pq.Int64Array
	for _, v := range args.UserSubPrice {
		userSubIDs = append(userSubIDs, v.ID)
	}
	//检查会员是否可用
	errCode, err = CheckProductSub(&ArgsCheckProductSub{
		ID:         0,
		OrgID:      data.OrgID,
		SortID:     args.SortID,
		Tags:       args.Tags,
		UserSubs:   userSubIDs,
		UserTicket: args.UserTicket,
		BuyCount:   0,
	})
	if err != nil {
		return
	}
	//检查票据是否可用
	errCode, err = CheckProductTicket(&ArgsCheckProductTicket{
		ID:         0,
		OrgID:      data.OrgID,
		SortID:     args.SortID,
		Tags:       args.Tags,
		UserSubs:   userSubIDs,
		UserTicket: args.UserTicket,
		BuyCount:   0,
	})
	if err != nil {
		return
	}
	//同一个组织下，禁止连续创建相同标题内容
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM mall_core WHERE org_id = $1 AND parent_id = 0 AND title = $2 AND create_at >= $3", args.OrgID, args.Title, CoreFilter.GetNowTimeCarbon().SubSeconds(5).Time)
	if err == nil && data.ID > 0 {
		errCode = "repeat_product"
		err = errors.New("repeat create product")
		return
	}
	//检查其他商品是否为选项类商品和本商户产品
	for _, v := range args.OtherOptions.DataList {
		if v.ProductID < 1 {
			continue
		}
		var vID int64
		err = Router2SystemConfig.MainDB.Get(&vID, "SELECT id FROM mall_core WHERE id = $1 AND delete_at < to_timestamp(1000000) AND publish_at > to_timestamp(1000000) AND product_type = 1 AND org_id = $2", v.ProductID, args.OrgID)
		if err != nil || vID < 1 {
			errCode = "product_not_option"
			err = errors.New("product not option")
			return
		}
	}
	//检查票据是否可用于该组织
	for _, v := range args.GivingTickets {
		if v.TicketConfigID < 1 {
			continue
		}
		err = UserTicket.CheckConfigOrg(args.OrgID, v.TicketConfigID)
		if err != nil {
			errCode = "ticket_not_org_self"
			err = errors.New("ticket not org self")
			return
		}
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "mall_core", "INSERT INTO mall_core (parent_id, org_id, product_type, is_virtual, sort, sort_id, tags, code, title, title_des, des, cover_file_ids, des_files, currency, price_real, price_expire_at, price, price_no_tax, integral, integral_price, integral_transport_free, user_sub_price, user_ticket, transport_id, address, warehouse_product_id, weight, count, buy_count, other_options, giving_tickets, params) VALUES (0, :org_id, :product_type, :is_virtual, 9999999, :sort_id, :tags, :code, :title, :title_des, :des, :cover_file_ids, :des_files, :currency, :price_real, :price_expire_at, :price, :price_no_tax, :integral, :integral_price, :integral_transport_free, :user_sub_price, :user_ticket, :transport_id, :address, :warehouse_product_id, :weight, :count, 0, :other_options, :giving_tickets, :params)", map[string]interface{}{
		"org_id":                  args.OrgID,
		"product_type":            args.ProductType,
		"is_virtual":              args.IsVirtual,
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
	}, &data)
	if err != nil {
		errCode = "insert"
		return
	}
	//清理缓冲
	if args.SortID > 0 {
		deleteProductSortCache(args.SortID)
	}
	//反馈
	return
}
