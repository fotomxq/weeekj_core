package TMSTransport

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsBind 分区和配送人员绑定关系
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
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID"`
	//分区ID
	MapAreaID int64 `db:"map_area_id" json:"mapAreaID"`
	//更多分区
	// 可以绑定更多分区，但分配性能会下降
	MoreMapAreaIDs pq.Int64Array `db:"more_map_area_ids" json:"moreMapAreaIDs" check:"ids" empty:"true"`
	//最近30天评价
	Level30Day int `db:"level_30_day" json:"level30Day"`
	//最近30天里程数
	KM30Day int `db:"km_30_day" json:"km30Day"`
	//最近30天累计任务累计耗时
	Time30Day int64 `db:"time_30_day" json:"time30Day"`
	//最近30天任务量
	Count30Day int `db:"count_30_day" json:"count30Day"`
	//最近30天完成任务量
	CountFinish30Day int `db:"count_finish_30_day" json:"countFinish30Day"`
	//最近1天评价等数据
	Level1Day       int   `db:"level_1_day" json:"level1Day"`
	KM1Day          int   `db:"km_1_day" json:"km1Day"`
	Time1Day        int64 `db:"time_1_day" json:"time1Day"`
	Count1Day       int   `db:"count_1_day" json:"count1Day"`
	CountFinish1Day int   `db:"count_finish_1_day" json:"countFinish1Day"`
	//当前未完成任务
	UnFinishCount int64 `db:"un_finish_count" json:"unFinishCount"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
