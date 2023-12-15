package OrgCoreCore

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsDeleteBind 删除某个绑定关系参数
type ArgsDeleteBind struct {
	//绑定ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 如果是给非ID，则必须给组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteBind 删除某个绑定关系
func DeleteBind(args *ArgsDeleteBind) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_core_bind", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//更新统计
	updateOrgBindAnalysis(args.OrgID)
	//推送NATS
	CoreNats.PushDataNoErr("/org/core/bind", "delete", args.ID, "bind", nil)
	//清理缓冲
	deleteBindCache(args.ID)
	return
}

// ArgsDeleteBindByUser 通过用户删除绑定关系
type ArgsDeleteBindByUser struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteBindByUser 删除某个绑定关系
func DeleteBindByUser(args *ArgsDeleteBindByUser) (err error) {
	//重新获取数据，进行推送数据
	var bindID int64
	err = Router2SystemConfig.MainDB.Get(&bindID, "SELECT id FROM org_core_bind WHERE user_id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", args.UserID, args.OrgID)
	if err != nil {
		return
	}
	//执行删除
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_core_bind", "user_id = :user_id AND org_id = :org_id", args)
	if err != nil {
		return
	}
	//更新统计
	updateOrgBindAnalysis(args.OrgID)
	//推送NATS
	CoreNats.PushDataNoErr("/org/core/bind", "delete", args.OrgID, "user", map[string]int64{
		"userID": args.UserID,
	})
	//清理缓冲
	deleteBindCache(bindID)
	return
}

// ArgsDeleteBindByOrg 通过组织删除成员关系参数
type ArgsDeleteBindByOrg struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteBindByOrg 通过组织删除成员关系
func DeleteBindByOrg(args *ArgsDeleteBindByOrg) (err error) {
	//重新获取数据，进行推送数据
	var bindIDs pq.Int64Array
	err = Router2SystemConfig.MainDB.Get(&bindIDs, "SELECT id FROM org_core_bind WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", args.OrgID)
	if err != nil {
		return
	}
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_core_bind", "org_id = :org_id", args)
	if err != nil {
		return
	}
	//更新统计
	updateOrgBindAnalysis(args.OrgID)
	//推送NATS
	CoreNats.PushDataNoErr("/org/core/bind", "delete", args.OrgID, "all", nil)
	//清理缓冲
	for _, v := range bindIDs {
		deleteBindCache(v)
	}
	return
}

// ArgsReturnBind 恢复绑定关系参数
type ArgsReturnBind struct {
	//绑定ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 如果是给非ID，则必须给组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// ReturnBind 恢复绑定关系
func ReturnBind(args *ArgsReturnBind) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_core_bind SET delete_at = to_timestamp(0) WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//更新统计
	updateOrgBindAnalysis(args.OrgID)
	//清理缓冲
	deleteBindCache(args.ID)
	return
}
