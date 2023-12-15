package ServiceUserInfo

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteInfo 删除列参数
type ArgsDeleteInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteInfo 删除列
func DeleteInfo(args *ArgsDeleteInfo) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_user_info", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//缓冲
	deleteInfoCache(args.ID)
	//推送nats
	pushNatsInfoStatus("delete", args.ID)
	//统计数据
	pushNatsAnalysis(args.OrgID)
	//日志
	appendLog(&argsAppendLog{
		InfoID:     args.ID,
		OrgID:      0,
		ChangeMark: "delete",
		ChangeDes:  "删除档案",
		OldDes:     "",
		NewDes:     "",
	})
	//反馈
	return
}

// ArgsReturnInfo 还原档案参数
type ArgsReturnInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// ReturnInfo 还原档案
func ReturnInfo(args *ArgsReturnInfo) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_user_info SET delete_at = to_timestamp(0), die_at = to_timestamp(0), out_at = to_timestamp(0) WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//缓冲
	deleteInfoCache(args.ID)
	//推送nats
	pushNatsInfoStatus("return", args.ID)
	//统计数据
	pushNatsAnalysis(args.OrgID)
	//日志
	appendLog(&argsAppendLog{
		InfoID:     args.ID,
		OrgID:      0,
		ChangeMark: "return",
		ChangeDes:  "还原档案",
		OldDes:     "",
		NewDes:     "",
	})
	//反馈
	return
}
