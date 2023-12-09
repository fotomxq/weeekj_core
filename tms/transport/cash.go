package TMSTransport

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	FinancePay "gitee.com/weeekj/weeekj_core/v5/finance/pay"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetCashList 获取现金收取列表参数
type ArgsGetCashList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//收付款类型
	// 0 收款 1 付款
	PayType int `db:"pay_type" json:"payType"`
}

// GetCashList 获取现金收取列表
func GetCashList(args *ArgsGetCashList) (dataList []FieldsCash, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.PayType > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "pay_type = :pay_type"
		maps["pay_type"] = args.PayType
	}
	if where == "" {
		where = "true"
	}
	tableName := "tms_transport_cash"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, bind_id, transport_id, pay_type, price FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetCashByTransport 获取配送单对应的收支记录参数
type ArgsGetCashByTransport struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id"`
}

// GetCashByTransport 获取配送单对应的收支记录
func GetCashByTransport(args *ArgsGetCashByTransport) (data FieldsCash, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, bind_id, transport_id, pay_type, price FROM tms_transport_cash WHERE org_id = $1 AND transport_id = $2", args.OrgID, args.TransportID)
	if err == nil && data.ID < 1 {
		err = errors.New("data not exist")
		return
	}
	return
}

// ArgsUpdateTransportPayClient 代客户确认配送单付款参数
type ArgsUpdateTransportPayClient struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//操作IP
	IP string
}

// UpdateTransportPayClient 代客户确认配送单付款
func UpdateTransportPayClient(args *ArgsUpdateTransportPayClient) (payData FinancePay.FieldsPayType, result interface{}, needResult bool, errCode string, err error) {
	//获取配送单
	var data FieldsTransport
	data, err = GetTransport(&ArgsGetTransport{
		ID:     args.ID,
		OrgID:  args.OrgID,
		InfoID: -1,
		UserID: -1,
	})
	if err != nil {
		errCode = "transport_not_exist"
		return
	}
	if data.DeleteAt.Unix() > 1000000 {
		errCode = "transport_not_exist"
		err = errors.New("transport is delete")
		return
	}
	if data.PayFinishAt.Unix() > 1000000 {
		errCode = "transport_is_pay"
		err = errors.New("transport is pay")
		return
	}
	if data.PayID < 1 {
		errCode = "pay_not_exist"
		err = errors.New("pay not exist")
		return
	}
	//检查支付请求
	payData, err = FinancePay.GetOne(&FinancePay.ArgsGetOne{
		ID:  data.PayID,
		Key: "",
	})
	if err != nil {
		errCode = "pay_not_exist"
		return
	}
	//检查支付状态
	if payData.Status != 0 {
		errCode = "pay_status"
		err = errors.New("pay status not 0")
		return
	}
	//客户端确认支付
	payData, result, needResult, errCode, err = FinancePay.UpdateStatusClient(&FinancePay.ArgsUpdateStatusClient{
		CreateInfo: payData.CreateInfo,
		ID:         payData.ID,
		Key:        "",
		Params:     nil,
		IP:         args.IP,
	})
	return
}

// ArgsUpdateTransportCash 完成配送单的现金收取参数
type ArgsUpdateTransportCash struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作组织人员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// UpdateTransportCash 完成配送单的现金收取
func UpdateTransportCash(args *ArgsUpdateTransportCash) (errCode string, err error) {
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
	if data.DeleteAt.Unix() > 1000000 {
		errCode = "transport_not_exist"
		err = errors.New("transport is delete")
		return
	}
	if data.PayFinishAt.Unix() > 1000000 {
		errCode = "transport_is_pay"
		err = errors.New("transport is pay")
		return
	}
	if data.PayID < 1 {
		errCode = "pay_not_exist"
		err = errors.New("pay not exist")
		return
	}
	//检查支付请求
	var payData FinancePay.FieldsPayType
	payData, err = FinancePay.GetOne(&FinancePay.ArgsGetOne{
		ID:  data.PayID,
		Key: "",
	})
	if err != nil {
		errCode = "pay_not_exist"
		return
	}
	//检查支付状态
	if payData.Status != 0 && payData.Status != 1 {
		errCode = "pay_status"
		err = errors.New("pay status not 0 or 1")
		return
	}
	if payData.PaymentChannel.System != "cash" && payData.TakeChannel.System != "cash" {
		errCode = "pay_not_cash"
		err = errors.New("pay not cash")
		return
	}
	//标记支付完成
	errCode, err = FinancePay.UpdateStatusFinish(&FinancePay.ArgsUpdateStatusFinish{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "org_bind",
			ID:     args.BindID,
			Mark:   "",
			Name:   "",
		},
		ID:     payData.ID,
		Key:    "",
		Params: nil,
	})
	if err != nil {
		return
	}
	//标记日志
	isRefund, _ := data.Params.GetValBool("isRefund")
	var isRefundInt int
	if isRefund {
		isRefundInt = 1
		paySystem, b := data.Params.GetVal("paySystem")
		if !b || paySystem == "" {
			data.Params = CoreSQLConfig.Set(data.Params, "paySystem", "tms_refund_cash")
		}
	} else {
		isRefundInt = 0
		paySystem, b := data.Params.GetVal("paySystem")
		if !b || paySystem == "" {
			data.Params = CoreSQLConfig.Set(data.Params, "paySystem", "tms_cash")
		}
	}
	err = createCash(data.ID, isRefundInt, data.OrgID, args.BindID, data.Price)
	if err != nil {
		errCode = "insert"
		return
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
		_ = appendLog(&argsAppendLog{
			OrgID:           args.OrgID,
			BindID:          args.BindID,
			TransportID:     args.ID,
			TransportBindID: args.BindID,
			Mark:            "pay_cash",
			Des:             fmt.Sprint("现金形式缴纳配送费"),
		})
	}
	return
}

// 创建新的记录
func createCash(transportID int64, payType int, orgID, bindID, price int64) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tms_transport_cash (org_id, bind_id, transport_id, pay_type, price) VALUES (:org_id,:bind_id,:transport_id,:pay_type,:price)", map[string]interface{}{
		"org_id":       orgID,
		"bind_id":      bindID,
		"transport_id": transportID,
		"pay_type":     payType,
		"price":        price,
	})
	return
}
