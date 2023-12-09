package TMSTransport

import (
	"errors"
	"fmt"
	ClassConfig "gitee.com/weeekj/weeekj_core/v5/class/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	OrgCoreCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceOrderMod "gitee.com/weeekj/weeekj_core/v5/service/order/mod"
	"time"
)

// ArgsCreateTransport 创建新配送单参数
type ArgsCreateTransport struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//当前配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//取货地址
	FromAddress CoreSQLAddress.FieldsAddress `db:"from_address" json:"fromAddress"`
	//收货地址
	ToAddress CoreSQLAddress.FieldsAddress `db:"to_address" json:"toAddress"`
	//订单ID
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//货物ID
	Goods FieldsTransportGoods `db:"goods" json:"goods"`
	//快递总重量
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//长宽
	Length int `db:"length" json:"length" check:"intThan0" empty:"true"`
	Width  int `db:"width" json:"width" check:"intThan0" empty:"true"`
	//货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	//配送费用
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//是否完成了缴费
	PayFinish bool `db:"pay_finish" json:"payFinish" check:"bool" empty:"true"`
	//期望送货时间
	TaskAt string `db:"task_at" json:"taskAt" check:"isoTime" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateTransport 创建新配送单
func CreateTransport(args *ArgsCreateTransport) (data FieldsTransport, errCode string, err error) {
	//检查货物
	if len(args.Goods) < 1 {
		errCode = "err_no_goods"
		err = errors.New("goods is empty")
		return
	}
	//修正配送费缴费状态
	if args.Price < 1 {
		args.PayFinish = true
	}
	//预期上门时间
	var taskAt time.Time
	taskAt, err = CoreFilter.GetTimeByISO(args.TaskAt)
	if err != nil {
		taskAt = time.Now()
	}
	//获取SN信息
	var snData FieldsTransport
	err = Router2SystemConfig.MainDB.Get(&snData, "SELECT id, create_at, sn, sn_day FROM tms_transport WHERE org_id = $1 ORDER BY id DESC LIMIT 1", args.OrgID)
	if err != nil {
		snData = FieldsTransport{}
		snData.SN = 1
		snData.SNDay = 1
	} else {
		snData.SN += 1
		if snData.CreateAt.Unix() >= CoreFilter.GetNowTimeCarbon().StartOfDay().Time.Unix() {
			snData.SNDay += 1
		} else {
			snData.SNDay = 1
		}
	}
	//是否为退货单
	isRefund := false
	isRefund, _ = args.Params.GetValBool("isRefund")
	//初始化数据
	var waitLog []argsAppendLog
	//如果没有分配配送人员，则直接自动分配
	if args.BindID < 1 {
		waitLog, args.Params, args.BindID = transportSelectBind(isRefund, args.OrgID, args.UserID, args.BindID, args.Goods, args.FromAddress, args.ToAddress, args.Params)
	}
	//检查配送员ID
	if args.BindID > 0 {
		_, err = GetBindByBindID(&ArgsGetBindByBindID{
			OrgID:  args.OrgID,
			BindID: args.BindID,
		})
		if err != nil {
			errCode = "err_no_service_org_bind"
			err = errors.New(fmt.Sprint("get bind by bind id, ", err))
			return
		}
	}
	//修订状态信息
	status := 0
	// 商户全局自动跳过取货状态的设置
	// 注意，在isRefund状态下，将失效
	transportAllowSkipPick := false
	transportAllowSkipPick, _ = OrgCoreCore.Config.GetConfigValBool(&ClassConfig.ArgsGetConfig{
		BindID:    args.OrgID,
		Mark:      "TransportAllowSkipPick",
		VisitType: "admin",
	})
	if args.BindID > 0 && transportAllowSkipPick && !isRefund {
		status = 2
	}
	//修正支付状态
	var payFinishAt time.Time
	if args.PayFinish && args.Price < 1 {
		payFinishAt = CoreFilter.GetNowTime()
		paySystem, b := args.Params.GetVal("paySystem")
		if !b || paySystem == "" {
			if args.Price < 1 {
				args.Params = CoreSQLConfig.Set(args.Params, "paySystem", "free")
			} else {
				args.Params = CoreSQLConfig.Set(args.Params, "paySystem", "order")
			}
		}
	} else {
		args.Params = CoreSQLConfig.Set(args.Params, "paySystem", "unkonw")
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tms_transport", "INSERT INTO tms_transport (create_at, org_id, bind_id, info_id, user_id, sn, sn_day, status, from_address, to_address, order_id, goods, weight, length, width, currency, price, pay_finish_at, pay_id, task_at, params) VALUES (:create_at, :org_id,:bind_id,:info_id,:user_id,:sn,:sn_day,:status,:from_address,:to_address,:order_id,:goods,:weight,:length,:width,:currency,:price,:pay_finish_at,0,:task_at,:params)", map[string]interface{}{
		"create_at":     CoreFilter.GetNowTime(),
		"org_id":        args.OrgID,
		"bind_id":       args.BindID,
		"info_id":       args.InfoID,
		"user_id":       args.UserID,
		"sn":            snData.SN,
		"sn_day":        snData.SNDay,
		"status":        status,
		"from_address":  args.FromAddress,
		"to_address":    args.ToAddress,
		"order_id":      args.OrderID,
		"goods":         args.Goods,
		"weight":        args.Weight,
		"length":        args.Length,
		"width":         args.Width,
		"currency":      args.Currency,
		"price":         args.Price,
		"pay_finish_at": payFinishAt,
		"task_at":       taskAt,
		"params":        args.Params,
	}, &data)
	if err != nil {
		errCode = "err_insert"
		err = errors.New(fmt.Sprint("insert data, ", err))
		return
	}
	//写入日志
	waitLog = append(waitLog, argsAppendLog{
		OrgID:           args.OrgID,
		BindID:          args.BindID,
		TransportID:     data.ID,
		TransportBindID: data.BindID,
		Mark:            "create",
		Des:             fmt.Sprint("创建配送单"),
	})
	for _, v := range waitLog {
		v.TransportID = data.ID
		_ = appendLog(&v)
	}
	//生成统计数据
	err = updateTransportBindAnalysis(data, data.BindID)
	if err != nil {
		CoreLog.Warn("update transport bind analysis failed, ", err)
		err = nil
	}
	if data.BindID > 0 {
		pushNatsAnalysisBind(data.BindID)
	}
	//推送配送单MQTT更新
	pushMQTTTransportUpdate(data.OrgID, data.BindID, data.ID, 0)
	//推送订单更新配送单请求
	if data.OrderID > 0 {
		//通知订单创建了服务单
		ServiceOrderMod.UpdateTransportID(ServiceOrderMod.ArgsUpdateTransportID{
			TMSType:     "transport",
			ID:          args.OrderID,
			SN:          data.SN,
			SNDay:       data.SNDay,
			Des:         fmt.Sprint("生成配送单ID[", data.ID, "]，SN[", data.SN, "]，当日SN[", data.SNDay, "]"),
			TransportID: data.ID,
		})
	}
	//反馈
	return
}
