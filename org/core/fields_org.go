package OrgCoreCore

import (
	"github.com/lib/pq"
	"time"
)

//FieldsOrg 组织主表
type FieldsOrg struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//所属用户
	// 掌管该数据的用户，创建人和根管理员，不可删除只能更换
	UserID int64 `db:"user_id" json:"userID"`
	//企业唯一标识码
	// 用于特殊识别和登陆识别等操作
	Key string `db:"key" json:"key"`
	//构架名称，或组织名称
	Name string `db:"name" json:"name"`
	//组织描述
	Des string `db:"des" json:"des"`
	//上级ID
	ParentID int64 `db:"parent_id" json:"parentID"`
	//上级控制权限限制
	ParentFunc pq.StringArray `db:"parent_func" json:"parentFunc"`
	//开通业务
	// 该内容只有总管理员或订阅能进行控制
	OpenFunc pq.StringArray `db:"open_func" json:"openFunc"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
}
