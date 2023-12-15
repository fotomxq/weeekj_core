package OrgMission

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

type FieldsMission struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//关联自动化
	AutoID int64 `db:"auto_id" json:"autoID"`
	//状态
	// 0 未完成; 1 已完成; 2 放弃; 3 删除或取消
	Status int `db:"status" json:"status"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//创建人
	CreateBindID int64 `db:"create_bind_id" json:"createBindID"`
	//执行人
	BindID int64 `db:"bind_id" json:"bindID"`
	//其他执行人
	OtherBindIDs pq.Int64Array `db:"other_bind_ids" json:"otherBindIDs"`
	//标题
	Title string `db:"title" json:"title"`
	//描述
	Des string `db:"des" json:"des"`
	//文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//执行时间
	StartAt time.Time `db:"start_at" json:"startAt"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt"`
	//是否需提醒
	// -1 不需要; 0 需要等待提醒中; >0 已经触发提醒的ID
	TipID int64 `db:"tip_id" json:"tipID"`
	//上级任务
	ParentID int64 `db:"parent_id" json:"parentID"`
	//级别
	Level int `db:"level" json:"level"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
