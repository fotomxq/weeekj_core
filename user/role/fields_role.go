package UserRole

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsRole 角色信息
type FieldsRole struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//角色类型
	RoleType int64 `db:"role_type" json:"roleType" check:"id"`
	//申请ID
	ApplyID int64 `db:"apply_id" json:"applyID" check:"id"`
	//用户ID
	// 允许为0，则该信息不属于任何用户，或不和任何用户关联
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
