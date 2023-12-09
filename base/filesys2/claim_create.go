package BaseFileSys2

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

type argsCreateClaim struct {
	//创建组织
	// 可选，指定后该文件归属于组织，用户ID将只是指引，没有操作权限
	OrgID int64 `db:"org_id" json:"orgID"`
	//创建用户
	// 必须指定创建的用户，如果组织失效，则文件将自动归属于用户
	UserID int64 `db:"user_id" json:"userID"`
	//是否公开？
	// 否则必须指定认证来源才能查看
	IsPublic bool `db:"is_public" json:"isPublic"`
	//文件结构体
	FileID int64 `db:"file_id" json:"fileID"`
	//文件自动过期时间
	// 过期将自动销毁该文件
	// null为永远不过期
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//描述或备注
	Des string `db:"des" json:"des"`
	//其他扩展信息
	Infos CoreSQLConfig.FieldsConfigsType `db:"infos" json:"infos"`
}

func createClaim(args *argsCreateClaim) (newClaimData FieldsFileClaim, newCoreData FieldsFile, errCode string, err error) {
	err = claimDB.Insert().SetFields([]string{"org_id", "user_id", "is_public", "file_id", "expire_at", "des", "infos"}).Add(args).ExecAndCheckID()
	if err != nil {
		return
	}
	return
}
