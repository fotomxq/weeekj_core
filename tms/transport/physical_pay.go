package TMSTransport

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinancePhysicalPay "github.com/fotomxq/weeekj_core/v5/finance/physical_pay"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsPayPhysical 实物支付配送单参数
type ArgsPayPhysical struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//抵扣物品数量集合
	Data []ArgsPayPhysicalData `json:"data"`
}

// ArgsPayPhysicalData 抵扣物品
type ArgsPayPhysicalData struct {
	//获取来源
	// 如果商品mark带有virtual标记，且订单商品全部带有该标记，订单将在付款后直接完成
	From CoreSQLFrom.FieldsFrom `db:"from" json:"from"`
	//给予标的物数量
	PhysicalCount int64 `db:"physical_count" json:"physicalCount" check:"int64Than0"`
}

// PayPhysical 实物支付配送单
func PayPhysical(args *ArgsPayPhysical) (errCode string, err error) {
	//获取配送单
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		errCode = "transport_not_exist"
		return
	}
	//采用实物抵扣
	var logData []FinancePhysicalPay.ArgsCreateLogData
	for _, v := range data.Goods {
		err = FinancePhysicalPay.CheckPhysicalByFrom(&FinancePhysicalPay.ArgsGetPhysicalByFrom{
			OrgID: args.OrgID,
			BindFrom: CoreSQLFrom.FieldsFrom{
				System: v.System,
				ID:     v.ID,
				Mark:   v.Mark,
				Name:   v.Name,
			},
		})
		if err != nil {
			errCode = "check_physical"
			err = errors.New(fmt.Sprint("have no support physical data, org id: ", args.OrgID, ", mall product: {", v.System, ", ", v.ID, ", ", v.Mark, ", ", v.Name, "}", ", err: ", err))
			return
		}
		isFind := false
		var physicalCount int64 = 0
		for _, v2 := range args.Data {
			if v2.From.System == v.System && v2.From.Mark == v.Mark && v2.From.ID == v.ID {
				physicalCount = v2.PhysicalCount
				isFind = true
				break
			}
		}
		if !isFind {
			errCode = "no_physical"
			err = errors.New("need more physical count")
			return
		}
		logData = append(logData, FinancePhysicalPay.ArgsCreateLogData{
			PhysicalCount: physicalCount,
			BindFrom: CoreSQLFrom.FieldsFrom{
				System: v.System,
				ID:     v.ID,
				Mark:   v.Mark,
				Name:   v.Name,
			},
			BindCount: int64(v.Count),
		})
	}
	var newIDs pq.Int64Array
	params := CoreSQLConfig.FieldsConfigsType{
		{
			Mark: "order_id",
			Val:  fmt.Sprint(data.OrderID),
		},
		{
			Mark: "tms_id",
			Val:  fmt.Sprint(data.ID),
		},
	}
	newIDs, err = FinancePhysicalPay.CreateLog(&FinancePhysicalPay.ArgsCreateLog{
		OrgID:  data.OrgID,
		BindID: data.BindID,
		UserID: data.UserID,
		System: "tms",
		Data:   logData,
		Params: params,
	})
	if err != nil {
		errCode = "create_log"
		return
	}
	//记录扩展参数数据
	var newIDsStr string
	for _, v := range newIDs {
		if newIDsStr == "" {
			newIDsStr = fmt.Sprint(v)
		} else {
			newIDsStr = fmt.Sprint(newIDsStr, ",", v)
		}
	}
	data.Params = CoreSQLConfig.Set(data.Params, "physical_log_ids", newIDsStr)
	data.Params = CoreSQLConfig.Set(data.Params, "paySystem", "tms_physical_pay")
	//标记完成配送单支付
	if data.PayFinishAt.Unix() > 100000 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), params = :params WHERE id = :id AND (org_id = :org_id OR :org_id < 1)", map[string]interface{}{
			"id":     args.ID,
			"org_id": args.OrgID,
			"params": data.Params,
		})
		if err != nil {
			errCode = "update_transport"
			return
		}
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport SET update_at = NOW(), pay_finish_at = NOW(), params = :params WHERE id = :id AND (org_id = :org_id OR :org_id < 1) AND pay_finish_at < to_timestamp(1000000)", map[string]interface{}{
			"id":     args.ID,
			"org_id": args.OrgID,
			"params": data.Params,
		})
		if err != nil {
			errCode = "update_transport"
			return
		}
	}
	_ = appendLog(&argsAppendLog{
		OrgID:           args.OrgID,
		BindID:          args.BindID,
		TransportID:     args.ID,
		TransportBindID: args.BindID,
		Mark:            "pay_cash",
		Des:             fmt.Sprint("实物形式缴纳配送费"),
	})
	return
}
