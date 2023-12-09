package BlogStuRead

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteContent 删除词条参数
type ArgsDeleteContent struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteContent 删除词条
func DeleteContent(args *ArgsDeleteContent) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "blog_stu_read_content", "id = :id AND (org_id = :org_id OR :org_id < 1)", args)
	if err != nil {
		return
	}
	deleteContentCacheByID(args.ID)
	return
}
