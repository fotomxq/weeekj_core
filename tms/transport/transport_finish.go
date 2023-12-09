package TMSTransport

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	MapMathArgs "gitee.com/weeekj/weeekj_core/v5/map/math/args"
	MapMathPoint "gitee.com/weeekj/weeekj_core/v5/map/math/point"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceOrderMod "gitee.com/weeekj/weeekj_core/v5/service/order/mod"
)

// ArgsUpdateTransportFinish 更新配送单状态为完成参数
type ArgsUpdateTransportFinish struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//协助操作人员ID
	OperateBindID int64 `json:"operateBindID"`
	//是否为订单退款完成
	IsOrderRefund bool `json:"isOrderRefund"`
}

// UpdateTransportFinish 更新配送单状态为完成
func UpdateTransportFinish(args *ArgsUpdateTransportFinish) (err error) {
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	//检查是否为退单
	isRefund, b := data.Params.GetValBool("isRefund")
	if !b {
		isRefund = false
	}
	//必须已经付款，否则拒绝完成配送
	if data.Price > 0 && data.PayFinishAt.Unix() < 1000000 {
		err = errors.New("no pay")
		return
	}
	if args.IsOrderRefund {
		data.Params = CoreSQLConfig.Set(data.Params, "orderRefund", "true")
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), finish_at = NOW(), status = 3, params = :params WHERE id = :id AND (org_id = :org_id OR :org_id < 1) AND status = 2", map[string]interface{}{
			"id":     args.ID,
			"org_id": args.OrgID,
			"params": data.Params,
		})
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), finish_at = NOW(), status = 3 WHERE id = :id AND (org_id = :org_id OR :org_id < 1) AND status = 2", map[string]interface{}{
			"id":     args.ID,
			"org_id": args.OrgID,
		})
	}
	if err != nil {
		return
	}
	//获取支付方式
	payFromSystem, b := data.Params.GetVal("paySystem")
	if !b || payFromSystem == "" {
		payFromSystem = "unkonw"
	}
	if args.OperateBindID > 0 {
		_ = appendLog(&argsAppendLog{
			OrgID:           args.OrgID,
			BindID:          args.OperateBindID,
			TransportID:     args.ID,
			TransportBindID: data.BindID,
			Mark:            "finish",
			Des:             fmt.Sprint("平台工作人员，协助完成更新状态为完成配送"),
		})
	} else {
		_ = appendLog(&argsAppendLog{
			OrgID:           args.OrgID,
			BindID:          args.BindID,
			TransportID:     args.ID,
			TransportBindID: data.BindID,
			Mark:            "finish",
			Des:             fmt.Sprint("更新状态为完成配送"),
		})
	}
	if data.BindID > 0 {
		var km int64
		km, err = MapMathPoint.GetDistance(&MapMathPoint.ArgsGetDistance{
			StartPoint: MapMathArgs.ParamsPoint{
				PointType: data.FromAddress.GetMapType(),
				Longitude: data.FromAddress.Longitude,
				Latitude:  data.FromAddress.Latitude,
			},
			EndPoint: MapMathArgs.ParamsPoint{
				PointType: data.ToAddress.GetMapType(),
				Longitude: data.ToAddress.Longitude,
				Latitude:  data.ToAddress.Latitude,
			},
		})
		if err != nil {
			km = 0
		}
		err = appendAnalysis(&argsAppendAnalysis{
			OrgID:       data.OrgID,
			BindID:      data.BindID,
			InfoID:      data.InfoID,
			UserID:      data.UserID,
			TransportID: data.ID,
			KM:          km,
			OverTime:    CoreFilter.GetNowTimeCarbon().Time.Unix() - data.CreateAt.Unix(),
			Level:       0,
		})
		if err != nil {
			return
		}
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport_bind SET time_1_day = time_1_day + :time_1_day, count_finish_1_day = count_finish_1_day + 1 WHERE bind_id = :bind_id", map[string]interface{}{
			"bind_id":    data.BindID,
			"time_1_day": CoreFilter.GetNowTimeCarbon().Time.Unix() - data.CreateAt.Unix(),
		})
		if err != nil {
			return
		}
		//货物统计
		var dataBindAnalysisGoods FieldsBindAnalysisGoods
		if err = Router2SystemConfig.MainDB.Get(&dataBindAnalysisGoods, "SELECT id, create_at, org_id, bind_id, goods FROM tms_transport_bind_analysis_goods WHERE org_id = $1 AND bind_id = $2 AND create_at >= $3 AND create_at <= $4 AND is_refund = $5 LIMIT 1", data.OrgID, data.BindID, CoreFilter.GetNowTimeCarbon().StartOfDay().Time, CoreFilter.GetNowTimeCarbon().EndOfDay().Time, isRefund); err != nil || dataBindAnalysisGoods.ID < 1 {
			//创建新的数据
			err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tms_transport_bind_analysis_goods", "INSERT INTO tms_transport_bind_analysis_goods (org_id, bind_id, is_refund, goods) VALUES (:org_id,:bind_id,:is_refund,:goods)", map[string]interface{}{
				"org_id":    data.OrgID,
				"bind_id":   data.BindID,
				"is_refund": isRefund,
				"goods":     FieldsBindAnalysisGoodsGoods{},
			}, &dataBindAnalysisGoods)
			if err != nil {
				err = errors.New(fmt.Sprint("create analysis bind goods, ", err))
				return
			}
		}
		for _, v := range data.Goods {
			isFind := false
			var orderData ServiceOrderMod.FieldsOrder
			orderData, err = ServiceOrderMod.GetByID(&ServiceOrderMod.ArgsGetByID{
				ID:     data.OrderID,
				OrgID:  data.OrgID,
				UserID: -1,
			})
			if err != nil {
				err = nil
				continue
			}
			var goodPrice int64 = 0
			for _, v2 := range orderData.Goods {
				if v.System == v2.From.System && v.ID == v2.From.ID {
					goodPrice = v2.Count * v2.Price
					for _, v3 := range v2.Exemptions {
						goodPrice -= v3.Price
					}
					if goodPrice < 1 {
						goodPrice = 0
					}
				}
			}
			for k2, v2 := range dataBindAnalysisGoods.Goods {
				if v.System == v2.System && v.ID == v2.ID && data.OrderID > 0 && v2.PaySystem == payFromSystem {
					if orderData.CreateFrom == v2.FromSystem {
						isFind = true
						dataBindAnalysisGoods.Goods[k2].Count += int64(v.Count)
						dataBindAnalysisGoods.Goods[k2].Price += goodPrice
						break
					}
				}
			}
			if !isFind {
				dataBindAnalysisGoods.Goods = append(dataBindAnalysisGoods.Goods, FieldsBindAnalysisGoodsGood{
					System:     v.System,
					ID:         v.ID,
					FromSystem: orderData.CreateFrom,
					PaySystem:  payFromSystem,
					Count:      int64(v.Count),
					Price:      goodPrice,
				})
			}
		}
		if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport_bind_analysis_goods SET goods = :goods WHERE id = :id", map[string]interface{}{
			"goods": dataBindAnalysisGoods.Goods,
			"id":    dataBindAnalysisGoods.ID,
		}); err != nil {
			err = errors.New(fmt.Sprint("update analysis bind goods, ", err))
			return
		}
	}
	if err != nil {
		return
	}
	//推送MQTT更新
	pushMQTTTransportUpdate(data.OrgID, data.BindID, data.ID, 1)
	//推送nats更新
	pushNatsStatusUpdate("finish", data.ID, "配送单配送完成")
	//统计
	if data.BindID > 0 {
		pushNatsAnalysisBind(data.BindID)
	}
	//反馈
	return
}
