package ServiceOrder

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateFailed 更新到失败参数
type ArgsUpdateFailed struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//调整说明描述
	Des string `db:"des" json:"des"`
}

// UpdateFailed 更新到失败
func UpdateFailed(args *ArgsUpdateFailed) (err error) {
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "failed", args.Des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), delete_at = NOW(), status = 5, logs = logs || :log WHERE id = :id AND (status != 5 OR status != 6) AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", map[string]interface{}{
		"id":      args.ID,
		"org_id":  args.OrgID,
		"user_id": args.UserID,
		"log":     newLog,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(args.ID)
	//发起退款请求
	if _, err2 := RefundPay(&ArgsRefundPay{
		ID:          args.ID,
		OrgID:       args.OrgID,
		UserID:      args.UserID,
		OrgBindID:   args.OrgBindID,
		RefundPrice: -1,
		Des:         "订单失败",
	}); err2 != nil {
		CoreLog.Error("update order failed, refund pay, order id: ", args.ID, ", err: ", err2)
	}
	//反馈
	return
}
