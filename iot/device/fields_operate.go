package IOTDevice

import (
	"time"

	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
)

// FieldsOperate 授权信息
type FieldsOperate struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//授权过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//权限类型
	// all 全部权限; read 允许查看设备信息; write 允许编辑设备信息; mission 任务下达权限; operate 修改授权关系; associated 关联设备
	Permissions pq.StringArray `db:"permissions" json:"permissions"`
	//允许执行的动作
	// 将根据设备组的动作查询，如果存在则允许，否则将禁止执行该类动作
	Action pq.Int64Array `db:"action" json:"action"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//注册地
	Address string `db:"address" json:"address"`
	//组织分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//组织标签ID组
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
