package ServiceHousekeeping

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsBind 服务人员
type FieldsBind struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID"`
	//分区ID
	MapAreaID int64 `db:"map_area_id" json:"mapAreaID"`
	//服务累计收款
	AllTakePrice int64 `db:"all_take_price" json:"allTakePrice"`
	//累计服务次数
	AllLogCount int64 `db:"all_log_count" json:"allLogCount"`
	//累计评分
	AllLevel int64 `db:"all_level" json:"allLevel"`
	//当前未完成任务
	UnFinishCount int64 `db:"un_finish_count" json:"unFinishCount"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
