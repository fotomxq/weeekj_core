package TMSUserRunning

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceOrderMod "github.com/fotomxq/weeekj_core/v5/service/order/mod"
	"github.com/lib/pq"
	"time"
)

// ArgsCreateMission 创建新的任务参数
type ArgsCreateMission struct {
	//跑腿单类型
	// 0 帮我送 ; 1 帮我买; 2 帮我取
	RunType int `db:"run_type" json:"runType" check:"intThan0" empty:"true"`
	//期望上门时间
	WaitAt string `db:"wait_at" json:"waitAt" check:"isoTime"`
	//物品类型
	GoodType string `db:"good_type" json:"goodType" check:"mark"`
	//关联组织
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//关联用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//关联订单ID
	// 可能没有关联订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//商品等待缴纳费用
	OrderWaitPrice int64 `db:"order_wait_price" json:"orderWaitPrice" check:"price" empty:"true"`
	//等待缴纳的费用
	RunWaitPrice int64 `db:"run_wait_price" json:"runWaitPrice" check:"price" empty:"true"`
	//跑腿费是否货到付款
	RunPayAfter bool `db:"run_pay_after" json:"runPayAfter" check:"bool"`
	//订单是否已经缴纳了所有费用
	OrderPayAllPrice bool `json:"orderPayAllPrice" check:"bool"`
	//商品内容描述
	OrderDes string `db:"order_des" json:"orderDes" check:"des" min:"1" max:"3000" empty:"true"`
	//跑腿单描述信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1200" empty:"true"`
	//物品重量
	GoodWidget int `db:"good_widget" json:"goodWidget" check:"intThan0" empty:"true"`
	//发货地址
	FromAddress CoreSQLAddress.FieldsAddress `db:"from_address" json:"fromAddress" check:"address_data" empty:"true"`
	//送货地址
	ToAddress CoreSQLAddress.FieldsAddress `db:"to_address" json:"toAddress" check:"address_data" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params" empty:"true"`
}

// CreateMission 创建新的任务
func CreateMission(args *ArgsCreateMission) (data FieldsMission, err error) {
	//修正参数
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//跑腿费用不能少于1
	if args.RunWaitPrice == 0 {
		err = errors.New("run wait price less 1")
		return
	}
	//自动计算跑腿费用
	if args.RunWaitPrice < 0 {
		sumData := GetRunPrice(&ArgsGetRunPrice{
			WaitAt:      args.WaitAt,
			GoodType:    args.GoodType,
			GoodWidget:  args.GoodWidget,
			FromAddress: args.FromAddress,
			ToAddress:   args.ToAddress,
		})
		args.Params = CoreSQLConfig.Set(args.Params, "betweenM", sumData.BetweenM)
		args.Params = CoreSQLConfig.Set(args.Params, "betweenKM", sumData.BetweenKM)
		args.Params = CoreSQLConfig.Set(args.Params, "betweenPrice", sumData.BetweenPrice)
		args.Params = CoreSQLConfig.Set(args.Params, "waitPrice", sumData.WaitPrice)
		args.Params = CoreSQLConfig.Set(args.Params, "widgetPrice", sumData.WidgetPrice)
	}
	//构建提取代码
	var takeCode string
	takeCode, err = CoreFilter.GetRandStr3(6)
	if err != nil {
		return
	}
	//跑腿单缴费时间
	runPayAt := time.Time{}
	var runPrice int64
	//如果存在订单，则获取订单数据
	var orderWaitPrice int64
	var orderData ServiceOrderMod.FieldsOrder
	var orderPrice int64
	if args.OrderID > 0 {
		orderData, err = ServiceOrderMod.GetByID(&ServiceOrderMod.ArgsGetByID{
			ID:     args.OrderID,
			OrgID:  -1,
			UserID: -1,
		})
		if err != nil {
			return
		}
		if orderData.TransportPayAfter {
			orderWaitPrice = orderData.PriceTotal
		}
	} else {
		orderData.Des = args.OrderDes
		orderData.PriceTotal = args.OrderWaitPrice
		orderWaitPrice = args.OrderWaitPrice
	}
	//订单支付时间
	var orderPayAt time.Time
	if orderData.PayStatus == 1 {
		orderPrice = orderData.PriceTotal
		orderWaitPrice = 0
		orderPayAt = CoreFilter.GetNowTime()
		if args.OrderPayAllPrice {
			args.RunWaitPrice = 0
			for _, v := range orderData.PriceList {
				if v.PriceType == 1 {
					runPrice = v.Price
				}
			}
			runPayAt = CoreFilter.GetNowTime()
			orderPrice = orderData.PriceTotal - runPrice
		}
	}
	//获取上门时间
	var waitAt time.Time
	waitAt, err = CoreFilter.GetTimeByISO(args.WaitAt)
	if err != nil {
		waitAt = CoreFilter.GetNowTime()
	}
	//创建新的任务
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tms_user_running_mission", "INSERT INTO tms_user_running_mission (wait_at, good_type, take_code, run_type, org_id, user_id, order_id, role_id, run_pay_at, run_pay_id, run_price, run_wait_price, run_pay_list, run_pay_after, order_pay_after, order_wait_price, order_price, order_pay_at, order_pay_id, service_price, des, order_des_files, order_des, good_widget, from_address, to_address, logs, params) VALUES (:wait_at,:good_type,:take_code,:run_type,:org_id,:user_id,:order_id,:role_id,:run_pay_at,0,:run_price,:run_wait_price,:run_pay_list,:run_pay_after,:order_pay_after,:order_wait_price,:order_price,:order_pay_at,:order_pay_id,:service_price,:des,:order_des_files,:order_des,:good_widget,:from_address,:to_address,:logs,:params)", map[string]interface{}{
		"take_code":        takeCode,
		"wait_at":          waitAt,
		"good_type":        args.GoodType,
		"run_type":         args.RunType,
		"org_id":           args.OrgID,
		"user_id":          args.UserID,
		"order_id":         args.OrderID,
		"role_id":          0,
		"run_pay_at":       runPayAt,
		"run_price":        runPrice,
		"run_wait_price":   args.RunWaitPrice,
		"run_pay_list":     pq.Int64Array{},
		"run_pay_after":    args.RunPayAfter,
		"order_pay_after":  orderData.TransportPayAfter,
		"order_wait_price": orderWaitPrice,
		"order_price":      orderPrice,
		"order_pay_at":     orderPayAt,
		"order_pay_id":     orderData.PayID,
		"service_price":    0,
		"des":              args.Des,
		"order_des_files":  pq.Int64Array{},
		"order_des":        orderData.Des,
		"good_widget":      args.GoodWidget,
		"from_address":     args.FromAddress,
		"to_address":       args.ToAddress,
		"logs":             FieldsMissionLogs{},
		"params":           args.Params,
	}, &data)
	//通知订单关联
	if args.OrderID > 0 {
		ServiceOrderMod.UpdateTransportID(ServiceOrderMod.ArgsUpdateTransportID{
			TMSType:     "running",
			ID:          args.OrderID,
			SN:          0,
			SNDay:       0,
			Des:         fmt.Sprint("生成跑腿单ID[", data.ID, "]"),
			TransportID: data.ID,
		})
	}
	//反馈
	return
}
