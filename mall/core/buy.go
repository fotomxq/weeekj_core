package MallCore

import (
	"errors"
	"fmt"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLGPS "github.com/fotomxq/weeekj_core/v5/core/sql/gps"
	MallLogMod "github.com/fotomxq/weeekj_core/v5/mall/log/mod"
	MapArea "github.com/fotomxq/weeekj_core/v5/map/area"
	MapMathArgs "github.com/fotomxq/weeekj_core/v5/map/math/args"
	MapMathPoint "github.com/fotomxq/weeekj_core/v5/map/math/point"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	OrgMap "github.com/fotomxq/weeekj_core/v5/org/map"
	ServiceOrderWaitFields "github.com/fotomxq/weeekj_core/v5/service/order/wait_fields"
	TMSUserRunning "github.com/fotomxq/weeekj_core/v5/tms/user_running"
	UserIntegral "github.com/fotomxq/weeekj_core/v5/user/integral"
	UserSubscription "github.com/fotomxq/weeekj_core/v5/user/subscription"
	UserTicket "github.com/fotomxq/weeekj_core/v5/user/ticket"
	"github.com/lib/pq"
	"strconv"
	"strings"
)

//购物前后的预备处理模块
// 1\验证可用性
// 2\处理交易开始，减少库存部分

// ArgsGetProductPrice 验证可用性并计算最终价格参数
type ArgsGetProductPrice struct {
	//商品ID列
	Products []ArgsGetProductPriceProduct `db:"products" json:"products"`
	//商户ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//会员配置ID
	// 只能指定一个
	UserSubID int64 `db:"user_sub_id" json:"userSubID" check:"id" empty:"true"`
	//票据
	// 可以使用的票据列，具体的配置在票据配置内进行设置
	// 票据分平台和商户，平台票据需参与活动才能使用，否则将自动禁止设置和后期使用
	UserTicket pq.Int64Array `db:"user_ticket" json:"userTicket" check:"ids" empty:"true"`
	//是否使用积分
	UseIntegral bool `db:"use_integral" json:"useIntegral" check:"bool"`
	//收货地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//是否绕过对库存限制
	SkipProductCountLimit bool `json:"skipProductCountLimit" check:"bool" empty:"true"`
	//配送方式
	// 0 self 自运营服务; 1 自提; 2 running 跑腿服务; 3 housekeeping 家政服务
	TransportType int `db:"transport_type" json:"transportType"`
}

type ArgsGetProductPriceProduct struct {
	//商品ID
	ID int64 `db:"id" json:"id" check:"id"`
	//选项Key
	// 如果给空，则该商品必须也不包含选项
	OptionKey string `db:"option_key" json:"optionKey" check:"mark" empty:"true"`
	//购买数量
	// 如果为0，则只判断单价的价格
	BuyCount int `db:"buy_count" json:"buyCount" check:"int64Than0"`
}

// DataProductPrice 外部用的计算汇总模块
// 拆分不同商品及优惠，并减免对应的费用
// 仅保留必要的信息结构，构建的外部专用信息模块
type DataProductPrice struct {
	//商品结构
	ProductList []FieldsCore `json:"productList"`
	//分类结构列
	SortList []ClassSort.FieldsSort `json:"sortList"`
	//原始价格
	BeforePrice int64 `json:"beforePrice"`
	//最终价格
	LastPrice int64 `json:"lastPrice"`
	//商品最终价格
	LastPriceProduct int64 `json:"lastPriceProduct"`
	//可用于会员的价格部分
	// 仅包含分类通用会员折扣部分，商品强制指定费用不在此处计算
	LastPriceSub int64 `json:"lastPriceSub"`
	//积分影响的费用
	LastPriceIntegral int64 `json:"lastPriceIntegral"`
	//票据影响的价格部分
	LastPriceTicket int64 `json:"lastPriceTicket"`
	//总的配送费用
	LastPriceTransport int64 `json:"lastPriceTransport"`
	//可用于会员的价格部分减免前
	BeforePriceSub int64 `json:"beforePriceSub"`
	//货物清单
	Goods ServiceOrderWaitFields.FieldsGoods `db:"goods" json:"goods"`
	//订单总的抵扣
	// 例如满减活动，不局限于个别商品的活动
	Exemptions ServiceOrderWaitFields.FieldsExemptions `db:"exemptions" json:"exemptions"`
	//无法享受会员的商品列
	NotSubProducts []int64 `json:"notSubProducts"`
	//用户订阅结构
	UserSubConfig UserSubscription.FieldsConfig `json:"userSubConfig"`
	//票据配置列
	TicketConfigs []UserTicket.FieldsConfig `json:"ticketConfigs"`
	//票据对应关系
	Tickets []DataProductPriceTicket `json:"tickets"`
	//用户最初积分
	BeforeUserIntegral int64 `json:"beforeUserIntegral"`
	//用户剩余积分
	LastUserIntegral int64 `json:"lastUserIntegral"`
	//配送方式
	TransportSystem string `db:"transport_system" json:"transportSystem"`
	//配送模版及费用计算
	Transports []DataProductPriceTransport `json:"transports"`
	//特殊信息结构体
	Infos []DataProductPriceInfo `json:"infos"`
}

// DataProductPriceInfo 信息结构体
type DataProductPriceInfo struct {
	//信息标识码
	Mark string `json:"mark"`
	//关联的商品
	ProductID int64 `json:"productID"`
	//是否为隐藏项
	IsHide bool `json:"isHide"`
	//影响的费用
	Price int64 `json:"price"`
	//消息
	Msg string `json:"msg"`
}

// DataProductPriceTicket 票据对应多个商品设计
type DataProductPriceTicket struct {
	//票据ID
	TicketID int64 `json:"ticketID"`
	//可用商品
	Products []int64 `json:"products"`
	//该票据存在于分类列
	SortIDs []int64 `json:"sortIDs"`
	//是否用于订单整体抵扣
	UseOrder bool `json:"useOrder"`
	//预计使用数量
	NeedCount int `json:"needCount"`
	//用户持有总量
	UserCount int `json:"userCount"`
	//商品价格
	Price int64 `json:"price"`
}

// DataProductPriceTransport 配送费用计算规则
type DataProductPriceTransport struct {
	//计费模版结构
	ConfigData FieldsTransport `json:"configData"`
	//产品ID
	// 免费产品部分将直接跳过
	ProductIDs []int64 `json:"productIDs"`
	//总重量
	Weight int `json:"weight"`
	//总件数
	BuyCount int `json:"buyCount"`
	//配送距离
	DistanceM float64 `json:"distanceM"`
	//费用合计
	Price int64 `json:"price"`
}

// GetProductPrice 验证可用性并计算最终价格
// 只能计算同一个商户下的数据
func GetProductPrice(args *ArgsGetProductPrice) (resultData DataProductPrice, errCode string, err error) {
	/**
	模块计算规则
	1、获取配置包
	2、获取商户设置
	3、遍历和处理商品
	4、继续处理订单总的折扣
	5、计算总的配送费
	*/
	//初始化数据集合
	resultData.Exemptions = ServiceOrderWaitFields.FieldsExemptions{}
	//拦截没有货物的请求
	if len(args.Products) < 1 {
		errCode = "err_buy_empty"
		err = errors.New("goods is empty")
		return
	}
	//////////////////////////////////////////////////////////////////////////////////////////////
	//1、获取配置包
	//////////////////////////////////////////////////////////////////////////////////////////////
	//获取用户的积分
	var userIntegral UserIntegral.FieldsIntegral
	userIntegral, err = UserIntegral.GetUser(&UserIntegral.ArgsGetUser{
		OrgID:  args.OrgID,
		UserID: args.UserID,
	})
	if err != nil {
		//如果积分不存在，则自动计算为0
		err = nil
	}
	resultData.BeforeUserIntegral = userIntegral.Count + 0
	resultData.LastUserIntegral = userIntegral.Count + 0
	//获取用户订阅配置
	resultData.UserSubConfig, err = UserSubscription.GetConfigByID(&UserSubscription.ArgsGetConfigByID{
		ID:    args.UserSubID,
		OrgID: args.OrgID,
	})
	if err != nil {
		err = nil
		//找不到则忽略后续处理，将自动跳过会员部分
	}
	//获取票据配置列
	resultData.TicketConfigs, err = UserTicket.GetConfigMore(&UserTicket.ArgsGetConfigMore{
		IDs:        args.UserTicket,
		HaveRemove: false,
	})
	if err != nil {
		err = nil
		//找不到则忽略处理，后续将跳过该部分设计
	}
	//////////////////////////////////////////////////////////////////////////////////////////////
	//2、获取商户设置
	//////////////////////////////////////////////////////////////////////////////////////////////
	//获取商户配送费强制性设置
	var transportOutAreaPrice int64
	transportOutAreaPrice, err = OrgCore.Config.GetConfigValInt64(&ClassConfig.ArgsGetConfig{
		BindID:    args.OrgID,
		Mark:      "TransportOutAreaPrice",
		VisitType: "admin",
	})
	if err != nil {
		err = nil
		transportOutAreaPrice = -1
	}
	//////////////////////////////////////////////////////////////////////////////////////////////
	//3、遍历和处理商品
	//////////////////////////////////////////////////////////////////////////////////////////////
	//遍历商品
	for _, vProduct := range args.Products {
		//计算商品最终价格
		errCode, err = getProductLastPrice(&vProduct, &resultData, args.OrgID, args.UserID, args.UserSubID, args.UserTicket, args.Address, args.TransportType, transportOutAreaPrice, args.UseIntegral, args.SkipProductCountLimit)
		if err != nil {
			err = errors.New(fmt.Sprint("get product last price, ", err))
			return
		}
	}
	//////////////////////////////////////////////////////////////////////////////////////////////
	//4、继续处理订单总的折扣
	//////////////////////////////////////////////////////////////////////////////////////////////
	//计算票据抵扣费用
	for kTicket, vTicket := range resultData.Tickets {
		//跳过非订单抵扣部分
		if !vTicket.UseOrder {
			continue
		}
		if vTicket.NeedCount > vTicket.UserCount {
			continue
		}
		//计算等待减免的总价格
		var waitPrice = resultData.LastPriceProduct + 0
		/**
		for _, v := range groupTicket {
			if v.TicketID != vTicket.TicketID {
				continue
			}
			for _, v2 := range resultData.Goods {
				if v.ProductID == v2.From.ID {
					waitPrice += v2.Price
				}
			}
		}
		if waitPrice < 1 {
			continue
		}
		*/
		//找出票据配置
		var vTicketConfig UserTicket.FieldsConfig
		for _, vConfig := range resultData.TicketConfigs {
			if vConfig.ID != vTicket.TicketID {
				continue
			}
			vTicketConfig = vConfig
		}
		//减免价格
		if vTicketConfig.ExemptionPrice > 0 {
			if vTicketConfig.ExemptionMinPrice > 0 {
				if waitPrice > vTicketConfig.ExemptionMinPrice {
					waitPrice = waitPrice - vTicketConfig.ExemptionPrice
				}
				if waitPrice < vTicketConfig.ExemptionMinPrice {
					waitPrice = vTicketConfig.ExemptionMinPrice
				}
			} else {
				waitPrice = waitPrice - vTicketConfig.ExemptionPrice
			}
		}
		if vTicketConfig.ExemptionDiscount > 0 {
			if vTicketConfig.ExemptionMinPrice > 0 {
				if waitPrice > vTicketConfig.ExemptionMinPrice {
					waitPrice = int64(float64(waitPrice) * (float64(vTicketConfig.ExemptionDiscount) / 100))
				}
				if vTicket.Price < vTicketConfig.ExemptionMinPrice {
					waitPrice = vTicketConfig.ExemptionMinPrice
				}
			} else {
				waitPrice = int64(float64(waitPrice) * (float64(vTicketConfig.ExemptionDiscount) / 100))
			}
		}
		var appendPrice = resultData.LastPriceProduct - waitPrice
		if appendPrice < 1 {
			appendPrice = 0
		}
		if appendPrice > 0 {
			resultData.Tickets[kTicket].Price = waitPrice
			resultData.Tickets[kTicket].NeedCount += 1
			resultData.Infos = append(resultData.Infos, DataProductPriceInfo{
				Mark:      "user_ticket",
				ProductID: 0,
				IsHide:    false,
				Price:     vTicket.Price,
				Msg:       "票据抵扣订单费用",
			})
			resultData.Exemptions = append(resultData.Exemptions, ServiceOrderWaitFields.FieldsExemption{
				System:   "user_ticket",
				ConfigID: vTicket.TicketID,
				Name:     vTicketConfig.Title,
				Des:      "票据折扣",
				Count:    1,
				Price:    waitPrice,
			})
			resultData.LastPriceTicket += appendPrice
			resultData.LastPriceProduct -= appendPrice
			resultData.LastPrice -= appendPrice
		}
	}
	//////////////////////////////////////////////////////////////////////////////////////////////
	//5、计算总的配送费
	//////////////////////////////////////////////////////////////////////////////////////////////
	resultData.TransportSystem = getTransportType(args.TransportType)
	switch args.TransportType {
	case 2:
		//根据配送方式识别配送费用
		//计算配送费用
		if transportOutAreaPrice > 0 {
			//检查是否超出分区
			var areaList []MapArea.FieldsArea
			areaList, err = MapArea.CheckPointInAreas(&MapArea.ArgsCheckPointInAreas{
				MapType: args.Address.MapType,
				Point: CoreSQLGPS.FieldsPoint{
					Longitude: args.Address.Longitude,
					Latitude:  args.Address.Latitude,
				},
				OrgID:    args.OrgID,
				IsParent: false,
				Mark:     "tms",
			})
			if err != nil || len(areaList) < 1 {
				//超出分区后，取消错误，并增加超区价格
				err = nil
				resultData.LastPriceTransport += transportOutAreaPrice
				resultData.LastPrice += transportOutAreaPrice
				resultData.Infos = append(resultData.Infos, DataProductPriceInfo{
					Mark:      "out_area",
					ProductID: 0,
					IsHide:    false,
					Price:     transportOutAreaPrice,
					Msg:       "超出配送区域，增加配送价格",
				})
			}
		} else {
			//计算配送模版费用
			for kTransport, vTransport := range resultData.Transports {
				var transportPrice int64
				transportPrice, err = getProductTransport(&vTransport, vTransport.Weight, vTransport.DistanceM, vTransport.BuyCount)
				if err != nil {
					errCode = "transport_price"
					err = errors.New(fmt.Sprint("get product transport, ", err))
					return
				}
				if transportPrice < 1 {
					continue
				}
				resultData.Transports[kTransport].Price = transportPrice
				resultData.LastPriceTransport += transportPrice
				resultData.LastPrice += transportPrice
			}
		}
	case 3:
		//跑腿费用
		var productAddress CoreSQLAddress.FieldsAddress
		if len(args.Products) > 0 {
			productAddress, _ = getProductAddress(args.Products[0].ID)
		}
		//如果组织配置、地图信息，获取发货地距离
		if productAddress.Latitude < 1 && productAddress.Longitude < 1 {
			orgConfigAddressGPS := OrgCore.Config.GetConfigValNoErr(args.OrgID, "OrderSendDefaultAddressGPS")
			if orgConfigAddressGPS != "" {
				orgConfigAddressGPSList := strings.Split(orgConfigAddressGPS, ",")
				if len(orgConfigAddressGPSList) > 1 {
					productAddress.Latitude, _ = strconv.ParseFloat(orgConfigAddressGPSList[0], 64)
					productAddress.Longitude, _ = strconv.ParseFloat(orgConfigAddressGPSList[1], 64)
				}
			}
		}
		if productAddress.Latitude < 1 && productAddress.Longitude < 1 {
			orgMapData, _ := OrgMap.GetMapByOrg(&OrgMap.ArgsGetMapByOrg{
				OrgID:   args.OrgID,
				IsAudit: true,
			})
			if orgMapData.ID > 0 {
				productAddress.Latitude = orgMapData.Latitude
				productAddress.Longitude = orgMapData.Longitude
			}
		}
		var goodsIDs []int64
		goodsCount := 0
		goodsWeight := 0
		for _, v := range resultData.ProductList {
			goodsWeight += v.Weight
			goodsCount += 1
			goodsIDs = append(goodsIDs, v.ID)
		}
		runningResult := TMSUserRunning.GetRunPrice(&TMSUserRunning.ArgsGetRunPrice{
			WaitAt:      "",
			GoodType:    "order",
			GoodWidget:  goodsWeight,
			FromAddress: productAddress,
			ToAddress:   args.Address,
		})
		//合计费用到配送费
		resultData.Transports = append(resultData.Transports, DataProductPriceTransport{
			ConfigData: FieldsTransport{},
			ProductIDs: goodsIDs,
			Weight:     goodsWeight,
			BuyCount:   goodsCount,
			DistanceM:  float64(runningResult.BetweenM),
			Price:      runningResult.TotalPrice,
		})
		resultData.LastPriceTransport += runningResult.TotalPrice
		resultData.LastPrice += runningResult.TotalPrice
	}
	//修正价格
	if resultData.LastPriceProduct < 1 {
		resultData.LastPriceProduct = 0
	}
	if resultData.LastPrice < 1 {
		resultData.LastPrice = 0
	}
	//记录购物车行为
	for _, v := range resultData.ProductList {
		MallLogMod.AppendLog(args.UserID, "", args.OrgID, v.ID, 2)
	}
	//反馈
	return
}

// getProductTransport 计算商品的配送费用
func getProductTransport(configData *DataProductPriceTransport, weight int, distanceM float64, buyCount int) (lastPrice int64, err error) {
	//处理免费条件
	if configData.ConfigData.FreeUnit > 0 {
		switch configData.ConfigData.FreeType {
		case 0:
			break
		case 1:
			if buyCount <= configData.ConfigData.FreeUnit {
				return
			}
		case 2:
			if weight <= configData.ConfigData.FreeUnit {
				return
			}
		case 3:
			if int(distanceM/1000) <= configData.ConfigData.FreeUnit {
				return
			}
		}
	}
	//处理首费
	switch configData.ConfigData.Rules {
	case 0:
		//免费配送
		break
	case 1:
		//按件计算
		if buyCount < configData.ConfigData.RulesUnit {
			break
		}
		lastPrice = configData.ConfigData.RulesPrice
	case 2:
		//按重量计算
		if weight < configData.ConfigData.RulesUnit {
			break
		}
		lastPrice = configData.ConfigData.RulesPrice
	case 3:
		//按公里数计算
		//如果符合免费标准，则跳出
		if distanceM/1000 < float64(configData.ConfigData.RulesUnit) {
			break
		}
		lastPrice = configData.ConfigData.RulesPrice
	}
	//计算叠加费用
	var addUnit int64
	switch configData.ConfigData.AddType {
	case 0:
	case 1:
		addUnit = int64(buyCount) / int64(configData.ConfigData.AddUnit)
	case 2:
		addUnit = int64(weight) / int64(configData.ConfigData.AddUnit)
	case 3:
		addUnit = int64((distanceM/float64(1000) - float64(configData.ConfigData.RulesUnit)) / float64(configData.ConfigData.AddUnit))
	}
	if addUnit > 1 {
		lastPrice = lastPrice + configData.ConfigData.AddPrice*addUnit
	}
	//反馈
	return
}

// 在商品中查询票据是否可用
func findTicketInProduct(productTickets []int64, findTicket int64, productSortID int64, orgID int64) (isFind bool) {
	for _, vTicketID := range productTickets {
		if vTicketID != findTicket {
			continue
		}
		isFind = true
		return
	}
	if productSortID < 1 {
		return
	}
	var sortParam string
	var err error
	sortParam, err = Sort.GetParam(&ClassSort.ArgsGetParam{
		ID:     productSortID,
		BindID: orgID,
		Mark:   "user_tickets",
	})
	if err != nil {
		return
	}
	//拆分数据结构体
	var sortParams []string
	sortParams = strings.Split(sortParam, ",")
	for _, v := range sortParams {
		var vInt64 int64
		vInt64, err = CoreFilter.GetInt64ByString(v)
		if err != nil {
			err = nil
			continue
		}
		if vInt64 == findTicket {
			isFind = true
			return
		}
	}
	return
}

// 获取商品的最终费用
func getProductLastPrice(vProduct *ArgsGetProductPriceProduct, resultData *DataProductPrice, orgID, userID int64, useUserSubID int64, useUserTicketIDs []int64, buyAddress CoreSQLAddress.FieldsAddress, transportType int, transportOutAreaPrice int64, useIntegral bool, skipProductCountLimit bool) (errCode string, err error) {
	//初始化
	exemptions := ServiceOrderWaitFields.FieldsExemptions{}
	//////////////////////////////////////////////////////////////////////////////////////////////
	//获取基本信息
	//////////////////////////////////////////////////////////////////////////////////////////////
	//获取商品
	productData, b := getProductByIDCanUse(vProduct.ID, orgID)
	if !b {
		resultData.Infos = append(resultData.Infos, DataProductPriceInfo{
			Mark:      "product_not_exist",
			ProductID: vProduct.ID,
			IsHide:    false,
			Price:     0,
			Msg:       fmt.Sprint("商品不存在，商品ID: ", vProduct.ID),
		})
		errCode = "err_mall_product_not_exist"
		err = errors.New("product not exist")
		return
	}
	if productData.ParentID > 0 {
		resultData.Infos = append(resultData.Infos, DataProductPriceInfo{
			Mark:      "product_out",
			ProductID: vProduct.ID,
			IsHide:    false,
			Price:     0,
			Msg:       fmt.Sprint(productData.Title, "，商品已下架"),
		})
		errCode = "err_mall_product_not_publish"
		err = errors.New("product not publish")
		return
	}
	if !skipProductCountLimit {
		var appendInfo DataProductPriceInfo
		appendInfo, errCode, b = buyCheckProductCount(&productData, vProduct.OptionKey, vProduct.BuyCount)
		if !b {
			resultData.Infos = append(resultData.Infos, appendInfo)
			err = errors.New("product not have count")
			return
		}
	}

	//////////////////////////////////////////////////////////////////////////////////////////////
	//获取商品基础折扣价
	//////////////////////////////////////////////////////////////////////////////////////////////
	//商品价格
	var productPrice = productData.Price + 0
	//附加商品的当前价格
	if productData.PriceExpireAt.Unix() >= CoreFilter.GetNowTime().Unix() {
		productPrice = productData.PriceReal
	} else {
		productPrice = productData.Price
	}

	//////////////////////////////////////////////////////////////////////////////////////////////
	//商品选项处理
	//////////////////////////////////////////////////////////////////////////////////////////////
	//如果商品存在选项
	if len(productData.OtherOptions.DataList) > 0 {
		//如果该商品没有进行选择，则跳出
		if vProduct.OptionKey == "" {
			//自动默认选择第一个商品
			/**
			errCode = "no_product_option"
			err = errors.New("product have options, buy no select option key")
			*/
			for _, vOption := range productData.OtherOptions.DataList {
				vProduct.OptionKey = vOption.Key
				if vProduct.OptionKey != "" {
					break
				}
			}
		}
		//遍历选项，找到对应的选项
		for _, v := range productData.OtherOptions.DataList {
			//跳过非本选项商品
			if v.Key != vProduct.OptionKey {
				continue
			}
			sort1Name := ""
			for k2, v2 := range productData.OtherOptions.Sort1.Options {
				if k2 == v.Sort1 {
					sort1Name = v2
				}
			}
			sort2Name := ""
			for k2, v2 := range productData.OtherOptions.Sort2.Options {
				if k2 == v.Sort2 {
					sort2Name = v2
				}
			}
			//覆盖商品的标题
			productTitle := fmt.Sprint(productData.Title, "#", sort1Name, "_", sort2Name)
			//如果存在商品，则覆盖本商品
			if v.ProductID > 0 {
				//重新获取商品
				productData, b = getProductByIDCanUse(vProduct.ID, orgID)
				if !b {
					resultData.Infos = append(resultData.Infos, DataProductPriceInfo{
						Mark:      "product_not_exist",
						ProductID: productData.ID,
						IsHide:    false,
						Price:     0,
						Msg:       fmt.Sprint("商品不存在，商品ID: ", productData.ID),
					})
					errCode = "err_mall_product_not_exist"
					err = errors.New("product not exist")
					return
				}
				if productData.ParentID > 0 {
					resultData.Infos = append(resultData.Infos, DataProductPriceInfo{
						Mark:      "product_out",
						ProductID: productData.ID,
						IsHide:    false,
						Price:     0,
						Msg:       fmt.Sprint(productTitle, "，商品已下架"),
					})
					errCode = "err_mall_product_not_publish"
					err = errors.New("product not publish")
					return
				}
			} else {
				//如果不存在商品，则仅覆盖价格
				if v.PriceExpireAt.Unix() >= CoreFilter.GetNowTime().Unix() {
					productPrice = v.PriceReal
				} else {
					productPrice = v.Price
				}
			}
			//保留旧的商品名称
			productData.Title = productTitle
		}
	}

	//////////////////////////////////////////////////////////////////////////////////////////////
	//计算会员价
	//////////////////////////////////////////////////////////////////////////////////////////////
	//计算会员抵扣费用
	exemptions, productPrice, err = getProductLastPriceSub(vProduct, &productData, resultData, userID, useUserSubID, productPrice, exemptions)
	if err != nil {
		errCode = "err_mall_product_user_sub"
		return
	}

	//////////////////////////////////////////////////////////////////////////////////////////////
	//计算积分价格
	//////////////////////////////////////////////////////////////////////////////////////////////
	transportFree := false
	//计算积分抵扣
	exemptions, productPrice, transportFree, err = getProductLastPriceIntegral(vProduct, &productData, resultData, useIntegral, productPrice, exemptions)
	if err != nil {
		errCode = "err_mall_product_user_integral"
		return
	}

	//////////////////////////////////////////////////////////////////////////////////////////////
	//计算票据价格
	//////////////////////////////////////////////////////////////////////////////////////////////
	exemptions, err = getProductLastPriceTicket(vProduct, &productData, resultData, userID, useUserTicketIDs, productPrice, exemptions)
	if err != nil {
		errCode = "err_mall_product_user_ticket"
		return
	}

	//////////////////////////////////////////////////////////////////////////////////////////////
	//计算配送费
	//////////////////////////////////////////////////////////////////////////////////////////////
	switch transportType {
	case 0:
		//0 self 自运营服务
		errCode, err = getProductLastPriceTransport(vProduct, &productData, resultData, transportFree, transportOutAreaPrice, buyAddress)
		if err != nil {
			return
		}
	}

	//////////////////////////////////////////////////////////////////////////////////////////////
	//计算最后价格
	//////////////////////////////////////////////////////////////////////////////////////////////
	//检查商品是否为虚拟商品，订单内必须有且只有一个才会生效
	var goodMark string
	if productData.IsVirtual {
		goodMark = "virtual"
	} else {
		virtual, b := productData.Params.GetValBool("virtual")
		if b && virtual {
			goodMark = "virtual"
		}
	}
	otherType, b := productData.Params.GetVal("otherType")
	if b && otherType != "" {
		goodMark = otherType
	}
	//记录商品信息
	resultData.ProductList = append(resultData.ProductList, productData)
	if productPrice < 0 {
		productPrice = 0
	}
	if productData.Price < 0 {
		productData.Price = 0
	}
	if productData.PriceReal < 0 {
		productData.PriceReal = 0
	}
	if productData.PriceExpireAt.Unix() > CoreFilter.GetNowTime().Unix() && productData.PriceReal != productData.Price {
		exemptions = append(exemptions, ServiceOrderWaitFields.FieldsExemption{
			System:   "mall",
			ConfigID: 0,
			Name:     "",
			Des:      "商品限时折扣",
			Count:    int64(vProduct.BuyCount),
			Price:    int64(vProduct.BuyCount) * (productData.Price - productData.PriceReal),
		})
	}
	resultData.Goods = append(resultData.Goods, ServiceOrderWaitFields.FieldsGood{
		From: CoreSQLFrom.FieldsFrom{
			System: "mall",
			ID:     productData.ID,
			Mark:   goodMark,
			Name:   productData.Title,
		},
		OptionKey:  vProduct.OptionKey,
		Count:      int64(vProduct.BuyCount),
		Price:      productPrice,
		Exemptions: exemptions,
	})
	//计算之前的价格
	resultData.BeforePrice += productData.Price * int64(vProduct.BuyCount)
	if resultData.BeforePrice < 0 {
		resultData.BeforePrice = 0
	}
	//记录最终价格
	resultData.LastPrice += productPrice * int64(vProduct.BuyCount)
	if resultData.LastPrice < 0 {
		resultData.LastPrice = 0
	}
	resultData.LastPriceProduct += productPrice * int64(vProduct.BuyCount)
	if resultData.LastPriceProduct < 0 {
		resultData.LastPriceProduct = 0
	}
	//反馈
	return
}

// 计算会员价格抵扣部分
func getProductLastPriceSub(vProduct *ArgsGetProductPriceProduct, productData *FieldsCore, resultData *DataProductPrice, userID int64, useUserSubID int64, productPrice int64, exemptions ServiceOrderWaitFields.FieldsExemptions) (newExemptions ServiceOrderWaitFields.FieldsExemptions, newPrice int64, err error) {
	//初始化
	newExemptions = exemptions
	newPrice = productPrice
	//如果可用会员为空，则退出
	if resultData.UserSubConfig.ID < 1 {
		return
	}
	// 是否找到的订阅
	findSub := false
	// 采用配置来减免费用
	needConfigEx := false
	var appendPrice int64 = 0
	for _, vSub := range productData.UserSubPrice {
		if vSub.ID != useUserSubID {
			continue
		}
		//检查用户持有会员情况
		if b := UserSubscription.CheckSub(&UserSubscription.ArgsCheckSub{
			ConfigID: vSub.ID,
			UserID:   userID,
		}); !b {
			continue
		}
		appendPrice = productPrice - vSub.Price
		//productPrice = vSub.Price
		findSub = true
		needConfigEx = false
	}
	if !findSub && productData.SortID > 0 {
		//如果没有找到，则需要在分类中查询
		// 如果找到，则在分类中查询
		var sortParam string
		sortParam, err = Sort.GetParam(&ClassSort.ArgsGetParam{
			ID:     productData.SortID,
			BindID: productData.OrgID,
			Mark:   "user_subs",
		})
		if err == nil {
			//拆分数据结构体
			var sortParams []string
			sortParams = strings.Split(sortParam, ",")
			for _, v := range sortParams {
				var vInt64 int64
				vInt64, err = CoreFilter.GetInt64ByString(v)
				if err != nil {
					err = nil
					continue
				}
				if vInt64 == useUserSubID {
					findSub = true
					needConfigEx = true
					break
				}
			}
		} else {
			err = nil
		}
	}
	//采用会员配置抵扣费用
	if needConfigEx {
		var newProductPrice = productPrice + 0
		if resultData.UserSubConfig.ExemptionPrice > 0 {
			if resultData.UserSubConfig.ExemptionMinPrice > 0 {
				if newProductPrice > resultData.UserSubConfig.ExemptionMinPrice {
					newProductPrice = newProductPrice - resultData.UserSubConfig.ExemptionPrice
				}
				if newProductPrice < resultData.UserSubConfig.ExemptionMinPrice {
					newProductPrice = resultData.UserSubConfig.ExemptionMinPrice
				}
			} else {
				newProductPrice = newProductPrice - resultData.UserSubConfig.ExemptionPrice
			}
		}
		if resultData.UserSubConfig.ExemptionDiscount > 0 {
			if resultData.UserSubConfig.ExemptionMinPrice > 0 {
				if newProductPrice > resultData.UserSubConfig.ExemptionMinPrice {
					newProductPrice = int64(float64(newProductPrice) * (float64(resultData.UserSubConfig.ExemptionDiscount) / 100))
				}
				if newProductPrice < resultData.UserSubConfig.ExemptionMinPrice {
					newProductPrice = resultData.UserSubConfig.ExemptionMinPrice
				}
			} else {
				newProductPrice = int64(float64(newProductPrice) * (float64(resultData.UserSubConfig.ExemptionDiscount) / 100))
			}
		}
		appendPrice = productPrice - newProductPrice
		if appendPrice < 1 {
			appendPrice = 0
		}
	}
	//如果还没有找到，则标记为无法使用会员商品
	if !findSub {
		resultData.NotSubProducts = append(resultData.NotSubProducts, productData.ID)
		resultData.Infos = append(resultData.Infos, DataProductPriceInfo{
			Mark:      "user_sub_not",
			ProductID: productData.ID,
			IsHide:    true,
			Price:     0,
			Msg:       "无法使用会员抵扣",
		})
	} else {
		if appendPrice > 0 {
			//记录信息
			resultData.Infos = append(resultData.Infos, DataProductPriceInfo{
				Mark:      "user_sub",
				ProductID: productData.ID,
				IsHide:    false,
				Price:     productPrice,
				Msg:       "会员抵扣",
			})
			exemptions = append(exemptions, ServiceOrderWaitFields.FieldsExemption{
				System:   "user_sub",
				ConfigID: resultData.UserSubConfig.ID,
				Name:     resultData.UserSubConfig.Title,
				Des:      "会员抵扣",
				Count:    int64(vProduct.BuyCount),
				Price:    int64(vProduct.BuyCount) * appendPrice,
			})
			resultData.LastPriceSub += int64(vProduct.BuyCount) * appendPrice
			productPrice -= appendPrice
		}
	}
	//反馈数据
	newExemptions = exemptions
	newPrice = productPrice
	return
}

// 计算积分抵扣部分
func getProductLastPriceIntegral(vProduct *ArgsGetProductPriceProduct, productData *FieldsCore, resultData *DataProductPrice, useIntegral bool, productPrice int64, exemptions ServiceOrderWaitFields.FieldsExemptions) (newExemptions ServiceOrderWaitFields.FieldsExemptions, newPrice int64, transportFree bool, err error) {
	//初始化
	newExemptions = exemptions
	newPrice = productPrice
	//如果价格已经低于1，则退出
	if productPrice < 1 {
		return
	}
	//是否可用积分？
	if !useIntegral {
		return
	}
	//商品是否可兑换积分
	if productData.Integral < 1 || productData.IntegralPrice < 1 {
		return
	}
	//用户持有积分少于商品积分，则不允许
	if resultData.LastUserIntegral < productData.Integral {

	}
	//开始处理
	if productData.IntegralTransportFree {
		//积分如果包含配送费，则直接跳过配送费计算
		// 记录最终价格
		transportFree = true
	}
	resultData.LastUserIntegral -= productData.Integral
	resultData.Infos = append(resultData.Infos, DataProductPriceInfo{
		Mark:      "integral",
		ProductID: productData.ID,
		IsHide:    false,
		Price:     productPrice,
		Msg:       "积分抵扣费用",
	})
	exemptions = append(exemptions, ServiceOrderWaitFields.FieldsExemption{
		System:   "user_integral",
		ConfigID: 0,
		Name:     "",
		Des:      "积分抵扣费用",
		Count:    int64(vProduct.BuyCount),
		Price:    int64(vProduct.BuyCount) * productData.IntegralPrice,
	})
	productPrice = productPrice - productData.IntegralPrice
	resultData.LastPriceIntegral += int64(vProduct.BuyCount) * productData.IntegralPrice
	//反馈数据
	newExemptions = exemptions
	newPrice = productPrice
	return
}

// 计算配送费用部分
func getProductLastPriceTransport(vProduct *ArgsGetProductPriceProduct, productData *FieldsCore, resultData *DataProductPrice, transportFree bool, transportOutAreaPrice int64, buyAddress CoreSQLAddress.FieldsAddress) (errCode string, err error) {
	//配送费免费则跳出
	if transportFree {
		return
	}
	//其他配送费用支持
	if transportOutAreaPrice > 0 {
		return
	}
	//如果不存在配送模板，则不计算配送费用
	if productData.TransportID < 1 {
		return
	}
	//找到商品配送模版
	var transportTempData FieldsTransport
	transportTempData, err = GetTransportID(&ArgsGetTransportID{
		ID:    productData.TransportID,
		OrgID: productData.OrgID,
	})
	if err != nil {
		//找不到配送模版
		errCode = "err_mall_core_transport_template_not_exist"
		err = errors.New(fmt.Sprint("not find transport template data, ", err))
		return
	}
	//模板规则为0则不计算退出
	if transportTempData.Rules == 0 {
		return
	}
	//计算配送费用
	// 计算配送距离
	var productAddress CoreSQLAddress.FieldsAddress
	productAddress, err = getProductAddress(productData.ID)
	if err != nil {
		err = errors.New(fmt.Sprint("get product address failed, ", err))
		errCode = "err_mall_core_product_address"
		return
	}
	var distanceM float64
	distanceM, err = MapMathPoint.GetDistanceM(&MapMathPoint.ArgsGetDistance{
		StartPoint: MapMathArgs.ParamsPoint{
			PointType: productAddress.GetMapType(),
			Longitude: productAddress.Longitude,
			Latitude:  productAddress.Latitude,
		},
		EndPoint: MapMathArgs.ParamsPoint{
			PointType: buyAddress.GetMapType(),
			Longitude: buyAddress.Longitude,
			Latitude:  buyAddress.Latitude,
		},
	})
	// 获取配送模版
	isFindTransportConfig := false
	transportConfigKey := 0
	for kConfig, vConfig := range resultData.Transports {
		if vConfig.ConfigData.ID == productData.TransportID && distanceM == vConfig.DistanceM {
			isFindTransportConfig = true
			transportConfigKey = kConfig
			break
		}
	}
	if !isFindTransportConfig {
		resultData.Transports = append(resultData.Transports, DataProductPriceTransport{
			ConfigData: transportTempData,
			ProductIDs: []int64{productData.ID},
			Weight:     0,
			BuyCount:   0,
			DistanceM:  distanceM,
			Price:      0,
		})
		transportConfigKey = len(resultData.Transports) - 1
	} else {
		resultData.Transports[transportConfigKey].ProductIDs = append(resultData.Transports[transportConfigKey].ProductIDs, productData.ID)
	}
	//叠加计算重量和购买数量
	resultData.Transports[transportConfigKey].Weight += productData.Weight
	resultData.Transports[transportConfigKey].BuyCount += vProduct.BuyCount
	return
}

// 计算票据费用部分
// 此方法不会修改商品的标的价格， 而是直接修改最终价格递减处理
func getProductLastPriceTicket(vProduct *ArgsGetProductPriceProduct, productData *FieldsCore, resultData *DataProductPrice, userID int64, useUserTicketIDs []int64, productPrice int64, exemptions ServiceOrderWaitFields.FieldsExemptions) (newExemptions ServiceOrderWaitFields.FieldsExemptions, err error) {
	//初始化
	newExemptions = exemptions
	//如果价格已经低于1，则退出
	if productPrice < 1 {
		return
	}
	//计算票据的价格
	// 等待抵扣用于订单的配置ID列
	for _, vUesUserTicketID := range useUserTicketIDs {
		//排除重复的配置ID
		isRepeat := false
		for _, vEx := range exemptions {
			if vEx.ConfigID == vUesUserTicketID {
				isRepeat = true
				break
			}
		}
		if isRepeat {
			continue
		}
		//在指定范围内查询票据
		isFindTicket := findTicketInProduct(productData.UserTicket, vUesUserTicketID, productData.SortID, productData.OrgID)
		if !isFindTicket {
			continue
		}
		//获取票据配置
		var vTicketConfig UserTicket.FieldsConfig
		for _, vConfig := range resultData.TicketConfigs {
			if vConfig.ID == vUesUserTicketID {
				vTicketConfig = vConfig
				break
			}
		}
		//找不到配置则跳过
		if vTicketConfig.ID < 1 {
			continue
		}
		//订单抵扣则跳过
		if vTicketConfig.UseOrder {
			continue
		}
		//将找到的数据，写入到票据预备结构内
		// 在已存在的列队中查询，如果存在则插入票据
		// 该步骤主要检查票据使用和用户持有量
		isFindInTickets := false
		noTicket := false
		ticketKey := -1
		useTicketCount := 0
		for kTicket, vTicket := range resultData.Tickets {
			if vTicket.TicketID != vUesUserTicketID {
				continue
			}
			//票据少于1张，标记无票据
			if resultData.Tickets[kTicket].UserCount < 1 {
				noTicket = true
				continue
			}
			isFindInTickets = true
			resultData.Tickets[kTicket].Products = append(resultData.Tickets[kTicket].Products, productData.ID)
			ticketKey = kTicket
			break
		}
		if noTicket {
			continue
		}
		//如果该票据已使用量达到用户持有量，则跳过
		if ticketKey != -1 {
			if resultData.Tickets[ticketKey].NeedCount >= resultData.Tickets[ticketKey].UserCount {
				continue
			}
		}
		if !isFindInTickets {
			//检查用户持有票据情况
			var userTicketCount int64
			userTicketCount, err = UserTicket.GetTicketCount(&UserTicket.ArgsGetTicketCount{
				ConfigID: vTicketConfig.ID,
				UserID:   userID,
			})
			if err != nil {
				err = nil
				userTicketCount = 0
			}
			//如果用户没有持有票据，则退出
			if userTicketCount < 1 {
				noTicket = true
			} else {
				//将票据信息写入总集合
				resultData.Tickets = append(resultData.Tickets, DataProductPriceTicket{
					TicketID:  vUesUserTicketID,
					Products:  []int64{productData.ID},
					SortIDs:   []int64{},
					UseOrder:  vTicketConfig.UseOrder,
					NeedCount: 0,
					UserCount: int(userTicketCount),
					Price:     0,
				})
				ticketKey = len(resultData.Tickets) - 1
			}
		}
		if noTicket || ticketKey == -1 {
			continue
		}
		//计算所需使用的票据
		if resultData.Tickets[ticketKey].UserCount < vProduct.BuyCount {
			useTicketCount = resultData.Tickets[ticketKey].UserCount
		} else {
			useTicketCount = vProduct.BuyCount
		}
		//开始抵扣费用
		var newProductPrice = productPrice + 0
		if vTicketConfig.ExemptionPrice > 0 {
			if vTicketConfig.ExemptionMinPrice > 0 {
				if newProductPrice > vTicketConfig.ExemptionMinPrice {
					newProductPrice = newProductPrice - vTicketConfig.ExemptionPrice
				}
				if newProductPrice < vTicketConfig.ExemptionMinPrice {
					newProductPrice = vTicketConfig.ExemptionMinPrice
				}
			} else {
				newProductPrice = newProductPrice - vTicketConfig.ExemptionPrice
			}
		}
		if vTicketConfig.ExemptionDiscount > 0 {
			if vTicketConfig.ExemptionMinPrice > 0 {
				if newProductPrice > vTicketConfig.ExemptionMinPrice {
					newProductPrice = int64(float64(newProductPrice) * (float64(vTicketConfig.ExemptionDiscount) / 100))
				}
				if newProductPrice < vTicketConfig.ExemptionMinPrice {
					newProductPrice = vTicketConfig.ExemptionMinPrice
				}
			} else {
				newProductPrice = int64(float64(newProductPrice) * (float64(vTicketConfig.ExemptionDiscount) / 100))
			}
		}
		if newProductPrice < 1 {
			newProductPrice = 0
		}
		var appendPrice = productPrice - newProductPrice
		if appendPrice < 1 {
			appendPrice = 0
		}
		//抵扣金额大于0才能生效，否则跳过
		if appendPrice > 0 {
			//标记折扣信息
			resultData.Infos = append(resultData.Infos, DataProductPriceInfo{
				Mark:      "user_ticket",
				ProductID: productData.ID,
				IsHide:    false,
				Price:     appendPrice,
				Msg:       "票据抵扣",
			})
			exemptions = append(exemptions, ServiceOrderWaitFields.FieldsExemption{
				System:   "user_ticket",
				ConfigID: vTicketConfig.ID,
				Name:     vTicketConfig.Title,
				Des:      "票据抵扣费用",
				Count:    int64(useTicketCount),
				Price:    int64(useTicketCount) * appendPrice,
			})
			//总结票据抵扣的总费用
			resultData.LastPriceTicket += int64(useTicketCount) * appendPrice
			//总结该票据使用次数
			resultData.Tickets[ticketKey].NeedCount += useTicketCount
			//productPrice -= appendPrice * int64(useTicketCount)
			//如果递减价格超出当前售价，则递减价格归结为售价
			if appendPrice > productPrice {
				appendPrice = productPrice
			}
			resultData.LastPrice -= int64(useTicketCount) * appendPrice
			resultData.LastPriceProduct -= int64(useTicketCount) * appendPrice
		}
	}
	//反馈数据
	newExemptions = exemptions
	return
}

// 检查商品的库存
func buyCheckProductCount(productData *FieldsCore, optionKey string, buyCount int) (appendInfo DataProductPriceInfo, errCode string, b bool) {
	//检查是否存在选项
	if len(productData.OtherOptions.DataList) > 0 {
		for _, vOption := range productData.OtherOptions.DataList {
			//找到对应的key，计算库存量
			if vOption.Key != optionKey {
				continue
			}
			if vOption.Count < buyCount {
				appendInfo = DataProductPriceInfo{
					Mark:      "product_count",
					ProductID: productData.ID,
					IsHide:    false,
					Price:     0,
					Msg:       fmt.Sprint(productData.Title, "，商品库存不足"),
				}
				errCode = "err_mall_product_count"
				return
			}
		}
	} else {
		//计算商品总的库存
		if productData.Count < buyCount {
			appendInfo = DataProductPriceInfo{
				Mark:      "product_count",
				ProductID: productData.ID,
				IsHide:    false,
				Price:     0,
				Msg:       fmt.Sprint(productData.Title, "，商品库存不足"),
			}
			errCode = "err_mall_product_count"
			return
		}
	}
	b = true
	return
}
