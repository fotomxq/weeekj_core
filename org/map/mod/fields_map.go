package OrgMapMod

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsMap 商户地址结构
type FieldsMap struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//审核时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//上级ID
	// 用于叠加展示
	ParentID int64 `db:"parent_id" json:"parentID"`
	//展示小图标
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//展示信息
	Name string `db:"name" json:"name"`
	//展示介绍信息
	Des string `db:"des" json:"des"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province"`
	//所属城市
	City int `db:"city" json:"city"`
	//街道详细信息
	Address string `db:"address" json:"address"`
	//地图制式
	// 0 WGS-84 / 1 GCJ-02 / 2 BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//广告点击次数
	AdCount int64 `db:"ad_count" json:"adCount"`
	//广告可用次数
	AdCountLimit int64 `db:"ad_count_limit" json:"adCountLimit"`
	//查看最短时间长度
	ViewTimeLimit int64 `db:"view_time_limit" json:"viewTimeLimit"`
	//扩展参数
	// adFileID 广告文件ID，用于投放广告
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
