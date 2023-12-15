package ServiceHousekeeping

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsCreateLog 创建新的请求参数
type ArgsCreateLog struct {
	//服务用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//预计上门时间
	NeedAt time.Time `db:"need_at" json:"needAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//其他服务参与人员
	OtherBinds pq.Int64Array `db:"other_binds" json:"otherBinds" check:"ids" empty:"true"`
	//服务项目商品ID
	MallProductID int64 `db:"mall_product_id" json:"mallProductID" check:"id"`
	//订单ID
	OrderID int64 `db:"order_id" json:"orderID" check:"id"`
	//货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	//单位报价
	// 服务负责人在收款前可以协商变更
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//是否支付
	PayAt time.Time `db:"pay_at" json:"payAt"`
	//客户备注
	Des string `db:"des" json:"des" check:"des" min:"1" max:"500" empty:"true"`
	//客户地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//服务配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateLog 创建新的请求
func CreateLog(args *ArgsCreateLog) (data FieldsLog, errCode string, err error) {
	for _, v := range args.OtherBinds {
		if v == args.BindID {
			errCode = "err_org_bind"
			err = errors.New("bind id in other bind ids")
			return
		}
	}
	if args.Price < 1 {
		args.Price = 0
		args.PayAt = CoreFilter.GetNowTime()
	}
	if args.BindID < 1 {
		var bindData FieldsBind
		bindData, err = getBindByMarketUserID(args.OrgID, args.UserID)
		if err != nil {
			err = nil
		} else {
			args.BindID = bindData.BindID
		}
	}
	//检查预约时间
	var needAt time.Time
	if args.ConfigID > 0 {
		if configData, b := CheckConfigTime(&ArgsCheckConfigTime{
			ID: args.ConfigID,
		}); !b {
			errCode = "err_time"
			err = errors.New("no time")
			return
		} else {
			needAt = configData.StartAt
		}
	} else {
		needAt = args.NeedAt
	}
	var sn, snDay int64
	sn, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "service_housekeeping_log", "id", "org_id = :org_id", map[string]interface{}{
		"org_id": args.OrgID,
	})
	if err != nil {
		err = nil
		sn = 0
	} else {
		sn += 1
	}
	snDay, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "service_housekeeping_log", "id", "org_id = :org_id AND create_at >= :start_at AND create_at <= :end_at", map[string]interface{}{
		"org_id":   args.OrgID,
		"start_at": CoreFilter.GetNowTimeCarbon().StartOfDay().StartOfHour().StartOfMinute().Time,
		"end_at":   CoreFilter.GetNowTimeCarbon().EndOfDay().EndOfHour().EndOfMinute().Time,
	})
	if err != nil {
		err = nil
		snDay = 0
	} else {
		snDay += 1
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_housekeeping_log", "INSERT INTO service_housekeeping_log(sn, sn_day, need_at, user_id, org_id, bind_id, other_binds, mall_product_id, order_id, currency, price, pay_at, des, address, config_id, params) VALUES(:sn, :sn_day, :need_at, :user_id, :org_id, :bind_id, :other_binds, :mall_product_id, :order_id, :currency, :price, :pay_at, :des, :address, :config_id, :params)", map[string]interface{}{
		"sn":              sn,
		"sn_day":          snDay,
		"user_id":         args.UserID,
		"need_at":         needAt,
		"org_id":          args.OrgID,
		"bind_id":         args.BindID,
		"other_binds":     args.OtherBinds,
		"mall_product_id": args.MallProductID,
		"order_id":        args.OrderID,
		"currency":        args.Currency,
		"price":           args.Price,
		"pay_at":          args.PayAt,
		"des":             args.Des,
		"address":         args.Address,
		"config_id":       args.ConfigID,
		"params":          args.Params,
	}, &data)
	if err != nil {
		errCode = "err_insert"
		return
	}
	//创建成员的统计数据
	if args.BindID > 0 {
		_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_bind SET update_at = NOW(), all_log_count = all_log_count + 1, un_finish_count = un_finish_count + 1 WHERE bind_id = :bind_id", map[string]interface{}{
			"bind_id": args.BindID,
		})
	}
	//更新配置统计
	_ = addConfig(args.ConfigID)
	//反馈
	return
}
