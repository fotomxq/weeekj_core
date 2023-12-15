package MarketCore

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"github.com/lib/pq"
	"time"
)

// FieldsBind 营销成员和客户关系
type FieldsBind struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//绑定的用户
	BindUserID int64 `db:"bind_user_id" json:"bindUserID"`
	//绑定的档案
	BindInfoID int64 `db:"bind_info_id" json:"bindInfoID"`
	//建立关系的渠道
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//客户备注
	Des string `db:"des" json:"des"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
