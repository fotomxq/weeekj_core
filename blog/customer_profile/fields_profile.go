package BlogCustomerProfile

import (
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsProfile 客户留存资料信息
type FieldsProfile struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//客户姓名
	Name string `db:"name" json:"name"`
	//客户联系地址组件
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//留言信息
	Msg string `db:"msg" json:"msg"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
