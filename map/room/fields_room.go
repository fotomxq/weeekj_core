package MapRoom

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

type FieldsRoom struct {
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
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//状态
	// 0 空闲; 1 有人; 2 退房; 3 不可用; 4 清理中
	Status int `db:"status" json:"status"`
	//服务呼叫状态
	// 0 no 无呼叫; 1 call 正在呼叫; 2 ok 已经应答并处置
	// 处置完成后将回归0状态
	ServiceStatus int `db:"service_status" json:"serviceStatus"`
	//服务工作人员
	ServiceBindID int64 `db:"service_bind_id" json:"serviceBindID"`
	//联动行政任务
	// 任务完成后将自动清除为0，否则将一直挂起
	ServiceMissionID int64 `db:"service_mission_id" json:"serviceMissionID"`
	//入驻人员列
	Infos pq.Int64Array `db:"infos" json:"infos"`
	//房间编号
	Code string `db:"code" json:"code"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//描述图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//背景色
	BgColor string `db:"bg_color" json:"bgColor" check:"color" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
