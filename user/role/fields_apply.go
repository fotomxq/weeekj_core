package UserRole

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsApply 申请角色
type FieldsApply struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//审核时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//申请描述
	AuditDes string `db:"audit_des" json:"auditDes" check:"des" min:"1" max:"600" empty:"true"`
	//拒绝原因
	AuditBanDes string `db:"audit_ban_des" json:"auditBanDes" check:"des" min:"1" max:"600" empty:"true"`
	//角色类型
	RoleType int64 `db:"role_type" json:"roleType" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//姓名
	Name string `db:"name" json:"name" check:"name"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//城市编码
	City string `db:"city" json:"city" check:"cityCode"`
	//性别
	// 0 男 1 女 2 未知
	Gender int `db:"gender" json:"gender" check:"gender"`
	//联系电话
	Phone string `db:"phone" json:"phone" check:"phone"`
	//个人照片
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//证件列
	CertFiles pq.Int64Array `db:"cert_files" json:"certFiles" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params" empty:"true"`
}
