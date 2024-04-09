package ServiceHousekeeping

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceOrderMod "github.com/fotomxq/weeekj_core/v5/service/order/mod"
	"github.com/lib/pq"
	"time"
)

// ArgsUpdateLogBind 更换服务人员参数
type ArgsUpdateLogBind struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//新的服务负责人
	NewBindID int64 `db:"new_bind_id" json:"newBindID" check:"id"`
	//其他服务参与人员
	OtherBinds pq.Int64Array `db:"other_binds" json:"otherBinds" check:"ids" empty:"true"`
}

// UpdateLogBind 更换服务人员
func UpdateLogBind(args *ArgsUpdateLogBind) (err error) {
	for _, v := range args.OtherBinds {
		if v == args.NewBindID {
			err = errors.New("bind id in other bind ids")
			return
		}
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_log SET update_at = NOW(), bind_id = :new_bind_id, other_binds = :other_binds WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:bind_id < 1 OR bind_id = :bind_id)", args)
	return
}

// ArgsUpdateLogOldToNewBind 批量更改未完成服务单服务人员参数
type ArgsUpdateLogOldToNewBind struct {
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//旧配送人员
	OldBindID int64 `db:"old_bind_id" json:"oldBindID" check:"id"`
	//新配送员
	NewBindID int64 `db:"new_bind_id" json:"newBindID" check:"id"`
}

// UpdateLogOldToNewBind 批量更改未完成服务单服务人员
func UpdateLogOldToNewBind(args *ArgsUpdateLogOldToNewBind) (err error) {
	//禁止自己转自己
	if args.OldBindID == args.NewBindID {
		err = errors.New("old and new is same")
		return
	}
	//修改新的服务人员
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_log SET update_at = NOW(), bind_id = :new_bind_id WHERE bind_id = :old_bind_id AND delete_at < to_timestamp(1000000) AND finish_at < to_timestamp(1000000) AND org_id = :org_id", args)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateLogPrice 变更价格参数
type ArgsUpdateLogPrice struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//费用组成
	PriceList ServiceOrderMod.FieldsPrices `db:"price_list" json:"priceList"`
	//单位报价
	// 服务负责人在收款前可以协商变更
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
}

// UpdateLogPrice 变更价格
func UpdateLogPrice(args *ArgsUpdateLogPrice) (err error) {
	if args.Price < 0 {
		err = errors.New("price less 0")
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_log SET update_at = NOW(), price = :price WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:bind_id < 1 OR bind_id = :bind_id) AND pay_at < to_timestamp(1000000)", args)
	if err == nil {
		var data FieldsLog
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, order_id FROM service_housekeeping_log WHERE id = $1", args.ID)
		if err != nil {
			return
		}
		if data.OrderID > 0 {
			ServiceOrderMod.UpdatePrice(ServiceOrderMod.ArgsUpdatePrice{
				ID:        data.OrderID,
				OrgID:     data.OrgID,
				OrgBindID: args.BindID,
				PriceList: args.PriceList,
			})
		}
	} else {
		err = errors.New(fmt.Sprint("id: ", args.ID, ", org id: ", args.OrgID, ", bind id: ", args.BindID, ", new price: ", args.Price, ", err: ", err))
	}
	return
}

// ArgsUpdateLogNeedAt 修改新的上门时间参数
type ArgsUpdateLogNeedAt struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//预计上门时间
	NeedAt string `db:"need_at" json:"needAt" check:"isoTime" empty:"true"`
	//服务配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
}

// UpdateLogNeedAt 修改新的上门时间
func UpdateLogNeedAt(args *ArgsUpdateLogNeedAt) (err error) {
	//检查预约时间
	var needAt time.Time
	if args.NeedAt != "" {
		needAt, err = CoreFilter.GetTimeByISO(args.NeedAt)
		if err != nil {
			return
		}
	} else {
		if configData, b := CheckConfigTime(&ArgsCheckConfigTime{
			ID: args.ConfigID,
		}); !b {
			err = errors.New("no time")
			return
		} else {
			needAt = configData.StartAt
		}
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_log SET update_at = NOW(), need_at = :need_at WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:bind_id < 1 OR bind_id = :bind_id) AND finish_at < to_timestamp(1000000)", map[string]interface{}{
		"id":      args.ID,
		"org_id":  args.OrgID,
		"bind_id": args.BindID,
		"need_at": needAt,
	})
	if err != nil {
		return
	}
	//更新配置统计
	_ = addConfig(args.ConfigID)
	return
}

// ArgsUpdateLogFinish 标记完成服务参数
type ArgsUpdateLogFinish struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
}

// UpdateLogFinish 标记完成服务
func UpdateLogFinish(args *ArgsUpdateLogFinish) (err error) {
	//更新服务单为完成
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_log SET update_at = NOW(), finish_at = NOW() WHERE id = :id AND pay_at > to_timestamp(1000000) AND (:org_id < 1 OR org_id = :org_id) AND (:bind_id < 1 OR bind_id = :bind_id) AND finish_at < to_timestamp(1000000)", args)
	if err != nil {
		return
	}
	//获取服务单信息
	var data FieldsLog
	data, err = getLogID(args.ID)
	if err != nil {
		return
	}
	//广播服务单
	pushNatsUpdateStatus("finish", data.ID, "服务单完成")
	//反馈
	return
}

// 通知nats更新服务单
func pushNatsUpdateStatus(action string, id int64, des string) {
	CoreNats.PushDataNoErr("service_housekeeping_update", "/service/housekeeping/update", action, id, "", map[string]interface{}{
		"des": des,
	})
}
