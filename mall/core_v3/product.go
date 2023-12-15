package MallCoreV3

import (
	"errors"
	BaseQiniu "github.com/fotomxq/weeekj_core/v5/base/qiniu"
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserSubscription "github.com/fotomxq/weeekj_core/v5/user/subscription"
	UserTicket "github.com/fotomxq/weeekj_core/v5/user/ticket"
	"time"

	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
)

// ArgsGetProductData 获取产品信息参数
type ArgsGetProductData struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

type DataGetProductData struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//商品类型
	// 0 普通商品; 1 关联选项商品
	ProductType int `db:"product_type" json:"productType" check:"intThan0" empty:"true"`
	//是否为虚拟商品
	// 不会发生配送处理
	IsVirtual bool `db:"is_virtual" json:"isVirtual" check:"bool"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织名称
	OrgName string `json:"orgName"`
	//排序
	Sort int `db:"sort" json:"sort"`
	//发布时间
	PublishAt time.Time `db:"publish_at" json:"publishAt"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//排序名称
	SortName string `json:"sortName"`
	//标签名称
	TagName map[int64]string `json:"tagName"`
	//商品条形码编码
	Code string `db:"code" json:"code"`
	//标题
	Title string `db:"title" json:"title"`
	//副标题
	TitleDes string `db:"title_des" json:"titleDes"`
	//商品描述
	Des string `db:"des" json:"des"`
	//封面
	CoverFileURLs []string `json:"coverFileURLs"`
	//描述图组
	DesFiles []string `json:"desFiles"`
	//货币
	Currency int `db:"currency" json:"currency"`
	//实际费用
	PriceReal int64 `db:"price_real" json:"priceReal"`
	//折扣截止
	PriceExpireAt time.Time `db:"price_expire_at" json:"priceExpireAt"`
	//折扣前费用
	Price int64 `db:"price" json:"price"`
	//积分价格
	Integral int64 `db:"integral" json:"integral"`
	//积分最多抵扣费用
	IntegralPrice int64 `db:"integral_price" json:"integralPrice"`
	//积分包含配送费免费
	IntegralTransportFree bool `db:"integral_transport_free" json:"integralTransportFree"`
	//会员价格
	// 会员配置分平台和商户，平台会员需参与活动才能使用，否则将禁止设置和后期使用
	UserSubPrice []DataUserSubPrice `db:"user_sub_price" json:"userSubPrice"`
	//票据
	// 可以使用的票据列，具体的配置在票据配置内进行设置
	// 票据分平台和商户，平台票据需参与活动才能使用，否则将自动禁止设置和后期使用
	UserTicket map[int64]string `db:"user_ticket" json:"userTicket"`
	//配送费计费模版ID
	TransportID int64 `db:"transport_id" json:"transportID"`
	//唯一送货地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//货物重量
	// 单位g
	Weight int `db:"weight" json:"weight"`
	//总库存
	Count int `db:"count" json:"count"`
	//销量
	BuyCount int `db:"buy_count" json:"buyCount"`
	//关联附加选项
	OtherOptions DataMargeOtherOption `db:"other_options" json:"otherOptions"`
	//给与票据列
	// 和赠礼区别在于，赠礼不可退，此票据会跟随订单取消设置是否退还
	GivingTickets []DataGivingTicket `db:"giving_tickets" json:"givingTickets"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// DataUserSubPrice 会员价格
type DataUserSubPrice struct {
	//会员ID
	ID int64 `db:"id" json:"id"`
	//会员名称
	Name string `db:"name" json:"name"`
	//标定价格
	Price int64 `db:"price" json:"price"`
}

type DataGivingTicket struct {
	//票据ID
	TicketConfigID int64 `db:"ticket_config_id" json:"ticketConfigID" check:"id"`
	//票据名称
	Name string `db:"name" json:"name"`
	//张数
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//是否可退
	IsRefund bool `db:"is_refund" json:"isRefund" check:"bool"`
}

type DataMargeOtherOption struct {
	//分类1
	Sort1 DataOtherOptionSort `db:"sort1" json:"sort1"`
	//分类2
	Sort2 DataOtherOptionSort `db:"sort2" json:"sort2"`
	//数据集合
	DataList []DataMargeOtherOptionChild `db:"data_list" json:"dataList"`
}

type DataMargeOtherOptionChild struct {
	//分类1的选项key
	Sort1 int `db:"sort1" json:"sort1" check:"intThan0" empty:"true"`
	//分类2的选项key
	Sort2 int `db:"sort2" json:"sort2" check:"intThan0" empty:"true"`
	//商品ID
	// 可以给0，则必须声明其他项目内容
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//选项key
	Key string `db:"key" json:"key" check:"mark"`
	//实际费用
	PriceReal int64 `db:"price_real" json:"priceReal" check:"price" empty:"true"`
	//折扣截止
	PriceExpireAt time.Time `db:"price_expire_at" json:"priceExpireAt" check:"isoTime" empty:"true"`
	//折扣前费用
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//封面
	CoverFileURL string `json:"coverFileURL"`
	//总库存
	Count int `db:"count" json:"count" check:"intThan0" empty:"true"`
	//商品条形码编码
	Code string `db:"code" json:"code"`
}

// GetProductData 获取产品信息
func GetProductData(args *ArgsGetProductData) (data DataGetProductData, err error) {
	var rawData FieldsProduct
	rawData, err = GetProductTop(&ArgsGetProduct{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	if rawData.DeleteAt.Unix() > 1000000 || rawData.PublishAt.Unix() < 1000000 {
		err = errors.New("no data")
		return
	}
	otherOptions := DataMargeOtherOption{
		Sort1: DataOtherOptionSort{
			Name:    rawData.OtherOptions.Sort1.Name,
			Options: rawData.OtherOptions.Sort1.Options,
		},
		Sort2: DataOtherOptionSort{
			Name:    rawData.OtherOptions.Sort2.Name,
			Options: rawData.OtherOptions.Sort2.Options,
		},
		DataList: []DataMargeOtherOptionChild{},
	}
	for k := 0; k < len(rawData.OtherOptions.DataList); k++ {
		v := rawData.OtherOptions.DataList[k]
		appendData := DataMargeOtherOptionChild{
			Sort1:         v.Sort1,
			Sort2:         v.Sort2,
			ProductID:     v.ProductID,
			Key:           v.Key,
			PriceReal:     v.PriceReal,
			PriceExpireAt: v.PriceExpireAt,
			Price:         v.Price,
			CoverFileURL:  "",
			Count:         v.Count,
			Code:          v.Code,
		}
		if v.CoverFileID > 0 {
			appendData.CoverFileURL, _ = BaseQiniu.GetPublicURLStr(v.CoverFileID)
		}
		otherOptions.DataList = append(otherOptions.DataList, appendData)
	}
	data = DataGetProductData{
		ID:                    rawData.ID,
		CreateAt:              rawData.CreateAt,
		UpdateAt:              rawData.UpdateAt,
		ProductType:           rawData.ProductType,
		IsVirtual:             rawData.IsVirtual,
		OrgID:                 rawData.OrgID,
		OrgName:               "",
		Sort:                  rawData.Sort,
		PublishAt:             rawData.PublishAt,
		SortID:                rawData.SortID,
		SortName:              "",
		TagName:               map[int64]string{},
		Code:                  rawData.Code,
		Title:                 rawData.Title,
		TitleDes:              rawData.TitleDes,
		Des:                   rawData.Des,
		CoverFileURLs:         []string{},
		DesFiles:              []string{},
		Currency:              rawData.Currency,
		PriceReal:             rawData.PriceReal,
		PriceExpireAt:         rawData.PriceExpireAt,
		Price:                 rawData.Price,
		Integral:              rawData.Integral,
		IntegralPrice:         rawData.IntegralPrice,
		IntegralTransportFree: rawData.IntegralTransportFree,
		UserSubPrice:          []DataUserSubPrice{},
		UserTicket:            map[int64]string{},
		TransportID:           rawData.TransportID,
		Address:               rawData.Address,
		Weight:                rawData.Weight,
		Count:                 rawData.Count,
		BuyCount:              rawData.BuyCount,
		OtherOptions:          otherOptions,
		GivingTickets:         []DataGivingTicket{},
		Params:                rawData.Params,
	}
	if rawData.OrgID > 0 {
		data.OrgName, _ = OrgCoreCore.GetOrgName(&OrgCoreCore.ArgsGetOrg{
			ID: rawData.OrgID,
		})
	}
	if rawData.SortID > 0 {
		data.SortName, _ = Sort.GetName(&ClassSort.ArgsGetID{
			ID:     rawData.ID,
			BindID: -1,
		})
	}
	if len(rawData.Tags) > 0 {
		data.TagName = Tags.GetByIDNames(rawData.Tags, -1, 100)
	}
	if len(rawData.CoverFileIDs) > 0 {
		data.CoverFileURLs, _ = BaseQiniu.GetPublicURLStrs(rawData.CoverFileIDs)
	}
	if len(rawData.DesFiles) > 0 {
		data.DesFiles, _ = BaseQiniu.GetPublicURLStrs(rawData.DesFiles)
	}
	if len(rawData.UserSubPrice) > 0 {
		var subIDs []int64
		for _, v := range rawData.UserSubPrice {
			subIDs = append(subIDs, v.ID)
		}
		userSubConfigs, _ := UserSubscription.GetConfigMore(&UserSubscription.ArgsGetConfigMore{
			IDs:        subIDs,
			HaveRemove: false,
		})
		if len(userSubConfigs) > 0 {
			for _, v := range rawData.UserSubPrice {
				for _, v2 := range userSubConfigs {
					if v.ID != v2.ID {
						continue
					}
					data.UserSubPrice = append(data.UserSubPrice, DataUserSubPrice{
						ID:    v2.ID,
						Name:  v2.Title,
						Price: v.Price,
					})
				}
			}
		}
	}
	if len(rawData.UserTicket) > 0 {
		data.UserTicket, _ = UserTicket.GetConfigMoreMap(&UserTicket.ArgsGetConfigMore{
			IDs:        rawData.UserTicket,
			HaveRemove: false,
		})
	}
	if len(rawData.GivingTickets) > 0 {
		var ticketIDs []int64
		for _, v := range rawData.GivingTickets {
			ticketIDs = append(ticketIDs, v.TicketConfigID)
		}
		if len(ticketIDs) > 0 {
			userTicket, _ := UserTicket.GetConfigMoreMap(&UserTicket.ArgsGetConfigMore{
				IDs:        rawData.UserTicket,
				HaveRemove: false,
			})
			for _, v := range rawData.GivingTickets {
				for k2, v2 := range userTicket {
					if v.TicketConfigID != k2 {
						continue
					}
					data.GivingTickets = append(data.GivingTickets, DataGivingTicket{
						TicketConfigID: v.TicketConfigID,
						Name:           v2,
						Count:          v.Count,
						IsRefund:       v.IsRefund,
					})
				}
			}
		}
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
func GetProduct(args *ArgsGetProduct) (data FieldsProduct, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, parent_id, org_id, product_type, is_virtual, sort, publish_at, sort_id, tags, code, title, title_des, des, cover_file_ids, des_files, currency, price_real, price_expire_at, price, integral, integral_price, integral_transport_free, user_sub_price, user_ticket, transport_id, address, warehouse_product_id, weight, count, buy_count, other_options, giving_tickets, params FROM mall_core WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// GetProductTop 是否采用追溯的商品获取产品
func GetProductTop(args *ArgsGetProduct) (data FieldsProduct, err error) {
	data, err = GetProduct(args)
	if err != nil {
		return
	}
	if data.ParentID > 0 {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, parent_id, org_id, product_type, is_virtual, sort, publish_at, sort_id, tags, code, title, title_des, des, cover_file_ids, des_files, currency, price_real, price_expire_at, price, integral, integral_price, integral_transport_free, user_sub_price, user_ticket, transport_id, address, warehouse_product_id, weight, count, buy_count, other_options, giving_tickets, params FROM mall_core WHERE parent_id = $1", data.ParentID)
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
func GetProducts(args *ArgsGetProducts) (dataList []FieldsProduct, err error) {
	err = CoreSQLIDs.GetIDsOrgAndDelete(&dataList, "mall_core", "id, create_at, update_at, delete_at, parent_id, org_id, product_type, is_virtual, sort, publish_at, sort_id, tags, code, title, title_des, des, cover_file_ids, des_files, currency, price_real, price_expire_at, price, integral, integral_price, integral_transport_free, user_sub_price, user_ticket, transport_id, address, warehouse_product_id, weight, count, buy_count, other_options, giving_tickets, params", args.IDs, args.OrgID, args.HaveRemove)
	var newDataList []FieldsProduct
	for _, v := range dataList {
		if v.ParentID < 1 {
			newDataList = append(newDataList, v)
			continue
		}
		var data FieldsProduct
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, parent_id, org_id, product_type, is_virtual, sort, publish_at, sort_id, tags, code, title, title_des, des, cover_file_ids, des_files, currency, price_real, price_expire_at, price, integral, integral_price, integral_transport_free, user_sub_price, user_ticket, transport_id, address, warehouse_product_id, weight, count, buy_count, other_options, giving_tickets, params FROM mall_core WHERE parent_id = $1 AND ($2 = true OR delete_at < to_timestamp(1000000))", data.ParentID, args.HaveRemove)
		if err != nil {
			err = nil
			continue
		}
		newDataList = append(newDataList, data)
	}
	if len(newDataList) < 1 {
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
