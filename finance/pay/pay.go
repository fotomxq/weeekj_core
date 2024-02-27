package FinancePay

import (
	"encoding/json"
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime2 "github.com/fotomxq/weeekj_core/v5/core/sql/time2"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	FinanceLog "github.com/fotomxq/weeekj_core/v5/finance/log"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetList 获取请求列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//状态 必须填写
	// 0 wait 客户端发起付款，并正在支付中
	// 1 client 客户端完成支付，等待服务端验证
	// 2 failed 交易失败，服务端主动取消交易或其他原因取消交易
	// 3 finish 交易成功
	// 4 remove 交易销毁
	// 5 expire 交易过期
	// 6 refund 发起退款申请
	// 7 refundAudit 退款审核通过，等待处理中
	// 8 refundFailed 退款失败
	// 9 refundFinish 退款完成
	Status []int
	//付款人来源
	PaymentCreate CoreSQLFrom.FieldsFrom `db:"payment_create" json:"paymentCreate"`
	//支付方式
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	PaymentChannel CoreSQLFrom.FieldsFrom `db:"payment_channel" json:"paymentChannel"`
	//支付方的来源
	// 留空则代表平台方，否则为商户或加盟商
	PaymentFrom CoreSQLFrom.FieldsFrom `db:"payment_from" json:"paymentFrom"`
	//收款人来源
	TakeCreate CoreSQLFrom.FieldsFrom `db:"take_create" json:"takeCreate"`
	//收款渠道
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	TakeChannel CoreSQLFrom.FieldsFrom `db:"take_channel" json:"takeChannel"`
	//收款方来源
	// 留空则代表平台方，否则为商户或加盟商
	TakeFrom CoreSQLFrom.FieldsFrom `db:"take_from" json:"takeFrom"`
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//最小金额
	// -1则跳过
	MinPrice int64
	//最大金额
	// -1则跳过
	MaxPrice int64
	//查询时间范围
	TimeBetween CoreSQLTime2.FieldsCoreTime `db:"time_between" json:"timeBetween"`
	//是否需要退款相关参数
	//是否发起了退款
	//退款发起的金额
	//扩展数据查询
	Params CoreSQLConfig.FieldsConfigType
	//支付失败后的代码
	// 用于系统识别错误类型
	FailedCode string `db:"failed_code" json:"failedCode"`
	//是否为历史
	IsHistory bool
	//搜索
	Search string
}

// GetList 获取请求列表
func GetList(args *ArgsGetList) (dataList []FieldsPayType, dataCount int64, err error) {
	if len(args.Status) < 1 {
		args.Status = []int{0, 1, 2, 6, 7, 8}
	}
	where := "status = ANY(:status)"
	maps := map[string]interface{}{
		"status": pq.Array(args.Status),
	}
	where, maps, err = args.PaymentCreate.GetListAnd("payment_create", "payment_create", where, maps)
	if err != nil {
		return
	}
	where, maps, err = args.PaymentChannel.GetListAnd("payment_channel", "payment_channel", where, maps)
	if err != nil {
		return
	}
	where, maps, err = args.PaymentFrom.GetListAnd("payment_from", "payment_from", where, maps)
	if err != nil {
		return
	}
	where, maps, err = args.TakeCreate.GetListAnd("take_create", "take_create", where, maps)
	if err != nil {
		return
	}
	where, maps, err = args.TakeChannel.GetListAnd("take_channel", "take_channel", where, maps)
	if err != nil {
		return
	}
	where, maps, err = args.TakeFrom.GetListAnd("take_from", "take_from", where, maps)
	if err != nil {
		return
	}
	where, maps, err = args.CreateInfo.GetListAnd("create_info", "create_info", where, maps)
	if err != nil {
		return
	}
	if args.MinPrice > 0 {
		where = where + " AND price >= :min_price"
		maps["min_price"] = args.MinPrice
	}
	if args.MaxPrice > 0 {
		where = where + " AND price <= :max_price"
		maps["max_price"] = args.MaxPrice
	}
	where, maps = args.TimeBetween.GetBetweenByTimeAnd("create_at", where, maps)
	if args.Params.Mark != "" {
		var paramsJSON []byte
		paramsJSON, err = json.Marshal([]CoreSQLConfig.FieldsConfigType{
			{
				Mark: args.Params.Mark,
				Val:  args.Params.Val,
			},
		})
		if err != nil {
			return
		}
		where = where + " AND params @> :params"
		maps["params"] = string(paramsJSON)
	}
	if args.FailedCode != "" {
		where = where + " AND failed_code = :failed_code"
		maps["failed_code"] = args.FailedCode
	}
	if args.Search != "" {
		where = where + " AND (id ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%' OR key ILIKE '%' || :search || '%' OR failed_message ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "finance_pay"
	if args.IsHistory {
		tableName = tableName + "_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, expire_at, key, status, currency, price, refund_price, refund_send, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des, failed_code, failed_message, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "expire_at", "price", "refund_price", "failed_code"},
	)
	//禁止反馈rand
	for k, v := range dataList {
		for k2, v2 := range v.Params {
			if v2.Mark == "rand" {
				dataList[k].Params[k2].Val = ""
			}
		}
	}
	//反馈
	return
}

// ArgsGetOne 通过任意一个数据查询参数
type ArgsGetOne struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//key
	// key只能获取当前列表内数据，不会从历史表调用
	Key string `json:"key"`
	//是否为敏感模式
	// 敏感模式下将隐藏关键内容
	IsSecret bool `json:"isSecret"`
}

// GetOne 通过任意一个数据查询
func GetOne(args *ArgsGetOne) (data FieldsPayType, err error) {
	if args.ID > 0 {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, key, status, currency, price, refund_price, refund_send, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des, failed_code, failed_message, params FROM finance_pay WHERE id = $1", args.ID)
		if err != nil {
			err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, key, status, currency, price, refund_price, refund_send, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des, failed_code, failed_message, params FROM finance_pay_history WHERE id = $1", args.ID)
		}
	} else {
		if args.Key != "" {
			err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, key, status, currency, price, refund_price, refund_send, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des, failed_code, failed_message, params FROM finance_pay WHERE key = $1", args.Key)
			if err != nil {
				err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, key, status, currency, price, refund_price, refund_send, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des, failed_code, failed_message, params FROM finance_pay_history WHERE key = $1", args.Key)
			}
		} else {
			err = errors.New("id or key is empty")
			return
		}
	}
	if args.IsSecret {
		//禁止反馈rand
		for k2, v2 := range data.Params {
			if v2.Mark == "rand" {
				data.Params[k2].Val = ""
			}
		}
	}
	return
}

// ArgsGetID 通过ID查询数据参数
type ArgsGetID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//是否为敏感模式
	// 敏感模式下将隐藏关键内容
	IsSecret bool `json:"isSecret"`
}

// GetID 通过ID查询数据
func GetID(args *ArgsGetID) (data FieldsPayType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, key, status, currency, price, refund_price, refund_send, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des, failed_code, failed_message, params FROM finance_pay WHERE id = $1", args.ID)
	if err != nil || data.ID < 1 {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, key, status, currency, price, refund_price, refund_send, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des, failed_code, failed_message, params FROM finance_pay_history WHERE id = $1", args.ID)
	}
	if args.IsSecret {
		//禁止反馈rand
		for k2, v2 := range data.Params {
			if v2.Mark == "rand" {
				data.Params[k2].Val = ""
			}
		}
	}
	return
}

// getPayByID 获取支付ID
func getPayByID(id int64) (data FieldsPayType) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, key, status, currency, price, refund_price, refund_send, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des, failed_code, failed_message, params FROM finance_pay WHERE id = $1", id)
	if err != nil || data.ID < 1 {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, key, status, currency, price, refund_price, refund_send, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des, failed_code, failed_message, params FROM finance_pay_history WHERE id = $1", id)
		if err != nil {
			return
		}
	}
	return
}

// ArgsCheckPaymentFrom 检查支付创建来源正确性参数
type ArgsCheckPaymentFrom struct {
	//支付ID
	ID int64 `db:"id" json:"id" check:"id"`
	//收款人或收款渠道信息
	TakeFrom CoreSQLFrom.FieldsFrom `db:"take_from" json:"takeFrom"`
}

// CheckPaymentFrom 检查支付创建来源正确性
func CheckPaymentFrom(args *ArgsCheckPaymentFrom) (err error) {
	type dataType struct {
		//ID
		ID int64 `db:"id" json:"id"`
	}
	takeFrom := CoreSQLFrom.FieldsFromOnlyID{
		System: args.TakeFrom.System,
		ID:     args.TakeFrom.ID,
	}
	var data dataType
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM finance_pay WHERE id = $1 AND (payment_create @> $2 OR payment_from @> $2)", args.ID, takeFrom)
	if err == nil && data.ID > 0 {
		return
	}
	if err == nil {
		err = errors.New("no data")
	}
	return
}

// ArgsCheckTakeFrom 检查支付请求的渠道是否符合条件参数
type ArgsCheckTakeFrom struct {
	//支付ID
	ID int64 `db:"id" json:"id" check:"id"`
	//收款人或收款渠道信息
	TakeFrom CoreSQLFrom.FieldsFrom `db:"take_from" json:"takeFrom"`
}

// CheckTakeFrom 检查支付请求的渠道是否符合条件
func CheckTakeFrom(args *ArgsCheckTakeFrom) (err error) {
	type dataType struct {
		//ID
		ID int64 `db:"id" json:"id"`
	}
	takeFrom := CoreSQLFrom.FieldsFromOnlyID{
		System: args.TakeFrom.System,
		ID:     args.TakeFrom.ID,
	}
	var data dataType
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM finance_pay WHERE id = $1 AND (take_create @> $2 OR take_from @> $2)", args.ID, takeFrom)
	if err == nil && data.ID > 0 {
		return
	}
	if err == nil {
		err = errors.New("no data")
	}
	return
}

// makeShortKey 生成新的key
func makeShortKey(tryCount int) (string, error) {
	if tryCount > 50 {
		return "", errors.New("try get short key too many, pls take more len")
	}
	newKey, err := CoreFilter.GetRandStr3(shortKeyLen)
	if err != nil {
		return "", err
	}
	var data FieldsPayType
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM finance_pay WHERE key = $1", newKey)
	if err != nil {
		return newKey, nil
	} else {
		tryCount += 1
		return makeShortKey(tryCount)
	}
}

// saveFinanceLog 记录finance log
func saveFinanceLog(status int, createInfo CoreSQLFrom.FieldsFrom, data *FieldsPayType) error {
	return FinanceLog.Create(&FinanceLog.ArgsLogCreate{
		PayID:          data.ID,
		Hash:           "",
		Key:            data.Key,
		Status:         status,
		Currency:       data.Currency,
		Price:          data.Price,
		PaymentCreate:  data.PaymentCreate,
		PaymentChannel: data.PaymentChannel,
		PaymentFrom:    data.PaymentFrom,
		TakeCreate:     data.TakeCreate,
		TakeChannel:    data.TakeChannel,
		TakeFrom:       data.TakeFrom,
		CreateInfo:     createInfo,
		Des:            data.Des,
	})
}

// checkPaySystem 检查支付system
func checkPaySystem(paySystem string) error {
	switch paySystem {
	case "cash":
	case "deposit":
	case "weixin":
	case "alipay":
	case "paypal":
	case "company_returned":
	default:
		return errors.New("pay system is error")
	}
	return nil
}

// checkStatus 检查交易状态
func checkStatus(status int) error {
	if status > -1 && status < 10 {
		return nil
	}
	return errors.New("status is error")
}

// changeDeposit 给目标账户变动资金
func changeDeposit(depositID int64, mark string, price int64) (errCode string, err error) {
	var depositData FinanceDeposit.FieldsDepositType
	depositData, err = FinanceDeposit.GetByID(&FinanceDeposit.ArgsGetByID{
		ID: depositID,
	})
	if err != nil {
		err = errors.New("deposit not exist, " + err.Error())
		return
	}
	depositData, errCode, err = FinanceDeposit.SetByID(&FinanceDeposit.ArgsSetByID{
		UpdateHash:      depositData.UpdateHash,
		ID:              depositData.ID,
		ConfigMark:      mark,
		AppendSavePrice: price,
	})
	if err != nil {
		err = errors.New("deposit update, " + err.Error())
		return
	}
	return
}

// argsUpdateStatus 更新支付状态参数
type argsUpdateStatus struct {
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom
	//ID
	ID int64
	//key
	Key string
	//验证上一个状态必须是
	// 可以指定多个，如果为空则不验证
	PrevStatus []int
	//新状态
	Status int
	//堆叠写入数据
	// a = :a
	// 如果存在:，则后续必须给予maps，否则将报错
	// 注意不要使用status和params，该两个为占位内容
	SetQuery string
	SetMaps  map[string]interface{}
	//补充扩展
	Params []CoreSQLConfig.FieldsConfigType
}

// updateStatus 更新支付状态
func updateStatus(args *argsUpdateStatus) (errCode string, err error) {
	//获取交易数据
	var data FieldsPayType
	data, err = GetOne(&ArgsGetOne{
		ID:  args.ID,
		Key: args.Key,
	})
	if err != nil {
		errCode = "data_not_exist"
		return
	}
	if args.PrevStatus != nil && len(args.PrevStatus) > 0 {
		isFind := false
		for _, v := range args.PrevStatus {
			if v == data.Status {
				isFind = true
				break
			}
		}
		if !isFind {
			errCode = "status_not_support"
			err = errors.New("status not support")
			return
		}
	}
	//更新状态
	if args.Params != nil && len(args.Params) > 0 {
		for _, v := range args.Params {
			data.Params = CoreSQLConfig.Set(data.Params, v.Mark, v.Val)
		}
	}
	maps := map[string]interface{}{
		"id":     data.ID,
		"status": args.Status,
		"params": data.Params,
	}
	if args.SetQuery != "" {
		for k, v := range args.SetMaps {
			maps[k] = v
		}
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_pay SET status = :status, params = :params "+args.SetQuery+" WHERE id = :id", maps)
	if err != nil {
		errCode = "update"
		err = errors.New("update by id, " + err.Error())
		return
	}
	err = saveFinanceLog(args.Status, args.CreateInfo, &data)
	if err != nil {
		CoreLog.Info("update status ", args.Status, ", create finance log, ", err)
		err = nil
	}
	//根据状态推送nats
	switch args.Status {
	case 2:
		//推送nats
		CoreNats.PushDataNoErr("/finance/pay/failed", "failed", data.ID, "", nil)
	case 4:
		//推送nats
		CoreNats.PushDataNoErr("/finance/pay/failed", "remove", data.ID, "", nil)
	case 5:
		//推送nats
		CoreNats.PushDataNoErr("/finance/pay/failed", "expire", data.ID, "", nil)
	}
	return
}
