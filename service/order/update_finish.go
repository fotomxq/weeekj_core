package ServiceOrder

import (
	"errors"
	"fmt"
	BaseEarlyWarning "gitee.com/weeekj/weeekj_core/v5/base/early_warning"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	MallCore "gitee.com/weeekj/weeekj_core/v5/mall/core"
	MallCoreMod "gitee.com/weeekj/weeekj_core/v5/mall/core/mod"
	MarketGivingBuyMall "gitee.com/weeekj/weeekj_core/v5/market/giving_buy_mall"
	MarketGivingNewUserMod "gitee.com/weeekj/weeekj_core/v5/market/giving_new_user/mod"
	MarketGivingUserSub "gitee.com/weeekj/weeekj_core/v5/market/giving_user_sub"
	OrgUserMod "gitee.com/weeekj/weeekj_core/v5/org/user/mod"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceOrderAnalysis "gitee.com/weeekj/weeekj_core/v5/service/order/analysis"
	TMSTransport "gitee.com/weeekj/weeekj_core/v5/tms/transport"
	UserIntegral "gitee.com/weeekj/weeekj_core/v5/user/integral"
)

// ArgsUpdateFinish 完成订单参数
type ArgsUpdateFinish struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//调整说明描述
	Des string `db:"des" json:"des"`
}

// UpdateFinish 完成订单
func UpdateFinish(args *ArgsUpdateFinish) (err error) {
	// 获取订单信息
	var orderData FieldsOrder
	orderData, err = GetByID(&ArgsGetByID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		UserID: args.UserID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("get order data, ", err))
		return
	}
	//检查订单是否已经完成？
	if orderData.Status == 4 {
		err = errors.New("replace")
		return
	}
	//更新订单
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "finish", args.Des)
	if err != nil {
		err = errors.New(fmt.Sprint("get log, ", err))
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), status = 4, logs = logs || :log WHERE id = :id AND status != 4 AND price_pay = true", map[string]interface{}{
		"id":  args.ID,
		"log": newLog,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update order finish, ", err))
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//处理新用户奖励部分
	updateFinishNewUserGiving(&orderData)
	//处理赠礼环节
	updateFinishGiving(&orderData)
	//更新支付渠道扩展标记
	// 找到订单支付的渠道
	var payFromSystem string
	if orderData.Price > 0 {
		if orderData.TransportID > 0 && orderData.PayID < 1 {
			var tmsData TMSTransport.FieldsTransport
			tmsData, err = TMSTransport.GetTransport(&TMSTransport.ArgsGetTransport{
				ID:     orderData.TransportID,
				OrgID:  -1,
				InfoID: -1,
				UserID: -1,
			})
			var b bool
			payFromSystem, b = tmsData.Params.GetVal("paySystem")
			if !b {
				payFromSystem = ""
			}
		}
	}
	//补充记录支付渠道
	if payFromSystem != "" {
		orderData.Params = CoreSQLConfig.Set(orderData.Params, "paySystem", payFromSystem)
	}
	//修改订单的扩展参数
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET params = :params WHERE id = :id", map[string]interface{}{
		"id":     orderData.ID,
		"params": orderData.Params,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update order params, ", err))
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//统计数据
	var goods []ServiceOrderAnalysis.ArgsCreateAnalysisGood
	var exemptions []ServiceOrderAnalysis.ArgsCreateAnalysisExemption
	for _, v := range orderData.Goods {
		goods = append(goods, ServiceOrderAnalysis.ArgsCreateAnalysisGood{
			From:  v.From,
			Count: v.Count,
		})
		for _, v2 := range v.Exemptions {
			exemptions = append(exemptions, ServiceOrderAnalysis.ArgsCreateAnalysisExemption{
				System:   v2.System,
				ConfigID: v2.ConfigID,
				Count:    v2.Count,
				Price:    v2.Price,
			})
		}
	}
	for _, v := range orderData.Exemptions {
		exemptions = append(exemptions, ServiceOrderAnalysis.ArgsCreateAnalysisExemption{
			System:   v.System,
			ConfigID: v.ConfigID,
			Count:    v.Count,
			Price:    v.Price,
		})
	}
	err = ServiceOrderAnalysis.CreateAnalysis(&ServiceOrderAnalysis.ArgsCreateAnalysis{
		OrgID:      args.OrgID,
		UserID:     args.UserID,
		SystemMark: orderData.SystemMark,
		CreateFrom: orderData.CreateFrom,
		Currency:   orderData.Currency,
		Price:      orderData.Price,
		Goods:      goods,
		Exemptions: exemptions,
	})
	if err != nil {
		CoreLog.Warn("create order and create analysis, ", err)
		err = nil
	}
	//推送nats
	CoreNats.PushDataNoErr("/service/order/update", "finish", orderData.ID, "", nil)
	//更新组织用户数据
	if orderData.OrgID > 0 && orderData.UserID > 0 {
		OrgUserMod.PushUpdateUserData(orderData.OrgID, orderData.UserID)
	}
	//反馈
	return
}

// 处理完成订单新用户注册奖励
func updateFinishNewUserGiving(orderData *FieldsOrder) {
	MarketGivingNewUserMod.PushNewUserBuy(orderData.UserID, orderData.Price, true)
}

// 完成订单的赠礼处理环节
func updateFinishGiving(orderData *FieldsOrder) {
	var err error
	//获取推荐人
	referrerUserID, _ := orderData.Params.GetValInt64("referrerUserID")
	referrerBindID, _ := orderData.Params.GetValInt64("referrerBindID")
	//处理赠礼设计
	var productIDs, sortIDs, tags []int64
	haveMall := false
	for _, v := range orderData.Goods {
		switch v.From.System {
		case "mall":
			//购买商品行为
			haveMall = true
			var vProductData MallCoreMod.FieldsCore
			vProductData, err = MallCoreMod.GetProduct(&MallCoreMod.ArgsGetProduct{
				ID:    v.From.ID,
				OrgID: orderData.OrgID,
			})
			if err != nil {
				err = errors.New(fmt.Sprint("get order product, id: ", v.From.ID, ", err: ", err))
				return
			}
			productIDs = append(productIDs, vProductData.ID)
			if vProductData.SortID > 0 {
				sortIDs = append(sortIDs, vProductData.SortID)
			}
			for _, v2 := range vProductData.Tags {
				tags = append(tags, v2)
			}
			//增加销量
			err = MallCore.UpdateProductBuy(v.From.ID, int(v.Count))
			if err != nil {
				CoreLog.Error("update product buy, mall id: ", v.From.ID, ", err: ", err)
				err = nil
			}
		case "user_sub":
			//购买用户会员行为
			forceMarketGivingSubID, b := orderData.Params.GetValInt64("force_market_giving_sub_id")
			if !b {
				forceMarketGivingSubID = 0
			}
			_, err = MarketGivingUserSub.CreateLog(&MarketGivingUserSub.ArgsCreateLog{
				OrgID:               orderData.OrgID,
				UserID:              orderData.UserID,
				ReferrerUserID:      referrerUserID,
				ReferrerBindID:      referrerBindID,
				PriceTotal:          orderData.Price,
				SubConfigID:         v.From.ID,
				SubBuyCount:         v.Count,
				LockGivingUserSubID: forceMarketGivingSubID,
			})
			if err != nil {
				CoreLog.Warn("create order, create market giving user sub, ", err)
				err = nil
			}
		case "user_integral":
			//购买积分行为
			// 不能是商户行为
			if orderData.OrgID > 0 {
				_ = BaseEarlyWarning.SendMod(&BaseEarlyWarning.ArgsSendMod{
					Mark: "MallOrderWarning",
					Contents: map[string]string{
						"OrgID":   fmt.Sprint(orderData.OrgID),
						"UserID":  fmt.Sprint(orderData.UserID),
						"OrderID": fmt.Sprint(orderData.ID),
					},
				})
				continue
			}
			// 根据商品数量给与该用户或商户积分
			userIntegralOrgOrUser, _ := orderData.Params.GetValBool("userIntegralOrgOrUser")
			if userIntegralOrgOrUser {
				err = UserIntegral.AddCount(&UserIntegral.ArgsAddCount{
					OrgID:    v.From.ID,
					UserID:   0,
					AddCount: v.Count,
					Des:      "购买平台积分",
				})
				if err != nil {
					CoreLog.Warn("create order, add org user integral failed, ", err)
					err = nil
				}
			} else {
				err = UserIntegral.AddCount(&UserIntegral.ArgsAddCount{
					OrgID:    0,
					UserID:   orderData.UserID,
					AddCount: v.Count,
					Des:      "购买平台积分",
				})
				if err != nil {
					CoreLog.Warn("create order, add org user integral failed, ", err)
					err = nil
				}
			}
		}
	}
	if haveMall {
		_, err = MarketGivingBuyMall.CreateLog(&MarketGivingBuyMall.ArgsCreateLog{
			OrgID:          orderData.OrgID,
			UserID:         orderData.UserID,
			ReferrerUserID: referrerUserID,
			ReferrerBindID: referrerBindID,
			PriceTotal:     orderData.Price,
			MallProductID:  productIDs,
			SortID:         sortIDs,
			Tag:            tags,
		})
		if err != nil {
			CoreLog.Warn("create order, create market giving buy mall, ", err)
			err = nil
		}
	}
}
