package FinanceSafe

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//安全审计模块
// 根据Log对相关模块进行审计工作

// ArgsGetList 获取安全日志参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//造成该事件的来源
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//交易发生双方
	//来源，支付方
	PaymentCreate CoreSQLFrom.FieldsFrom `db:"payment_create" json:"paymentCreate"`
	//支付方的来源
	// 留空则代表平台方，否则为商户或加盟商
	PaymentFrom CoreSQLFrom.FieldsFrom `db:"payment_from" json:"paymentFrom"`
	//目标，接收方
	TakeCreate CoreSQLFrom.FieldsFrom `db:"take_create" json:"takeCreate"`
	//收款方来源
	// 留空则代表平台方，否则为商户或加盟商
	TakeFrom CoreSQLFrom.FieldsFrom `db:"take_from" json:"takeFrom"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//日志ID
	PayLogID int64 `db:"pay_log_id" json:"payLogID"`
	//是否需要预警参数
	NeedAllowEW bool
	//是否启动预警参数
	AllowEW bool
	//错误代码
	Code string
	//是否为打开状态
	AllowOpen bool
	//搜索
	Search string
}

// GetList 获取安全日志
func GetList(args *ArgsGetList) (dataList []FieldsSafeType, dataCount int64, err error) {
	where := "allow_open = :allow_open"
	maps := map[string]interface{}{
		"allow_open": args.AllowOpen,
	}
	where, maps, err = args.CreateInfo.GetListAnd("create_info", "create_info", where, maps)
	if err != nil {
		return
	}
	where, maps, err = args.PaymentCreate.GetListAnd("payment_create", "payment_create", where, maps)
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
	where, maps, err = args.TakeFrom.GetListAnd("take_from", "take_from", where, maps)
	if err != nil {
		return
	}
	if args.NeedAllowEW {
		where = where + " AND allow_ew = :allow_ew"
		maps["allow_ew"] = args.AllowEW
	}
	if args.PayID > 0 {
		where = where + " AND pay_id = :pay_id"
		maps["pay_id"] = args.PayID
	}
	if args.PayLogID > 0 {
		where = where + " AND pay_log_id = :pay_log_id"
		maps["pay_log_id"] = args.PayLogID
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"finance_safe",
		"id",
		"SELECT id, create_at, payment_create, payment_from, take_create, take_from, pay_id, pay_log_id, message, code, need_ew, allow_ew, ew_template_mark, allow_open FROM finance_safe WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "code", "ew_template_mark"},
	)
	return
}

// ArgsUpdateDone 标记日志为已处理参数
type ArgsUpdateDone struct {
	//ID
	ID int64
	//目标，接收方
	TakeCreate CoreSQLFrom.FieldsFrom `db:"take_create" json:"takeCreate"`
	//收款方来源
	// 留空则代表平台方，否则为商户或加盟商
	TakeFrom CoreSQLFrom.FieldsFrom `db:"take_from" json:"takeFrom"`
}

// UpdateDone 标记日志为已处理
// 可以验证接收方及来源
func UpdateDone(args *ArgsUpdateDone) (err error) {
	where := "id = :id"
	maps := map[string]interface{}{
		"id": args.ID,
	}
	var newWhere string
	newWhere, maps, err = args.TakeCreate.GetList("take_create", "take_create", maps)
	if err != nil {
		return
	} else {
		where = where + " AND " + newWhere
	}
	newWhere, maps, err = args.TakeFrom.GetList("take_from", "take_from", maps)
	if err != nil {
		return
	} else {
		where = where + " AND " + newWhere
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_safe SET allow_open = true WHERE "+where, maps)
	return
}

// checkHaveRecord 检查某个数据是否存在异常记录？
func checkHaveRecord(createInfo CoreSQLFrom.FieldsFrom, code string) (b bool) {
	var data FieldsSafeType
	createInfoData, err := createInfo.GetRaw()
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, allow_open FROM finance_safe WHERE create_info @> $1 AND allow_open = true AND code = $2", createInfoData, code)
	if err != nil {
		return
	}
	return data.AllowOpen
}

// argsCreateRecord 记录一个异常点参数
type argsCreateRecord struct {
	//造成该事件的来源
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//交易发生双方
	//来源，支付方
	PaymentCreate CoreSQLFrom.FieldsFrom `db:"payment_create" json:"paymentCreate"`
	//支付方的来源
	// 留空则代表平台方，否则为商户或加盟商
	PaymentFrom CoreSQLFrom.FieldsFrom `db:"payment_from" json:"paymentFrom"`
	//目标，接收方
	TakeCreate CoreSQLFrom.FieldsFrom `db:"take_create" json:"takeCreate"`
	//收款方来源
	// 留空则代表平台方，否则为商户或加盟商
	TakeFrom CoreSQLFrom.FieldsFrom `db:"take_from" json:"takeFrom"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//日志ID
	PayLogID int64 `db:"pay_log_id" json:"payLogID"`
	//安全事件详细描述信息
	Message string `db:"message" json:"message"`
	//安全标识码
	// 用于其他语言翻译或前端传输
	Code string `db:"code" json:"code"`
	//是否需要发出预警消息
	NeedEW bool `db:"needEW" json:"needEW"`
}

// createRecord 记录一个异常点
func createRecord(args *argsCreateRecord) (err error) {
	configMark := ""
	switch args.Code {
	case "LogHash":
		configMark = "FinanceSafeEWTemplateByLogHash"
	case "PayLost":
		configMark = "FinanceSafeEWTemplateByPayLost"
	case "PayPrice":
		configMark = "FinanceSafeEWTemplateByPayPrice"
	case "PayLimit0":
		configMark = "FinanceSafeEWTemplateByPayLimit0"
	case "PayLimitMax":
		configMark = "FinanceSafeEWTemplateByPayLimitMax"
	case "PayFrequencyOneFrom":
		configMark = "FinanceSafeEWTemplateByPayFrequencyOneFrom"
	case "PayFrequencyOneTo":
		configMark = "FinanceSafeEWTemplateByPayFrequencyOneTo"
	case "PayFrequencyAll":
		configMark = "FinanceSafeEWTemplateByPayFrequencyAll"
	default:
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_safe (create_info, payment_create, payment_from, take_create, take_from, pay_id, pay_log_id, message, code, need_ew, allow_ew, ew_template_mark, allow_open) VALUES (:create_info, :payment_create, :payment_from, :take_create, :take_from, :pay_id, :pay_log_id, :message, :code, :need_ew, false, :ew_template_mark, true)", map[string]interface{}{
		"create_info":      args.CreateInfo,
		"payment_create":   args.PaymentCreate,
		"payment_from":     args.PaymentFrom,
		"take_create":      args.TakeCreate,
		"take_from":        args.TakeFrom,
		"pay_id":           args.PayID,
		"pay_log_id":       args.PayLogID,
		"message":          args.Message,
		"code":             args.Code,
		"need_ew":          args.NeedEW,
		"ew_template_mark": configMark,
	})
	return
}
