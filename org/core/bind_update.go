package OrgCoreCore

import (
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// 更新登陆时间
type argsUpdateBindLastTime struct {
	//绑定ID
	ID int64 `db:"id" json:"id" check:"id"`
}

func updateBindLastTime(args *argsUpdateBindLastTime) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core_bind SET last_at = NOW() WHERE id = :id", args)
	if err != nil {
		return
	}
	deleteBindCache(args.ID)
	return
}

// 更新统计信息
func updateOrgBindAnalysis(orgID int64) {
	var count int64
	//机构人数
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM org_core_bind WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", orgID)
	AnalysisAny2.AppendData("re", "org_bind_count", time.Time{}, orgID, 0, 0, 0, 0, count)
}
