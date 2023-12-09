package MapArea

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLGPS "gitee.com/weeekj/weeekj_core/v5/core/sql/gps"
	"time"
)

type FieldsArea struct {
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
	//分区标识码
	// 分类标记
	Mark string `db:"mark" json:"mark"`
	//归属关系
	// 可以作为行政分区和下级配送分区关系的设置，只有平台方能设置没有上级的分区
	// 其他分区必须指定行政分区作为上级，否则无法建立分区
	// 上级分区必须同属一个城市，且所有点不能超越范围
	ParentID int64 `db:"parent_id" json:"parentID"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//所属城市
	City int `db:"city" json:"city"`
	//地图制式
	// 0 / 1 / 2 / 3
	// WGS-84 / GCJ-02 / BD-09 / 2000-china
	MapType int `db:"map_type" json:"mapType"`
	//坐标系
	Points CoreSQLGPS.FieldsPoints `db:"points" json:"points"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
