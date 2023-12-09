package ServiceInfoExchange

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsTake 参与和领取
type FieldsTake struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//参与用户
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//信息ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id"`
	//备注信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"500" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
