package MallCore

import (
	"fmt"
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	ServiceOrderMod "github.com/fotomxq/weeekj_core/v5/service/order/mod"
	UserTicket "github.com/fotomxq/weeekj_core/v5/user/ticket"
	"github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
	"strings"
)

func subNats() {
	//请求赠送用户虚拟商品
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "商城核心请求赠送虚拟商品",
		Description:  "",
		EventSubType: "all",
		Code:         "mall_core_product_virtual",
		EventType:    "nats",
		EventURL:     "/mall/core/product_virtual",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("mall_core_product_virtual", "/mall/core/product_virtual", subNatsProductVirtual)
	//请求同步修改商品信息
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "商城核心请求修改商品信息",
		Description:  "",
		EventSubType: "all",
		Code:         "mall_core_product_update",
		EventType:    "nats",
		EventURL:     "/mall/core/product_update",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("mall_core_product_update", "/mall/core/product_update", subNatsUpdateProduct)
}

// 请求赠送用户虚拟商品
func subNatsProductVirtual(_ *nats.Msg, action string, id int64, _ string, data []byte) {
	switch action {
	case "send":
		//给用户发送虚拟商品
		//获取参数
		orderID := gjson.GetBytes(data, "orderID").Int()
		count := gjson.GetBytes(data, "count").Int()
		userID := gjson.GetBytes(data, "userID").Int()
		orgID := gjson.GetBytes(data, "orgID").Int()
		//获取商品
		productData := getProductByID(id)
		//检查关系
		if !CoreFilter.EqID2(orgID, productData.OrgID) {
			ServiceOrderMod.AddLog(orderID, "订单所属虚拟商品的不属于该组织")
			return
		}
		if count < 1 {
			ServiceOrderMod.AddLog(orderID, "订单购买的虚拟商品数量少于0，自动跳过")
			return
		}
		//检查商品是否为虚拟商品
		if !productData.IsVirtual {
			ServiceOrderMod.AddLog(orderID, "订单购买的商品不是虚拟商品")
			return
		}
		//检查赠送单元
		//已经赠送的内容
		newTickets := UserTicket.ArgsAddTickets{
			OrgID:  orgID,
			UserID: userID,
			Data:   []UserTicket.ArgsAddTicketsChild{},
		}
		//可退票据列
		var canRefundConfigIDs []int64
		//处理票据赠予环节
		if len(productData.GivingTickets) > 0 {
			for _, vTicket := range productData.GivingTickets {
				newTickets.Data = append(newTickets.Data, UserTicket.ArgsAddTicketsChild{
					ConfigID:    vTicket.TicketConfigID,
					Count:       vTicket.Count,
					UseFromName: fmt.Sprint("购买商品(", productData.Title, ")[", productData.ID, "]"),
				})
				if vTicket.IsRefund {
					canRefundConfigIDs = append(canRefundConfigIDs, vTicket.TicketConfigID)
				}
			}
		}
		//为用户新增票据
		if len(newTickets.Data) > 0 {
			newTickets.CanRefundConfigIDs = canRefundConfigIDs
			newTicketIDs, newTicketRefundIDs, err := UserTicket.AddTickets(&newTickets)
			if err != nil {
				ServiceOrderMod.AddLog(orderID, fmt.Sprint("虚拟商品给与票据失败，错误代码: ", err))
				return
			}
			if len(newTicketIDs) > 0 {
				var newTicketIDsStr []string
				var newTicketIDsRefundStr []string
				for _, v := range newTicketIDs {
					newTicketIDsStr = append(newTicketIDsStr, CoreFilter.GetStringByInt64(v))
				}
				for _, v := range newTicketRefundIDs {
					newTicketIDsRefundStr = append(newTicketIDsRefundStr, CoreFilter.GetStringByInt64(v))
				}
				//记录订单扩展参数
				orderParams := []CoreSQLConfig.FieldsConfigType{
					{
						Mark: "mall_virtual_user_tickets_can_refund",
						Val:  strings.Join(newTicketIDsRefundStr, ","),
					},
				}
				ServiceOrderMod.UpdateOrderParams(orderID, orderParams)
			}
		}
		//标记订单完成
		ServiceOrderMod.UpdateFinish(orderID, "虚拟商品已经给与用户，自动完成订单")
	}
}

// argsSubNatsUpdateProduct 修改商品信息参数
type argsSubNatsUpdateProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
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
	//货物重量
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//建议售价
	Price int64 `json:"price"`
	//不含税价格
	PriceNoTax int64 `db:"price_no_tax" json:"priceNoTax" check:"price" empty:"true"`
}

// subNatsUpdateProduct 修改商品信息
func subNatsUpdateProduct(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	appendLog := "mall core sub nats update product, "
	var args argsSubNatsUpdateProduct
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Error(appendLog, "get params, ", err)
		return
	}
	mallData := getProductByID(args.ID)
	if mallData.ID < 1 {
		CoreLog.Error(appendLog, "get mall data, id: ", args.ID)
		return
	}
	otherOptions, err := mallData.OtherOptions.GetData()
	if err != nil {
		CoreLog.Error(appendLog, "get mall other options, mall product id: ", args.ID, ", err: ", err)
		return
	}
	errCode, err := UpdateProduct(&ArgsUpdateProduct{
		ID:                    mallData.ID,
		OrgID:                 mallData.OrgID,
		IsVirtual:             mallData.IsVirtual,
		Sort:                  mallData.Sort,
		SortID:                mallData.SortID,
		Tags:                  mallData.Tags,
		Code:                  mallData.Code,
		Title:                 args.Title,
		TitleDes:              args.TitleDes,
		Des:                   args.Des,
		CoverFileIDs:          args.CoverFileIDs,
		DesFiles:              args.DesFiles,
		Currency:              mallData.Currency,
		PriceReal:             mallData.Price,
		PriceExpireAt:         CoreFilter.GetISOByTime(mallData.PriceExpireAt),
		Price:                 args.Price,
		PriceNoTax:            args.PriceNoTax,
		Integral:              mallData.Integral,
		IntegralPrice:         mallData.IntegralPrice,
		IntegralTransportFree: mallData.IntegralTransportFree,
		UserSubPrice:          mallData.UserSubPrice,
		UserTicket:            mallData.UserTicket,
		TransportID:           mallData.TransportID,
		Address:               mallData.Address,
		WarehouseProductID:    mallData.WarehouseProductID,
		Weight:                args.Weight,
		Count:                 mallData.Count,
		OtherOptions:          otherOptions,
		GivingTickets:         mallData.GivingTickets,
		Params:                mallData.Params,
		SyncERPProduct:        false,
	})
	if err != nil {
		CoreLog.Error(appendLog, "update product, ", errCode, ", err: ", err)
		return
	}
}
