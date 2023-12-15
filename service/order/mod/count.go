package ServiceOrderMod

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetUserOrderCount 获取用户下的订单数量参数
type ArgsGetUserOrderCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//是否完成
	IsFinish bool `db:"is_finish" json:"isFinish"`
}

// GetUserOrderCount 获取用户下的订单数量
func GetUserOrderCount(args *ArgsGetUserOrderCount) (count int64, err error) {
	if args.IsFinish {
		count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "service_order", "id", "org_id = :org_id AND user_id = :user_id AND status = 4", map[string]interface{}{
			"org_id":  args.OrgID,
			"user_id": args.UserID,
		})
		return
	} else {
		count, err = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "service_order", "id", "org_id = :org_id AND user_id = :user_id", map[string]interface{}{
			"org_id":  args.OrgID,
			"user_id": args.UserID,
		})
		return
	}
}
