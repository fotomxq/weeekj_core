package ClassQueue

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

type FieldsQueue struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//其他模块的ID
	ModID int64 `db:"mod_id" json:"modID"`
	//处理状态
	// 如果消息件存在多个状态，可使用，否则应及时删除该消息
	Status int `db:"status" json:"status"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
