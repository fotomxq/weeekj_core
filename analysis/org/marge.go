package AnalysisOrg

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	FinanceAnalysis "github.com/fotomxq/weeekj_core/v5/finance/analysis"
	MallCore "github.com/fotomxq/weeekj_core/v5/mall/core"
	MarketCore "github.com/fotomxq/weeekj_core/v5/market/core"
	OrgUser "github.com/fotomxq/weeekj_core/v5/org/user"
	ServiceOrder "github.com/fotomxq/weeekj_core/v5/service/order"
	ServiceOrderAnalysis "github.com/fotomxq/weeekj_core/v5/service/order/analysis"
	TMSTransport "github.com/fotomxq/weeekj_core/v5/tms/transport"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

// ArgsMarge 聚合统计参数
type ArgsMarge struct {
	//多个标识码组
	Marks []ArgsMargeMark `json:"marks"`
	//商户ID
	OrgID int64 `json:"orgID"`
}

type ArgsMargeMark struct {
	//标识码
	Mark string `json:"mark"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//数量限制
	// 部分统计支持
	// 数据最多反馈1000条
	Limit int64 `json:"limit"`
}

// DataMarge 聚合统计反馈结构
type DataMarge struct {
	//数据结构
	Marks []DataMargeMark `json:"marks"`
}

type DataMargeMark struct {
	//标识码
	Mark string `json:"mark"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//数量限制
	// 部分统计支持
	// 数据最多反馈1000条
	Limit int64 `json:"limit"`
	//数据集合
	Data interface{} `json:"data"`
}

// GetMarge 获取聚合统计
func GetMarge(args *ArgsMarge) (result DataMarge, err error) {
	//遍历聚合函数处理
	for _, vMark := range args.Marks {
		if vMark.Limit < 1 || vMark.Limit > 1000 {
			continue
		}
		var timeBetween CoreSQLTime.FieldsCoreTime
		timeBetween, err = CoreSQLTime.GetBetweenByISO(vMark.TimeBetween)
		if err != nil {
			return
		}
		vResult := DataMargeMark{
			Mark:        vMark.Mark,
			TimeBetween: vMark.TimeBetween,
			Limit:       vMark.Limit,
		}
		switch vMark.Mark {
		case "finance_operating_income_month":
			//财务营业收入
			var dataList []FinanceAnalysis.DataGetAnalysis
			dataList, err = FinanceAnalysis.GetAnalysis(&FinanceAnalysis.ArgsGetAnalysis{
				TimeBetween:    timeBetween,
				TimeType:       "month",
				PaymentCreate:  CoreSQLFrom.FieldsFrom{},
				PaymentChannel: CoreSQLFrom.FieldsFrom{},
				PaymentFrom:    CoreSQLFrom.FieldsFrom{},
				TakeCreate: CoreSQLFrom.FieldsFrom{
					System: "org",
					ID:     args.OrgID,
					Mark:   "",
					Name:   "",
				},
				TakeChannel: CoreSQLFrom.FieldsFrom{},
				TakeFrom:    CoreSQLFrom.FieldsFrom{},
				Currency:    86,
				IsHistory:   false,
			})
			if err != nil {
				err = nil
				continue
			} else {
				vResult.Data = dataList
			}
		case "finance_savings_price":
			//储蓄账户转入资金总量
			vResult.Data, err = FinanceAnalysis.GetAnalysisTakePrice(&FinanceAnalysis.ArgsGetAnalysisTakePrice{
				TakeChannel: CoreSQLFrom.FieldsFrom{
					System: "deposit",
					ID:     0,
					Mark:   "savings",
					Name:   "",
				},
				TakeFrom: CoreSQLFrom.FieldsFrom{
					System: "org",
					ID:     args.OrgID,
					Mark:   "",
					Name:   "",
				},
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "finance_deposit_price":
			//押金账户转入资金总量
			vResult.Data, err = FinanceAnalysis.GetAnalysisTakePrice(&FinanceAnalysis.ArgsGetAnalysisTakePrice{
				TakeChannel: CoreSQLFrom.FieldsFrom{
					System: "deposit",
					ID:     0,
					Mark:   "deposit",
					Name:   "",
				},
				TakeFrom: CoreSQLFrom.FieldsFrom{
					System: "org",
					ID:     args.OrgID,
					Mark:   "",
					Name:   "",
				},
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "user_sub_count":
			//用户订阅量总量
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderCount(&ServiceOrder.ArgsGetAnalysisSystemOrderCount{
				OrgID:       args.OrgID,
				SystemMark:  "user_sub",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "user_sub_order_count":
			//订单销量统计
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderCount(&ServiceOrder.ArgsGetAnalysisSystemOrderCount{
				OrgID:       args.OrgID,
				SystemMark:  "user_sub",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "user_new":
			//用户新增总量
			vResult.Data, err = UserCore.GetAnalysisOrgCount(&UserCore.ArgsGetAnalysisOrgCount{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "user_active":
			//用户活跃总量
			vResult.Data, err = OrgUser.GetAnalysisActiveCount(&OrgUser.ArgsGetAnalysisActiveCount{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "tms_count":
			//配送单总量
			vResult.Data, err = TMSTransport.GetAnalysisCount(&TMSTransport.ArgsGetAnalysisCount{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "tms_wait_count":
			//配送单未完成总量
			vResult.Data, err = TMSTransport.GetAnalysisWaitCount(&TMSTransport.ArgsGetAnalysisWaitCount{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "tms_time_count_day":
			//按天拆分的配送量统计
			vResult.Data, err = TMSTransport.GetAnalysisTimeCount(&TMSTransport.ArgsGetAnalysisTimeCount{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
				TimeType:    "day",
			})
			if err != nil {
				err = nil
			}
		case "tms_time_count_month":
			//按月拆分的配送量统计
			vResult.Data, err = TMSTransport.GetAnalysisTimeCount(&TMSTransport.ArgsGetAnalysisTimeCount{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
				TimeType:    "month",
			})
			if err != nil {
				err = nil
			}
		case "tms_bind_count":
			//最近1月配送员数据
			vResult.Data, err = TMSTransport.GetAnalysisBind(&TMSTransport.ArgsGetAnalysisBind{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
			}
		case "tms_bind_avg":
			//最近1月配送员数据平均值
			vResult.Data, err = TMSTransport.GetAnalysisBindAvg(&TMSTransport.ArgsGetAnalysisBind{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
			}
		case "tms_bind_all_count":
			//计算配送安排人次
			vResult.Data, err = TMSTransport.GetAnalysisTakeBindCount(&TMSTransport.ArgsGetAnalysisTakeBindCount{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
			}
		case "tms_all_price":
			//配送单缴纳费用
			vResult.Data, err = TMSTransport.GetAnalysisPrice(&TMSTransport.ArgsGetAnalysisPrice{
				OrgID:       args.OrgID,
				BindID:      -1,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
			}
		case "tms_take_cash":
			//配送单收取的现金
			vResult.Data, err = TMSTransport.GetAnalysisCashSum(&TMSTransport.ArgsGetAnalysisCashSum{
				OrgID:       args.OrgID,
				BindID:      -1,
				PayType:     0,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
			}
		case "tms_payment_cash":
			//配送单收取的现金
			vResult.Data, err = TMSTransport.GetAnalysisCashSum(&TMSTransport.ArgsGetAnalysisCashSum{
				OrgID:       args.OrgID,
				BindID:      -1,
				PayType:     1,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
			}
		case "mall_buy_count":
			//商品购买量排名
			vResult.Data, err = MallCore.GetAnalysisCount(&MallCore.ArgsGetAnalysisCount{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
				Limit:       vMark.Limit,
			})
			if err != nil {
				err = nil
			}
		case "mall_order_count":
			//商城商品订单销量总数
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderCount(&ServiceOrder.ArgsGetAnalysisSystemOrderCount{
				OrgID:       args.OrgID,
				SystemMark:  "mall",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_count":
			//订单总数
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderCount(&ServiceOrder.ArgsGetAnalysisSystemOrderCount{
				OrgID:       args.OrgID,
				SystemMark:  "",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_price":
			//订单费用总数
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderPrice(&ServiceOrder.ArgsGetAnalysisSystemOrderCount{
				OrgID:       args.OrgID,
				SystemMark:  "",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_time_month_price":
			//订单费用总数，分时间段
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderPriceTime(&ServiceOrder.ArgsGetAnalysisSystemOrderPriceTime{
				OrgID:       args.OrgID,
				SystemMark:  "",
				BetweenTime: vMark.TimeBetween,
				TimeType:    "month",
			})
			if err != nil {
				err = nil
			}
		case "order_mall_price":
			//订单商城费用总数
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderPrice(&ServiceOrder.ArgsGetAnalysisSystemOrderCount{
				OrgID:       args.OrgID,
				SystemMark:  "",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_time_month_mall_price":
			//订单商城费用总数，分时间段
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderPriceTime(&ServiceOrder.ArgsGetAnalysisSystemOrderPriceTime{
				OrgID:       args.OrgID,
				SystemMark:  "mall",
				BetweenTime: vMark.TimeBetween,
				TimeType:    "month",
			})
			if err != nil {
				err = nil
			}
		case "order_refund_count":
			//订单退货总数
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderRefund(&ServiceOrder.ArgsGetAnalysisSystemOrderCount{
				OrgID:       args.OrgID,
				SystemMark:  "",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_refund_price":
			//订单退货总金额
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderRefundPrice(&ServiceOrder.ArgsGetAnalysisSystemOrderCount{
				OrgID:       args.OrgID,
				SystemMark:  "",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_time_month_refund_price":
			//订单退货总金额，分月统计
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderRefundPriceTime(&ServiceOrder.ArgsGetAnalysisSystemOrderPriceTime{
				OrgID:       args.OrgID,
				SystemMark:  "",
				BetweenTime: vMark.TimeBetween,
				TimeType:    "month",
			})
			if err != nil {
				err = nil
			}
		case "order_mall_refund_price":
			//订单商城退货总金额
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderRefundPrice(&ServiceOrder.ArgsGetAnalysisSystemOrderCount{
				OrgID:       args.OrgID,
				SystemMark:  "",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_time_month_mall_refund_price":
			//订单商城退货总金额，分月统计
			vResult.Data, err = ServiceOrder.GetAnalysisSystemOrderRefundPriceTime(&ServiceOrder.ArgsGetAnalysisSystemOrderPriceTime{
				OrgID:       args.OrgID,
				SystemMark:  "mall",
				BetweenTime: vMark.TimeBetween,
				TimeType:    "month",
			})
			if err != nil {
				err = nil
			}
		case "order_exemption_user_sub_count":
			//订单会员使用次数
			vResult.Data, err = ServiceOrderAnalysis.GetOrgExemptionCount(&ServiceOrderAnalysis.ArgsGetOrgExemptionCount{
				OrgID:       args.OrgID,
				FromSystem:  "user_sub",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_exemption_user_sub_price":
			//订单会员减免费用合计
			vResult.Data, err = ServiceOrderAnalysis.GetOrgExemptionPrice(&ServiceOrderAnalysis.ArgsGetOrgExemptionCount{
				OrgID:       args.OrgID,
				FromSystem:  "user_sub",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_exemption_user_ticket_count":
			//订单票据使用次数
			vResult.Data, err = ServiceOrderAnalysis.GetOrgExemptionCount(&ServiceOrderAnalysis.ArgsGetOrgExemptionCount{
				OrgID:       args.OrgID,
				FromSystem:  "user_ticket",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_exemption_user_ticket_price":
			//订单票据减免费用合计
			vResult.Data, err = ServiceOrderAnalysis.GetOrgExemptionPrice(&ServiceOrderAnalysis.ArgsGetOrgExemptionCount{
				OrgID:       args.OrgID,
				FromSystem:  "user_ticket",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_exemption_user_integral_count":
			//订单积分使用次数
			vResult.Data, err = ServiceOrderAnalysis.GetOrgExemptionCount(&ServiceOrderAnalysis.ArgsGetOrgExemptionCount{
				OrgID:       args.OrgID,
				FromSystem:  "user_integral",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_exemption_user_integral_price":
			//订单积分减免费用合计
			vResult.Data, err = ServiceOrderAnalysis.GetOrgExemptionPrice(&ServiceOrderAnalysis.ArgsGetOrgExemptionCount{
				OrgID:       args.OrgID,
				FromSystem:  "user_integral",
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "order_price_all_month":
			//订单分月费用统计
			vResult.Data, err = ServiceOrderAnalysis.GetAnalysisOrg(&ServiceOrderAnalysis.ArgsGetAnalysisOrg{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
				TimeType:    "month",
				CreateFrom:  -1,
				Limit:       vMark.Limit,
			})
			if err != nil {
				err = nil
			}
		case "market_new":
			//新增推广人数
			vResult.Data, err = MarketCore.GetAnalysisNewBind(&MarketCore.ArgsGetAnalysisNewBind{
				OrgID:       args.OrgID,
				BindID:      -1,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "market_new_have_order":
			//新推广转化人数，发生交易行为
			vResult.Data, err = MarketCore.GetAnalysisNewBindHavePrice(&MarketCore.ArgsGetAnalysisNewBind{
				OrgID:       args.OrgID,
				BindID:      -1,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "market_new_order_price":
			//新推广金额合计
			vResult.Data, err = MarketCore.GetAnalysisNewBindPrice(&MarketCore.ArgsGetAnalysisNewBind{
				OrgID:       args.OrgID,
				BindID:      -1,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "market_core_rank_10":
			//推广人排名，人数排名前10名
			vResult.Data, _, err = MarketCore.GetAnalysisCountBind(&MarketCore.ArgsGetAnalysisCountBind{
				Pages: CoreSQLPages.ArgsDataList{
					Page: 1,
					Max:  10,
					Sort: "count_count",
					Desc: true,
				},
				OrgID:       args.OrgID,
				ConfigID:    -1,
				TimeBetween: timeBetween,
			})
			if err != nil {
				err = nil
			}
		case "market_core_gps":
			//推广客户群体的分布情况
			vResult.Data, err = MarketCore.GetAnalysisGPS(&MarketCore.ArgsGetAnalysisGPS{
				OrgID:       args.OrgID,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
			}
		}
		result.Marks = append(result.Marks, vResult)
	}
	return
}
