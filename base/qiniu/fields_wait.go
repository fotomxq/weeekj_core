package BaseQiniu

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"time"
)

type FieldsWait struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	// 超出一定时间会自动删除该数据
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//是否公开
	IsPublic bool `db:"is_public" json:"isPublic"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//扩展参数
	ClaimInfos CoreSQLConfig.FieldsConfigsType `db:"claim_infos" json:"claimInfos"`
	//描述
	Des string `db:"des" json:"des"`
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//IP
	IP string `db:"ip" json:"ip"`
	//绑定的文件ID
	FileClaimID int64 `db:"file_claim_id" json:"fileClaimID"`
}
