package FinanceLog

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime2 "github.com/fotomxq/weeekj_core/v5/core/sql/time2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

//财务模块
// 提供各类支付汇总、储蓄、交易结构体方案
// 外部需要存在一套审计维护程序，该审计程序将反向获取对应模块数据并验证是否匹配

var (
	//系统用混淆计算的hash盐
	appendHash string
)

// SetHash 设置混淆值
func SetHash(hash string) {
	appendHash = hash
}

func GetHash() string {
	return appendHash
}

// ArgsGetList 查看记录列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//交易短key
	// 在历史表中，该值可能发生重复，请勿以该值作为最终唯一判断
	// 用于微信、支付宝等接口对接时，采用的短Key处理机制
	Key string `db:"key" json:"key"`
	//最终状态 必须填写
	// wait 客户端发起付款，并正在支付中
	// client 客户端完成支付，等待服务端验证
	// failed 交易失败，服务端主动取消交易或其他原因取消交易
	// finish 交易成功
	// remove 交易销毁
	// expire 交易过期
	// refund 发起退款申请
	// refundAudit 退款审核通过，等待处理中
	// refundFailed 退款失败
	// refundFinish 退款完成
	Status []int `bson:"Status" json:"status"`
	//付款人来源
	PaymentCreate CoreSQLFrom.FieldsFrom `json:"paymentCreate"`
	//支付方式
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	PaymentChannel CoreSQLFrom.FieldsFrom `json:"paymentChannel"`
	//支付方的来源
	// 留空则代表平台方，否则为商户或加盟商
	PaymentFrom CoreSQLFrom.FieldsFrom `json:"paymentFrom"`
	//收款人来源
	TakeCreate CoreSQLFrom.FieldsFrom `json:"takeCreate"`
	//收款渠道
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	TakeChannel CoreSQLFrom.FieldsFrom `json:"takeChannel"`
	//收款方来源
	// 留空则代表平台方，否则为商户或加盟商
	TakeFrom CoreSQLFrom.FieldsFrom `json:"takeFrom"`
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom `json:"createInfo"`
	//时间段
	TimeBetween CoreSQLTime2.FieldsCoreTime `json:"time_between"`
	//是否为历史
	IsHistory bool `json:"is_history"`
	//搜索
	Search string `json:"search"`
}

// GetList 查看记录列表
func GetList(args *ArgsGetList) (dataList []FieldsLogType, dataCount int64, err error) {
	if len(args.Status) < 1 {
		args.Status = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	}
	where := "status = ANY(:status)"
	maps := map[string]interface{}{
		"status": pq.Array(args.Status),
	}
	if args.Key != "" {
		where = where + " AND key = :key"
		maps["key"] = args.Key
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
	var newWhere string
	newWhere, maps = args.TimeBetween.GetBetweenByTime("create_at", maps)
	if newWhere != "" {
		where = where + " AND " + newWhere
	}
	if args.Search != "" {
		where = where + " AND (des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "finance_log"
	if args.IsHistory {
		tableName = tableName + "_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, pay_id, hash, key, status, currency, price, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "price"},
	)
	if err != nil {
		return
	}
	return
}

// ArgsLogCreate 创建新的记录
type ArgsLogCreate struct {
	//支付渠道信息ID
	PayID int64 `db:"pay_id" json:"payID"`
	//混淆验证
	Hash string `db:"hash" json:"hash"`
	//交易短key
	// 在历史表中，该值可能发生重复，请勿以该值作为最终唯一判断
	// 用于微信、支付宝等接口对接时，采用的短Key处理机制
	Key string `db:"key" json:"key"`
	//最终状态
	// wait 客户端发起付款，并正在支付中
	// client 客户端完成支付，等待服务端验证
	// failed 交易失败，服务端主动取消交易或其他原因取消交易
	// finish 交易成功
	// remove 交易销毁
	// expire 交易过期
	// refund 发起退款申请
	// refundAudit 退款审核通过，等待处理中
	// refundFailed 退款失败
	// refundFinish 退款完成
	Status int `db:"status" json:"status"`
	//交易货币类型
	// 采用CoreCurrency匹配
	Currency int `db:"currency" json:"currency"`
	//交易金额
	Price int64 `db:"price" json:"price"`
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
	//操作原因
	Des string `db:"des" json:"des"`
}

func Create(args *ArgsLogCreate) (err error) {
	//hash加盐
	hashStr := fmt.Sprint(args.PayID, args.Key, args.Currency, args.Price, args.PaymentCreate, args.PaymentFrom, args.PaymentChannel, args.TakeCreate, args.TakeFrom, args.TakeChannel, appendHash)
	args.Hash, err = CoreFilter.GetSha256Str(hashStr)
	if err != nil {
		return err
	}
	//写入数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_log (pay_id, hash, key, status, currency, price, payment_create, payment_channel, payment_from, take_create, take_channel, take_from, create_info, des) VALUES (:pay_id,:hash,:key,:status,:currency,:price,:payment_create,:payment_channel,:payment_from,:take_create,:take_channel,:take_from,:create_info,:des)", args)
	if err != nil {
		return
	}
	CoreNats.PushDataNoErr("/finance/log/file", "", 0, "", nil)
	return
}
