package BaseFileSys2

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsFileClaim 认领结构体
// 文件的hash可重复被认领，节约存储空间
// 除了匹配hash外，还会匹配文件的type和size信息，以减少碰撞概率
// 采用sha256作为唯一标识码标准
type FieldsFileClaim struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//验证Hash
	UpdateHash string `db:"update_hash" json:"updateHash"`
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
	//最后访问时间
	VisitLastAt time.Time `db:"visit_last_at" json:"visitLastAt"`
	//访问次数
	VisitCount int `db:"visit_count" json:"visitCount"`
	//描述或备注
	Des string `db:"des" json:"des"`
	//其他扩展信息
	Infos CoreSQLConfig.FieldsConfigsType `db:"infos" json:"infos"`
}
