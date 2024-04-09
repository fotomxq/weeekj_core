package BlogCore

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteContent 删除词条参数
type ArgsDeleteContent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 文章可以是由用户发出，组织ID可以为0
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// DeleteContent 删除词条
func DeleteContent(args *ArgsDeleteContent) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "blog_core_content", "id = :id AND (org_id = :org_id OR :org_id < 1) AND (:user_id < 1 OR user_id = :user_id) AND (:bind_id < 1 OR bind_id = :bind_id)", args)
	if err != nil {
		return
	}
	deleteContentCacheByID(args.ID)
	CoreNats.PushDataNoErr("blog_core_delete", "/blog/core/delete", "", args.ID, "", nil)
	return
}

// ArgsReturnContent 还原内容参数
type ArgsReturnContent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 文章可以是由用户发出，组织ID可以为0
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// ReturnContent 还原内容
func ReturnContent(args *ArgsReturnContent) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE blog_core_content SET delete_at = TO_TIMESTAMP(0), audit_at = TO_TIMESTAMP(0), publish_at = TO_TIMESTAMP(0) WHERE id = :id AND (org_id = :org_id OR :org_id < 1) AND (:user_id < 1 OR user_id = :user_id)", args)
	if err != nil {
		return
	}
	deleteContentCacheByID(args.ID)
	return
}
